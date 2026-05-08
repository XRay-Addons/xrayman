package xrayservice

import (
	"context"
	"sync"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	xray "github.com/xtls/libxray/xray"
	"go.uber.org/zap"
)

type XRayService struct {
	log *zap.Logger
	mu  sync.RWMutex
}

// TODO: WithLogger
func New(log *zap.Logger) (*XRayService, error) {
	if log == nil {
		return nil, errdefs.NilCall()
	}

	return &XRayService{
		log: log,
	}, nil
}

func (s *XRayService) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := xray.StopXray(); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

func (s *XRayService) Start(ctx context.Context, config string) error {
	if s == nil {
		return errdefs.NilCall()
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// close existed service
	if err := xray.StopXray(); err != nil {
		err = xerr.WrapWithStack(err)
		s.log.Warn("xray stopping error, leak possible", zap.Error(err))
	}

	// start new instance
	err := xray.RunXrayFromJSON("", "", config)
	if err != nil {
		return err
	}

	return nil
}

func (s *XRayService) Stop(ctx context.Context) error {
	if s == nil {
		return errdefs.NilCall()
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := xray.StopXray(); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

func (s *XRayService) Status(ctx context.Context) (models.ServiceStatus, error) {
	if s == nil {
		return models.ServiceStatusStopped, errdefs.NilCall()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	if xray.GetXrayState() {
		return models.ServiceStatusRunning, nil
	} else {
		return models.ServiceStatusStopped, nil
	}
}

// TODO: TestConfig for bootstrap
