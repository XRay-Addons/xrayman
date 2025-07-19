package servicectl

import (
	"context"
	"fmt"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"go.uber.org/zap"
)

type XRayServiceCtl struct {
	configPath string
	supervisor supervisor.Supervisor
}

const serviceName = "xray"

func New(execPath, configPath string, log *zap.Logger) (*XRayServiceCtl, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: xray service init: logger", errdefs.ErrNilArgPassed)
	}
	command := []string{execPath, "-config", configPath}
	supervisor, err := supervisor.New(serviceName, command, log)
	if err != nil {
		return nil, fmt.Errorf("xray service init: %w", err)
	}
	return &XRayServiceCtl{
		configPath: configPath,
		supervisor: supervisor,
	}, nil
}

func (s *XRayServiceCtl) Close(ctx context.Context) error {
	if s == nil || s.supervisor == nil {
		return nil
	}
	if err := s.supervisor.Close(ctx); err != nil {
		return fmt.Errorf("xray service: close: %w", err)
	}
	return nil
}

func (s *XRayServiceCtl) Start(ctx context.Context, config string) error {
	if s == nil || s.supervisor == nil {
		return fmt.Errorf("%w: xray service: start", errdefs.ErrNilObjectCall)
	}
	err := os.WriteFile(s.configPath, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("%w: xray service: start: write config: %v", errdefs.ErrAccess, err)
	}
	if err := s.supervisor.Start(ctx); err != nil {
		return fmt.Errorf("xray service: start: %w", err)
	}
	return nil
}

func (s *XRayServiceCtl) Stop(ctx context.Context) error {
	if s == nil || s.supervisor == nil {
		return fmt.Errorf("%w: xray service: stop", errdefs.ErrNilObjectCall)
	}
	if err := s.supervisor.Stop(ctx); err != nil {
		return fmt.Errorf("xray service: stop: %w", err)
	}
	return nil
}

func (s *XRayServiceCtl) Status(ctx context.Context) (models.ServiceStatus, error) {
	if s == nil || s.supervisor == nil {
		return models.ServiceStopped, fmt.Errorf(
			"%w: xray service stop", errdefs.ErrNilObjectCall)
	}
	status, err := s.supervisor.Status(ctx)
	if err != nil {
		return models.ServiceStopped, fmt.Errorf("xray service status: %w", err)
	}
	return status, nil
}
