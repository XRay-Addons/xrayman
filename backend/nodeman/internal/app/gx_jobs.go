package app

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/statsman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/syncman"
	"go.uber.org/zap"
)

var backgroundSyncJob = gx.Options(
	gx.Provide(
		func(ps syncman.PoolSyncer, cfg *config.Config, l *zap.Logger) (*syncman.SyncMan, error) {
			return syncman.New(ps, cfg.StateSyncInterval, syncman.WithLogger(l))
		},
	),
	gx.Invoke(
		func(s *syncman.SyncMan, lc gx.Lifecycle) {
			lc.AppendJob(gx.Job{
				Name: "background sync",
				OnStart: func(context.Context) error {
					return s.Run()
				},
				OnStop: func(context.Context) error {
					s.Stop()
					return nil
				},
			})
		},
	),
)

var backgroundStatsJob = gx.Options(
	gx.Provide(
		func(ps statsman.StatsUpdater, cfg *config.Config, l *zap.Logger) (*statsman.StatsMan, error) {
			return statsman.New(ps, cfg.StatsSyncInterval, statsman.WithLogger(l))
		},
	),
	gx.Invoke(
		func(s *statsman.StatsMan, lc gx.Lifecycle) {
			lc.AppendJob(gx.Job{
				Name: "background stats",
				OnStart: func(context.Context) error {
					return s.Run()
				},
				OnStop: func(context.Context) error {
					s.Stop()
					return nil
				},
			})
		},
	),
)

var Jobs = gx.Module("jobs",
	backgroundSyncJob,
	backgroundStatsJob,
)
