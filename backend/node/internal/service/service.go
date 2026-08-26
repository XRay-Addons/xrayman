package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/internal/version"
)

type Service struct {
	serverCfg   ServerCfg
	clientCfg   ClientConfig
	xrayService XRayService
	xrayAPI     XRayAPI
	perf        Performance
}

func New(
	serverCfg ServerCfg,
	clientCfg ClientConfig,
	xrayService XRayService,
	xrayAPI XRayAPI,
	perf Performance,
) (*Service, error) {
	if serverCfg == nil {
		return nil, errdefs.NilArg("serverCfg")
	}
	if clientCfg == nil {
		return nil, errdefs.NilArg("clientCfg")
	}
	if xrayService == nil {
		return nil, errdefs.NilArg("xrayService")
	}
	if xrayAPI == nil {
		return nil, errdefs.NilArg("xrayAPI")
	}
	if perf == nil {
		return nil, errdefs.NilArg("perf")
	}

	return &Service{
		serverCfg:   serverCfg,
		clientCfg:   clientCfg,
		xrayService: xrayService,
		xrayAPI:     xrayAPI,
		perf:        perf,
	}, nil
}

func (s *Service) Start(ctx context.Context,
	params models.StartParams,
) (*models.StartResult, error) {
	if s == nil {
		return nil, errdefs.NilCall()
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
	// get server properties
	clientCfg, err := s.clientCfg.GetTemplate()
	if err != nil {
		return nil, err
	}
	// return node properties
	return &models.StartResult{
		ClientConfigTemplate: *clientCfg,
		Version:              version.Version,
	}, nil
}

func (s *Service) Stop(ctx context.Context) error {
	// stop server
	if err := s.xrayService.Stop(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) Status(
	ctx context.Context,
) (*models.StatusResult, error) {
	status, err := s.xrayService.Status(ctx)
	if err != nil {
		return nil, err
	}

	return &models.StatusResult{ServiceStatus: status}, nil
}

func (s *Service) EditUsers(ctx context.Context,
	params models.EditUsersParams,
) error {
	if err := s.xrayAPI.EditUsers(ctx, params.Add, params.Remove); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetStats(ctx context.Context) (*models.StatsResult, error) {
	usersStats, err := s.xrayAPI.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	perf, err := s.perf.GetPerformance(ctx)
	if err != nil {
		return nil, err
	}
	return &models.StatsResult{
		Users:       usersStats,
		Performance: *perf,
	}, nil
}
