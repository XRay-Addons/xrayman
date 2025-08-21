package service

import (
	"context"

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
		return nil, errdefs.NewNilArg("serverCfg")
	}
	if clientCfg == nil {
		return nil, errdefs.NewNilArg("clientCfg")
	}
	if xrayService == nil {
		return nil, errdefs.NewNilArg("xrayService")
	}
	if xrayAPI == nil {
		return nil, errdefs.NewNilArg("xrayAPI")
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
		return nil, errdefs.NewNilCall()
	}
	// get server config
	cfg, err := s.serverCfg.GetUsersCfg(params.Users)
	if err != nil {
		return nil, err
	}
	// start server
	if err = s.xrayService.Start(ctx, cfg); err != nil {
		return nil, err
	}
	// connect to server api
	if err = s.xrayAPI.Connect(ctx); err != nil {
		return nil, err
	}
	// get server properties
	clientCfg, err := s.clientCfg.Get()
	if err != nil {
		return nil, err
	}
	// return node properties
	return &models.StartResult{ClientCfg: *clientCfg}, nil
}

func (s *Service) Stop(ctx context.Context, params models.StopParams) (*models.StopResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	// disconnect from server api
	if err := s.xrayAPI.Disconnect(ctx); err != nil {
		return nil, err
	}
	if err := s.xrayService.Stop(ctx); err != nil {
		return nil, err
	}
	return &models.StopResult{}, nil
}

func (s *Service) Status(ctx context.Context, params models.StatusParams) (*models.StatusResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	status, err := s.xrayService.Status(ctx)
	if err != nil {
		return nil, err
	}
	return &models.StatusResult{ServiceStatus: status}, nil
}

func (s *Service) EditUsers(ctx context.Context, params models.EditUsersParams) (*models.EditUsersResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	// don't check status, xrayService.Status is too slow on osx
	// TODO: add linux support, use it, check status before grpc api call
	if err := s.xrayAPI.EditUsers(ctx, params.Add, params.Remove); err != nil {
		return nil, err
	}
	return &models.EditUsersResult{}, nil
}
