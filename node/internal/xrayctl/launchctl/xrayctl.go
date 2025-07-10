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
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/internal/retry"

	"go.uber.org/zap"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>xray</string>
    <key>ProgramArguments</key>
    <array>
      <string>{{.ExecPath}}</string>
      <string>-config</string>
      <string>{{.ConfigPath}}</string>
    </array>
    <key>RunAtLoad</key>
    <false/>
    
    <key>ProcessType</key>
    <string>Background</string>
  </dict>
</plist>`

const serviceName = "xray"

const plistLocation = "/Library/LaunchAgents/xray.plist"

const statusRegex = `(?m)^[ \t]*state = (.+?)$`

type XRayCtl struct {
	userDomain    string // gui/501
	plistLocation string // /Users/user/<plistLocation>

	// for initialization loop
	initialized atomic.Bool
	wg          sync.WaitGroup
	cancel      context.CancelFunc
}

func New(execPath, cfgPath string, log *zap.Logger) (*XRayCtl, error) {
	if log == nil {
		log = zap.NewNop()
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
		userDomain:    userDomain,
		plistLocation: filepath.Join(userHome, plistLocation),
		cancel:        cancel,
	}

	// init plist
	if err := xrayCtl.createPlistFile(execPath, cfgPath); err != nil {
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
		if err := stopService(ctl.userDomain, serviceName); err != nil {
			closeErrs = append(closeErrs, err)
		}
		if err := removeService(ctl.userDomain, serviceName); err != nil {
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

func (ctl *XRayCtl) createPlistFile(execPath, cfgPath string) error {
	plistContent, err := makePlist(execPath, cfgPath)
	if err != nil {
		return fmt.Errorf("create plist file: %w", err)
	}
	if err := os.WriteFile(ctl.plistLocation, plistContent, 0644); err != nil {
		return fmt.Errorf("write service plist: %w: %v", errdefs.ErrAccess, err)
	}
	return nil
}

func makePlist(execPath, cfgPath string) ([]byte, error) {
	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse service plist: %w: %v", errdefs.ErrIPE, err)
	}

	data := map[string]string{
		"ExecPath":   execPath,
		"ConfigPath": cfgPath,
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
		err := removeService(ctl.userDomain, serviceName)
		if err != nil {
			err = fmt.Errorf("init service: %w", err)
			log.Warn(err.Error())
			return err
		}

		// create new service
		err = createService(ctl.userDomain, ctl.plistLocation)
		if err != nil {
			err = fmt.Errorf("init service: %w", err)
			log.Warn(err.Error())
			return err
		}

		// mark as initialized
		ctl.initialized.Store(true)

		return nil
	}

	retry.RetryInfinite(ctx, initFn, 1*time.Second)
}

func (ctl *XRayCtl) Start(ctx context.Context) error {
	if err := ctl.checkServiceReady(); err != nil {
		return fmt.Errorf("xray start: %w", err)
	}

	// send start signal to service
	return startService(ctl.userDomain, serviceName)
}

func (ctl *XRayCtl) Stop(ctx context.Context) error {
	if err := ctl.checkServiceReady(); err != nil {
		return fmt.Errorf("xray stop: %w", err)
	}

	return stopService(ctl.userDomain, serviceName)
}

func (ctl *XRayCtl) Status(ctx context.Context) (models.ServiceStatus, error) {
	if err := ctl.checkServiceReady(); err != nil {
		return models.ServiceStopped, fmt.Errorf("xray status: %w", err)
	}

	return getServiceStatus(ctl.userDomain, serviceName)
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
