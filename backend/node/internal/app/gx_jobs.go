package app

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/http/server"
)

var httpServerJob = gx.Invoke(
	func(s *server.HttpServer) gx.Job {
		return gx.Job{
			Name: "http server",
			Run: func() error {
				return s.Listen()
			},
			Shutdown: func(ctx context.Context) error {
				return s.Shutdown(ctx)
			},
		}
	},
)

var Jobs = gx.Module("jobs",
	httpServerJob,
)
