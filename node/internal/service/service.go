package service

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type Service struct {
	serverCfg      ServerCfg
	clientCfg      ClientCfg
	xrayServiceCtl XRayServiceCtl
	xrayAPI        XRayAPI
}

func New(
	serverCfg ServerCfg,
	clientCfg ClientCfg,
	xrayServiceCtl XRayServiceCtl,
	xrayAPI XRayAPI,
) (*Service, error) {
	if serverCfg == nil {
		return nil, fmt.Errorf("%w: service init: serverCfg", errdefs.ErrNilArgPassed)
	}
	if clientCfg == nil {
		return nil, fmt.Errorf("%w: service init: clientCfg", errdefs.ErrNilArgPassed)
	}
	if xrayServiceCtl == nil {
		return nil, fmt.Errorf("%w: service init: xrayServiceCtl", errdefs.ErrNilArgPassed)
	}
	if xrayAPI == nil {
		return nil, fmt.Errorf("%w: service init: xrayAPI", errdefs.ErrNilArgPassed)
	}

	return &Service{
		serverCfg:      serverCfg,
		clientCfg:      clientCfg,
		xrayServiceCtl: xrayServiceCtl,
		xrayAPI:        xrayAPI,
	}, nil
}

func (s *Service) Start(
	ctx context.Context,
	users []models.User,
) (*models.NodeProperties, error) {
	if s == nil {
		return nil, fmt.Errorf("%w: service: start", errdefs.ErrNilObjectCall)
	}
	// get server config
	cfg, err := s.serverCfg.GetUsersCfg(users)
	if err != nil {
		return nil, fmt.Errorf("service start: %w", err)
	}
	// start server
	if err = s.xrayServiceCtl.Start(ctx, cfg); err != nil {
		return nil, fmt.Errorf("service start: %w", err)
	}
	// get server properties
	clientTemplate, err := s.clientCfg.GetClientCfgTemplate()
	if err != nil {
		return nil, fmt.Errorf("service start: %w", err)
	}
	// return node properties
	return &models.NodeProperties{
		ClientCfgTemplate: *clientTemplate,
	}, nil
}

func (s *Service) Stop(
	ctx context.Context,
) error {
	if s == nil {
		return fmt.Errorf("%w: service: start", errdefs.ErrNilObjectCall)
	}
	if err := s.xrayServiceCtl.Stop(ctx); err != nil {
		return fmt.Errorf("service stop: %w", err)
	}
	return nil
}

func (s *Service) Status(
	ctx context.Context,
) (*models.NodeStatus, error) {
	if s == nil {
		return nil, fmt.Errorf("%w: service: start", errdefs.ErrNilObjectCall)
	}
	status, err := s.xrayServiceCtl.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("service status: %w", err)
	}
	return &models.NodeStatus{
		ServiceStatus: status,
	}, nil
}

func (s *Service) EditUsers(
	ctx context.Context,
	add, remove []models.User,
) error {
	if s == nil {
		return fmt.Errorf("%w: service: start", errdefs.ErrNilObjectCall)
	}
	if err := s.xrayAPI.EditUsers(ctx, add, remove); err != nil {
		return fmt.Errorf("service: edit users: %w", err)
	}
	return nil
}
