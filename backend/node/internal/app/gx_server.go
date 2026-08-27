package app

import (
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

type RouterParams struct {
	gx.In
	ApiHandler http.Handler `name:"ogenserver-handler"`
	Log        *zap.Logger
}

var r = gx.ProvideNamed(
	func(p RouterParams) (http.Handler, error) {
		return router.New(
			router.WithHandler("/", p.ApiHandler),
			router.WithLogger(p.Log))
	},
	"router",
)

type ServerParams struct {
	gx.In
	Endpoint string       `name:"endpoint"`
	Router   http.Handler `name:"router"`
	TLS      *tls.Config
}

var s = gx.Provide(
	func(p ServerParams) (*server.HttpServer, error) {
		return server.New(p.Endpoint, p.Router, server.WithTLS(p.TLS))
	},
)

var Server = gx.Module("server",
	httpHandler,
	securityHandler,
	apiHandler,
	r,
	s,
)
