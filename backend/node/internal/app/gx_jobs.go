package app

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/http/server"
)

var httpServerJob = gx.Invoke(
	func(s *server.HttpServer, lc gx.Lifecycle) {
		lc.AppendJob(gx.Job{
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

var Jobs = gx.Module("jobs",
	httpServerJob,
)
