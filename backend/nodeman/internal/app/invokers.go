package app

/*var HelloMessageInvoker = fx.Invoke(
	func(cfg *config.Config, log *zap.Logger) {
		log.Warn(fmt.Sprintf("api available on %s via %s",
			cfg.ApiServicePath, cfg.ApiServiceUrl))
		log.Warn(fmt.Sprintf("user page available on %s via %s",
			cfg.UserSpaPath, cfg.UserSpaUrl))
		log.Warn(fmt.Sprintf("admin page available on %s via %s",
			cfg.AdminSpaPath, cfg.AdminSpaUrl))
	},
)

var HttpServerJob = fx.Invoke(
	func(s *server.HttpServer, lc fx.Lifecycle) {
		lc.AppendJob(fx.Job{
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

var BackgroundSyncJob = fx.Invoke(
	func(s *syncman.SyncMan, lc fx.Lifecycle) {
		lc.AppendJob(fx.Job{
			Name: "background sync",
			OnStart: func(context.Context) error {
				return s.Run()
			},
			OnStop: func(context.Context) error {
				return s.Stop()
			},
		})
	},
)

var BackgroundStatsJob = fx.Invoke(
	func(s *statsman.StatsMan, lc fx.Lifecycle) {
		lc.AppendJob(fx.Job{
			Name: "background stats",
			OnStart: func(context.Context) error {
				return s.Run()
			},
			OnStop: func(context.Context) error {
				return s.Stop()
			},
		})
	},
)

/*
// background stats
fx.core.AddRunner("background stats",
	func() (err error) {
		return statsJob.Run()
	},
	func(context.Context) error {
		return statsJob.Stop()
	},
)*/
