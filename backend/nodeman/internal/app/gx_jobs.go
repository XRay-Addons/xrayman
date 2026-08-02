package app

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/http/server"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/statsman"
	"github.com/XRay-Addons/xrayman/nodeman/internal/jobs/syncman"
	"go.uber.org/zap"
)

var httpServerJob = gx.Invoke(
	func(s *server.HttpServer, lc gx.Lifecycle) {
		lc.AppendJobEx(gx.Job{
			Name: "http server",
			OnStart: func(context.Context) error {
				return s.Listen()
			},
			OnStop: func(ctx context.Context) error {
				return s.Shutdown(ctx)
			},
		})
	},
)

var backgroundSyncJob = gx.Options(
	gx.Provide(
		func(ps syncman.PoolSyncer, cfg *config.Config, l *zap.Logger) (*syncman.SyncMan, error) {
			return syncman.New(ps, cfg.StateSyncInterval, syncman.WithLogger(l))
		},
	),
	gx.Invoke(
		func(s *syncman.SyncMan, lc gx.Lifecycle) {
			lc.AppendJobEx(gx.Job{
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
			lc.AppendJobEx(gx.Job{
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
	httpServerJob,
	backgroundSyncJob,
	backgroundStatsJob,
)
