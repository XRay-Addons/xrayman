package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type Service struct {
	cfg XRayCfg
	api XRayApi
	ctl XRayCtl
	perfCtl PerfCtl

	mu      sync.Mutex
}

func New(
	cfg XRayCfg,
	ctl XRayCtl,
	api XRayApi,
	perfCtl PerfCtl,
) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%w: cfg not exists", errdefs.ErrIPE)
	}
	if ctl == nil {
		return nil, fmt.Errorf("%w: ctl not exists", errdefs.ErrIPE)
	}
	if api == nil {
		return nil, fmt.Errorf("%w: api not exists", errdefs.ErrIPE)
	}
	if perfCtl == nil {
		return nil, fmt.Errorf("%w: perfctl not exists", errdefs.ErrIPE)
	}
	return &Service{
		cfg:     cfg,
		api:     api,
		ctl:     ctl,
		perfCtl: perfCtl,
	}, nil
}

func (s *Service) Start(ctx context.Context, users []models.User) (*models.Node, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// update config
	serverCfg, err := s.cfg.GetServerConfig(users)
	if err != nil {
		return nil, fmt.Errorf("get users config: %w", err)
	}

	// start service
	if err := s.ctl.Start(ctx, serverCfg); err != nil {
		return nil, fmt.Errorf("restart service: %w", err)
	}

	// ping service after restart, several attempt required
	restartPingRetries := 5
	if err := s.ping(ctx, restartPingRetries); err != nil {
		return nil, fmt.Errorf("ping service: %w", err)
	}
	return &models.Node{
		ClientConfig: s.cfg.GetClientConfig(),
	}, nil
}

func (s *Service) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ctl.Stop(ctx); err != nil {
		return fmt.Errorf("stop xray service: %w", err)
	}

	return nil
}

func (s *Service) Status(ctx context.Context) (*models.NodeStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, err := s.ctl.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("get xray status: %w", err)
	}
	cpuLoad, err := s.perfCtl.GetCPUUsage()
	if err != nil {
		return nil, fmt.Errorf("get cpu usage: %w", err)
	}
	ramLoad, err := s.perfCtl.GetRAMUsage()
	if err != nil {
		return nil, fmt.Errorf("get ram usage: %w", err)
	}

	return &models.NodeStatus{
		Status: status,
		CPULoad: cpuLoad,
		RAMLoad: ramLoad,
	}, nil
}

func (s *Service) AddUsers(ctx context.Context, users []models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.api.AddUsers(ctx, s.cfg.GetInbounds(), users); err != nil {
		return fmt.Errorf("add xray users: %w", err)
	}

	return nil
}

func (s *Service) DelUsers(ctx context.Context, users []models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.api.DelUsers(ctx, s.cfg.GetInbounds(), users); err != nil {
		return fmt.Errorf("del xray users: %w", err)
	}

	return nil
}

func (s *Service) ping(ctx context.Context, retries int) error {
	var err error
	for range retries {
		delay := 1 * time.Second
		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: ping attempts cancelled", errdefs.ErrCancelled)
		case <-timer.C:
			if err = s.api.Ping(ctx); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("ping last try: %w", err)
}
