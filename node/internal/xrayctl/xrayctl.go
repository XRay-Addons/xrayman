package xrayctl

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/shared/models"
	"github.com/coreos/go-systemd/v22/dbus"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

type XRayCtl struct {
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	cfgPath string

	mu sync.Mutex
}

const (
	xrayService = "xray.service"
)

func New(execPath, cfgPath string, log *zap.Logger) (*XRayCtl, error) {
	if log == nil {
		log = zap.NewNop()
	}

	serviceCfg, err := createServiceConfig(execPath, cfgPath)

	if err != nil {
		return nil, fmt.Errorf("create service config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	xrayCtl := &XRayCtl{
		cancel:  cancel,
		cfgPath: cfgPath,
	}

	// run connection loop
	xrayCtl.wg.Add(1)
	go xrayCtl.initServiceLoop(ctx, serviceCfg, log)

	return xrayCtl, nil
}

func (ctl *XRayCtl) Close(ctx context.Context) error {
	ctl.cancel()
	ctl.wg.Wait()

	if err := ctl.Stop(ctx); err != nil {
		return fmt.Errorf("close ctl: %w", err)
	}
	return nil
}

func (ctl *XRayCtl) Start(ctx context.Context, config string) error {
	err := ctl.runCtlFn(ctx, func(conn *dbus.Conn) (err error) {
		if err := os.WriteFile(ctl.cfgPath, []byte(config), 0o644); err != nil {
			return fmt.Errorf("%w: write xray config: %v", errdefs.ErrAccess, err)
		}

		systemdCh := make(chan string)
		_, err = conn.RestartUnitContext(ctx, xrayService, "replace", systemdCh)
		if err != nil {
			return fmt.Errorf("%w: start service: %v", errdefs.ErrService, err)
		}

		select {
		case status := <-systemdCh:
			if status == "done" {
				return nil
			}
			return fmt.Errorf("%w: start service: %s", errdefs.ErrService, status)
		case <-ctx.Done():
			return fmt.Errorf("%w: start service cancelled", errdefs.ErrService)
		}
	})
	if err != nil {
		return fmt.Errorf("start service: %w", err)
	}
	return nil
}

func (ctl *XRayCtl) Stop(ctx context.Context) error {
	err := ctl.runCtlFn(ctx, func(conn *dbus.Conn) error {
		systemdCh := make(chan string)
		_, err := conn.StopUnitContext(ctx, xrayService, "replace", systemdCh)
		if err != nil {
			return fmt.Errorf("%w: stop service: %v", errdefs.ErrService, err)
		}

		select {
		case status := <-systemdCh:
			if status == "done" {
				return nil
			}
			return fmt.Errorf("%w: start service: %s", errdefs.ErrService, status)
		case <-ctx.Done():
			return fmt.Errorf("%w: start service cancelled", errdefs.ErrService)
		}
	})
	if err != nil {
		return fmt.Errorf("stop service: %w", err)
	}
	return nil
}

func (ctl *XRayCtl) Status(ctx context.Context) (models.Status, error) {
	var isRunning bool
	err := ctl.runCtlFn(ctx, func(conn *dbus.Conn) error {
		unitStatus, err := conn.GetUnitPropertiesContext(ctx, xrayService)
		if err != nil {
			return fmt.Errorf("%w: is service running", errdefs.ErrService)
		}

		activeState, ok := unitStatus["ActiveState"].(string)
		if !ok {
			return fmt.Errorf("%w: ActiveState property not found", errdefs.ErrService)
		}
		isRunning = activeState == "active"
		return nil
	})
	if err != nil {
		return models.NotRunning, fmt.Errorf("is running service: %w", err)
	}
	if !isRunning {
		return models.NotRunning, nil
	}
	return models.Running, nil
}

type ServiceCfg struct {
	UnitSection struct {
		Description string `ini:"Description"`
		After       string `ini:"After"`
	} `ini:"Unit"`

	ServiceSection struct {
		Type      string `ini:"Type"`
		ExecStart string `ini:"ExecStart"`
	} `ini:"Service"`
}

func createServiceConfig(execPath, cfgPath string) (string, error) {
	cfg := ini.Empty()

	var serviceCfg ServiceCfg
	serviceCfg.UnitSection.Description = "xray service"
	serviceCfg.UnitSection.After = "network.target"
	serviceCfg.ServiceSection.Type = "simple"
	serviceCfg.ServiceSection.ExecStart = fmt.Sprintf("%s  run -config %s", execPath, cfgPath)

	if err := cfg.ReflectFrom(&serviceCfg); err != nil {
		return "", fmt.Errorf("%w: create config: %v", errdefs.ErrIPE, err)
	}

	var buf bytes.Buffer
	if _, err := cfg.WriteTo(&buf); err != nil {
		return "", fmt.Errorf("%w: config to buffer: %v", errdefs.ErrIPE, err)
	}

	return buf.String(), nil
}

func (ctl *XRayCtl) initServiceLoop(ctx context.Context, serviceCfg string, log *zap.Logger) {
	defer ctl.wg.Done()

	if ctl.initService(ctx, serviceCfg) == nil {
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ctl.initService(ctx, serviceCfg); err != nil {
				log.Error("init service", zap.Error(err))
			} else {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

const systemdConfsPath = "/run/systemd/system"

func (ctl *XRayCtl) initService(ctx context.Context, serviceCfg string) error {
	ctl.mu.Lock()
	defer ctl.mu.Unlock()

	confPath := filepath.Join(systemdConfsPath, xrayService)
	if err := os.WriteFile(confPath, []byte(serviceCfg), 0o644); err != nil {
		return fmt.Errorf("%w: write xray service config: %v", errdefs.ErrAccess, err)
	}

	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: dbus connection: %v", errdefs.ErrService, err)
	}
	if err := conn.ReloadContext(ctx); err != nil {
		return fmt.Errorf("%w: dbus reload: %v", errdefs.ErrService, err)
	}

	return nil
}

type ctlFn = func(conn *dbus.Conn) error

func (ctl *XRayCtl) runCtlFn(ctx context.Context, fn ctlFn) error {
	ctl.mu.Lock()
	defer ctl.mu.Unlock()

	// TODO: add retr
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: dbus connection: %v", errdefs.ErrService, err)
	}
	if err := fn(conn); err != nil {
		return fmt.Errorf("%w: service fn: %v", errdefs.ErrService, err)
	}
	return nil
}
