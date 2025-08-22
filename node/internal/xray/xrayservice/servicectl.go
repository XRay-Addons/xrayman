package xrayservice

import (
	"context"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor"
	"github.com/XRay-Addons/xrayman/node/internal/infra/supervisor/supervisorapi"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"go.uber.org/zap"
)

type XRayService struct {
	configPath string
	supervisor supervisorapi.Supervisor
}

const serviceName = "xray"

func New(execPath, configPath string, log *zap.Logger) (*XRayService, error) {
	if log == nil {
		return nil, errdefs.NewNilCall()
	}
	command := []string{execPath, "-config", configPath}
	supervisor, err := supervisor.New(serviceName, command, log)
	if err != nil {
		return nil, err
	}
	return &XRayService{
		configPath: configPath,
		supervisor: supervisor,
	}, nil
}

func (s *XRayService) Close(ctx context.Context) error {
	if s == nil || s.supervisor == nil {
		return nil
	}
	if err := s.supervisor.Close(ctx); err != nil {
		return err
	}
	return nil
}

func (s *XRayService) Start(ctx context.Context, config string) error {
	if s == nil || s.supervisor == nil {
		return errdefs.NewNilCall()
	}
	err := os.WriteFile(s.configPath, []byte(config), 0644)
	if err != nil {
		return errdefs.Wrap(err, errdefs.WithStack(),
			errdefs.WithFile(s.configPath))
	}
	if err := s.supervisor.Start(ctx); err != nil {
		return err
	}
	return nil
}

func (s *XRayService) Stop(ctx context.Context) error {
	if s == nil || s.supervisor == nil {
		return errdefs.NewNilCall()
	}
	if err := s.supervisor.Stop(ctx); err != nil {
		return err
	}
	return nil
}

func (s *XRayService) Status(ctx context.Context) (models.ServiceStatus, error) {
	if s == nil || s.supervisor == nil {
		return models.ServiceStopped, errdefs.NewNilCall()
	}
	status, err := s.supervisor.Status(ctx)
	if err != nil {
		return models.ServiceStopped, err
	}
	serviceStatus, err := convertStatus(status)
	if err != nil {
		return models.ServiceStopped, err
	}
	return serviceStatus, nil
}
