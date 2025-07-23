package service

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type Service struct {
	serverCfg   ServerCfg
	clientCfg   ClientCfg
	xrayService XRayService
	xrayAPI     XRayAPI
}

func New(
	serverCfg ServerCfg,
	clientCfg ClientCfg,
	xrayService XRayService,
	xrayAPI XRayAPI,
) (*Service, error) {
	if serverCfg == nil {
		return nil, fmt.Errorf("%w: service init: serverCfg", errdefs.ErrNilArgPassed)
	}
	if clientCfg == nil {
		return nil, fmt.Errorf("%w: service init: clientCfg", errdefs.ErrNilArgPassed)
	}
	if xrayService == nil {
		return nil, fmt.Errorf("%w: service init: xrayService", errdefs.ErrNilArgPassed)
	}
	if xrayAPI == nil {
		return nil, fmt.Errorf("%w: service init: xrayAPI", errdefs.ErrNilArgPassed)
	}

	return &Service{
		serverCfg:   serverCfg,
		clientCfg:   clientCfg,
		xrayService: xrayService,
		xrayAPI:     xrayAPI,
	}, nil
}

func (s *Service) Start(ctx context.Context, params models.StartParams) (*models.StartResult, error) {
	if s == nil {
		return nil, fmt.Errorf("%w: service: start", errdefs.ErrNilObjectCall)
	}
	// get server config
	cfg, err := s.serverCfg.GetUsersCfg(params.Users)
	if err != nil {
		return nil, fmt.Errorf("service start: %w", err)
	}
	// start server
	if err = s.xrayService.Start(ctx, cfg); err != nil {
		return nil, fmt.Errorf("service start: %w", err)
	}
	// connect to server api
	if err = s.xrayAPI.Connect(ctx); err != nil {
		return nil, fmt.Errorf("service api connect: %w", err)
	}
	// get server properties
	clientCfg, err := s.clientCfg.Get()
	if err != nil {
		return nil, fmt.Errorf("service start: %w", err)
	}
	// return node properties
	return &models.StartResult{ClientCfg: *clientCfg}, nil
}

func (s *Service) Stop(ctx context.Context, params models.StopParams) (*models.StopResult, error) {
	if s == nil {
		return nil, fmt.Errorf("%w: service: start", errdefs.ErrNilObjectCall)
	}
	// disconnect from server api
	if err := s.xrayAPI.Disconnect(ctx); err != nil {
		return nil, fmt.Errorf("service api disconnect: %w", err)
	}
	if err := s.xrayService.Stop(ctx); err != nil {
		return nil, fmt.Errorf("service stop: %w", err)
	}
	return &models.StopResult{}, nil
}

func (s *Service) Status(ctx context.Context, params models.StatusParams) (*models.StatusResult, error) {
	if s == nil {
		return nil, fmt.Errorf("%w: service: start", errdefs.ErrNilObjectCall)
	}
	status, err := s.xrayService.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("service status: %w", err)
	}
	return &models.StatusResult{ServiceStatus: status}, nil
}

func (s *Service) EditUsers(ctx context.Context, params models.EditUsersParams) (*models.EditUsersResult, error) {
	if s == nil {
		return nil, fmt.Errorf("%w: service: start", errdefs.ErrNilObjectCall)
	}
	// don't check status, xrayService.Status is too slow on osx
	// TODO: add linux support, use it, check status before grpc api call
	if err := s.xrayAPI.EditUsers(ctx, params.Add, params.Remove); err != nil {
		return nil, fmt.Errorf("service: edit users: %w", err)
	}
	return &models.EditUsersResult{}, nil
}
