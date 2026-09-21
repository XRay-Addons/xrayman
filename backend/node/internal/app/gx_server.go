package app

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/http/router"
	"github.com/XRay-Addons/xrayman/common/http/server"
	"github.com/XRay-Addons/xrayman/node/internal/http/api"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler/ogenserver"
	"github.com/XRay-Addons/xrayman/node/internal/http/security"
	"github.com/XRay-Addons/xrayman/node/internal/service"

	"go.uber.org/zap"
)

var httpHandler = gx.ProvideAnnotated(
	func(s *service.Service, l *zap.Logger) (*handler.Handler, error) {
		return handler.New(s, handler.WithLogger(l))
	},
	gx.As(new(ogenserver.Handler)),
)

var securityHandler = gx.ProvideAnnotated(
	security.New,
	gx.As(new(ogenserver.SecurityHandler)),
)

var apiHandler = gx.ProvideAnnotated(
	api.NewHandler,
	gx.ResultTags(`name:"ogenserver-handler"`),
)

type HttpRouterParams struct {
	gx.In
	ApiHandler http.Handler `name:"ogenserver-handler"`
	Log        *zap.Logger
}

var httpRouter = gx.ProvideNamed(
	func(p HttpRouterParams) (http.Handler, error) {
		return router.New(
			router.WithHandler("/", p.ApiHandler),
			router.WithLogger(p.Log))
	},
	"http-router",
)

type ServerParams struct {
	gx.In
	Endpoint string       `name:"endpoint"`
	Router   http.Handler `name:"http-router"`
	TLS      *tls.Config
	Log      *zap.Logger
}

var httpServer = gx.ProvideNamed(
	func(p ServerParams) (*server.HttpServer, error) {
		return server.New(p.Endpoint, p.Router,
			server.WithTLS(p.TLS), server.WithLog(p.Log))
	},
	"http-server",
)

type HttpServerJobParams struct {
	gx.In
	S *server.HttpServer `name:"http-server"`
}

var httpServerJob = gx.Invoke(
	func(p HttpServerJobParams, lc gx.Lifecycle) {
		lc.AppendJob(gx.Job{
			Name: "http server",
			Run: func() error {
				return p.S.Listen()
			},
			Shutdown: func(ctx context.Context) error {
				return p.S.Shutdown(ctx)
			},
		})
	},
)

var Server = gx.Module("server",
	httpHandler,
	securityHandler,
	apiHandler,
	httpRouter,
	httpServer,
	httpServerJob,
)
