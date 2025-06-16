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

	running bool
	mu      sync.Mutex
}

func New(cfg XRayCfg, api XRayApi, ctl XRayCtl) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%w: cfg not exists", errdefs.ErrIPE)
	}
	if api == nil {
		return nil, fmt.Errorf("%w: api not exists", errdefs.ErrIPE)
	}
	if ctl == nil {
		return nil, fmt.Errorf("%w: ctl not exists", errdefs.ErrIPE)
	}
	return &Service{
		cfg:     cfg,
		api:     api,
		ctl:     ctl,
		running: false,
	}, nil
}

func (s *Service) Start(ctx context.Context, users []models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// update config
	if err := s.cfg.SetUsers(users); err != nil {
		return fmt.Errorf("add config users: %w", err)
	}

	// start service
	if err := s.ctl.Start(ctx); err != nil {
		return fmt.Errorf("restart service: %w", err)
	}

	// ping service after restart, several attempt required
	restartPingRetries := 5
	if err := s.ping(ctx, restartPingRetries); err != nil {
		return fmt.Errorf("ping service: %w", err)
	}

	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ctl.Stop(ctx); err != nil {
		return fmt.Errorf("stop xray service: %w", err)
	}

	return nil
}

func (s *Service) Status(ctx context.Context) (models.Status, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, err := s.ctl.Status(ctx)
	if err != nil {
		return models.NotRunning, fmt.Errorf("get xray status: %w", err)
	}

	return status, nil
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
