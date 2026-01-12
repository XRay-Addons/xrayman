package launchctl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/retry"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/supervisorapi"

	"go.uber.org/zap"
)

type XRayCtl struct {
	serviceName   string
	userDomain    string // gui/501
	plistLocation string // /Users/user/<plistLocation>

	// for initialization loop
	initialized atomic.Bool
	wg          sync.WaitGroup
	cancel      context.CancelFunc
}

const plistDirectory = "/Library/LaunchAgents"

func New(serviceName string, command []string, log *zap.Logger) (*XRayCtl, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: init service: logger", errdefs.ErrNilObjectCall)
	}
	userDomain, err := userDomain()
	if err != nil {
		return nil, fmt.Errorf("init service: %w", err)
	}
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("init service: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	xrayCtl := &XRayCtl{
		serviceName:   serviceName + ".service",
		userDomain:    userDomain,
		plistLocation: filepath.Join(userHome, plistDirectory, serviceName+".plist"),
		cancel:        cancel,
	}

	// init plist
	if err := xrayCtl.createPlistFile(command); err != nil {
		return nil, fmt.Errorf("init xrayctl: %w", err)
	}

	// run installing service loop
	xrayCtl.wg.Add(1)
	go xrayCtl.createServiceLoop(ctx, log)

	return xrayCtl, nil
}

func (ctl *XRayCtl) Close(ctx context.Context) error {
	if ctl == nil {
		return nil
	}

	ctl.cancel()
	ctl.wg.Wait()

	var closeErrs []error

	// stop and remove service, if initialized
	if ctl.initialized.Load() {
		if err := stopService(ctl.userDomain, ctl.serviceName); err != nil {
			closeErrs = append(closeErrs, err)
		}
		if err := removeService(ctl.userDomain, ctl.serviceName); err != nil {
			closeErrs = append(closeErrs, err)
		}
		ctl.initialized.Store(false)
	}

	// remove service file
	if err := ctl.removePlistFile(); err != nil {
		closeErrs = append(closeErrs, err)
	}

	if closeErrs == nil {
		return nil
	}

	return fmt.Errorf("%w: close service: %v",
		errdefs.ErrService, errors.Join(closeErrs...))
}

func (ctl *XRayCtl) createPlistFile(command []string) error {
	plistContent, err := makePlist(ctl.serviceName, command)
	if err != nil {
		return fmt.Errorf("create plist file: %w", err)
	}
	if err := os.WriteFile(ctl.plistLocation, plistContent, 0644); err != nil {
		return fmt.Errorf("write service plist: %w: %v", errdefs.ErrAccess, err)
	}
	return nil
}

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>{{ .ServiceName }}</string>
    <key>ProgramArguments</key>
    <array>
    {{range .Command}}  <string>{{.}}</string>
    {{end}}</array>
    <key>RunAtLoad</key>
    <false/>
    
    <key>ProcessType</key>
    <string>Background</string>
  </dict>
</plist>`

func makePlist(serviceName string, command []string) ([]byte, error) {
	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse service plist: %w: %v", errdefs.ErrIPE, err)
	}

	data := map[string]interface{}{
		"ServiceName": serviceName,
		"Command":     command,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("make service plist: %w: %v", errdefs.ErrIPE, err)
	}
	return buf.Bytes(), nil
}

func (ctl *XRayCtl) removePlistFile() error {
	err := os.Remove(ctl.plistLocation)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return fmt.Errorf("remove service plist: %w: %v", errdefs.ErrAccess, err)
}

func (ctl *XRayCtl) createServiceLoop(ctx context.Context, log *zap.Logger) {
	defer ctl.wg.Done()

	initFn := func(ctx context.Context) error {
		// remove existed service
		err := removeService(ctl.userDomain, ctl.serviceName)
		if err != nil {
			log.Warn("retry: init service", zap.Error(err))
			return err
		}

		// create new service
		err = createService(ctl.userDomain, ctl.plistLocation)
		if err != nil {
			log.Warn("retry: init service", zap.Error(err))
			return err
		}

		// mark as initialized
		ctl.initialized.Store(true)

		return nil
	}

	retry.RetryInfinite(ctx, initFn, 250*time.Millisecond)
}

func (ctl *XRayCtl) Start(ctx context.Context) error {
	if err := ctl.checkServiceReady(); err != nil {
		return fmt.Errorf("service start: %w", err)
	}

	// send start signal to service
	if err := startService(ctl.userDomain, ctl.serviceName); err != nil {
		return err
	}

	// wait till service status will be Stopped or Started
	var status supervisorapi.ServiceStatus
	checkStatus := func(ctx context.Context) error {
		var err error
		status, err = ctl.Status(ctx)
		return err
	}
	if err := retry.RetryInfinite(ctx, checkStatus, 250*time.Millisecond); err != nil {
		return fmt.Errorf("service starting: %w", err)
	}
	if status != supervisorapi.StatusRunning {
		return fmt.Errorf("%w: failed to start service", errdefs.ErrService)
	}
	return nil
}

func (ctl *XRayCtl) Stop(ctx context.Context) error {
	if err := ctl.checkServiceReady(); err != nil {
		return fmt.Errorf("service stop: %w", err)
	}

	return stopService(ctl.userDomain, ctl.serviceName)
}

func (ctl *XRayCtl) Status(ctx context.Context) (supervisorapi.ServiceStatus, error) {
	if err := ctl.checkServiceReady(); err != nil {
		return supervisorapi.StatusUnknown, fmt.Errorf("service status: %w", err)
	}

	statusStr, err := getServiceStatus(ctl.userDomain, ctl.serviceName)
	if err != nil {
		return supervisorapi.StatusUnknown, fmt.Errorf("service status: %w", err)
	}

	switch statusStr {
	case "running":
		return supervisorapi.StatusRunning, nil
	case "not running":
		return supervisorapi.StatusStopped, nil
	default:
		return supervisorapi.StatusUnknown, fmt.Errorf(
			"%w: unknown ServiceStatus \"%s\"", errdefs.ErrService, statusStr)
	}
}

func (ctl *XRayCtl) checkServiceReady() error {
	if ctl == nil {
		return fmt.Errorf("%w: xrayctl", errdefs.ErrNilObjectCall)
	}
	if !ctl.initialized.Load() {
		return errdefs.ErrServiceNotReady
	}
	return nil
}

func userDomain() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("user domain: %w: %v", errdefs.ErrAccess, err)
	}
	return "gui/" + u.Uid, nil
}
