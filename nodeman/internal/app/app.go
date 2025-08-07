package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/clients/node/security"
	client "github.com/XRay-Addons/xrayman/nodeman/internal/clients/node"
	"github.com/XRay-Addons/xrayman/nodeman/internal/clients/node/security"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/router"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/server"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/tlscfg"
	a "github.com/XRay-Addons/xrayman/nodeman/internal/infra/app"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/keygen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodesyncer"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/poolmonitor"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/poolsyncer"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/service"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/memstorage"

	"go.uber.org/zap"
)

type App struct {
	app *a.App
}

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: app init: logger", errdefs.ErrNilArgPassed)
	}

	//var srvCfg *xraycfg.ServerCfg
	//var clientCfg *xraycfg.ClientCfg
	var tlsCfg *tls.Config
	//var xrayService *xrayservice.XRayService
	//var xrayAPI *xrayapi.XRayApi

	var storage *memstorage.Storage
	var kg *keygen.Keygen

	var httpClient *http.Client
	var poolClient client.PoolClient
	//var sec *security.SecurityFactory
	var nodeSyncer *nodesyncer.NodesSyncer

	var s *service.Service
	var h *handler.Handler
	//var sec api.SecurityHandler
	var r http.Handler

	var httpServer *server.HttpServer

	app := a.New(
		/*// server config
		a.WithComponent("server cfg",
			func() (err error) {
				srvCfg, err = xraycfg.NewServerCfg(cfg.XRayServer())
				return
			}, nil,
		),
		// client config
		a.WithComponent("client cfg",
			func() (err error) {
				clientCfg, err = xraycfg.NewClientCfg(cfg.XRayClient())
				return
			}, nil,
		),*/
		// TLS config
		a.WithComponent("tls cfg",
			func() (err error) {
				if !cfg.HasCerts() {
					log.Warn("xray dir contains NO certs, encryption disabled. use it only for testing!!!")
					return
				}
				tlsCfg, err = tlscfg.Load(cfg.NodemanCrt(), cfg.NodemanKey(), cfg.RootCrt())
				return
			}, nil,
		),

		// storage
		a.WithComponent("storage",
			func() error {
				storage = memstorage.New()
				return nil
			}, nil,
		),

		// keygen
		a.WithComponent("keygen",
			func() error {
				kg = keygen.New()
				return nil
			}, nil,
		),

		// http client
		a.WithComponent("http client", 
			func() (err error) {
				httpClient, err = httpclient.New()
				return
			},
			func() (err error) {
				httpClient.CloseIdleConnections()
			},
		),
		// pool client
		a.WithComponent("pool client",
			func() (err error) {
				clientFactory, err = client.NewPoolClient(client.WithHTTPClient(httpClient))
				return
			}, nil,
		),
		// node syncer
		a.WithComponent("node syncer",
			func() (err error) {
				nodesyncer.New()
			}
	
	)
		// pool syncer
		a.WithComponent("pool syncer",
			func() (err error) {
				poolsyncer.New(clientFactory,  )
			}
	)
		// security
		a.WithComponent("nodes security",
			func() error {
				sec = security.NewFactory()
				return nil
			}, nil,
		),

		// poolmon
		a.WithComponent("poolmon",
			func() error {
				poolmonitor.New()
			}


		// service
		a.WithComponent("service",
			func() (err error) {
				s, err = service.New(storage, kg, xrayService, xrayAPI)
				return
			}, nil,
		),
		// handler
		a.WithComponent("handler",
			func() (err error) {
				h, err = handler.New(s, log)
				return
			}, nil,
		),
		// router
		a.WithComponent("router",
			func() (err error) {
				r, err = router.New(h, router.WithLogger(log))
				return
			}, nil,
		),

		// http server
		a.WithComponent("http server",
			func() (err error) {
				httpServer, err = server.New(cfg.Endpoint, r)
				return
			}, nil,
		),

		a.WithRunner("http server",
			func() (err error) {
				return httpServer.Listen()
			},
			func(ctx context.Context) error {
				return httpServer.Shutdown(ctx)
			},
		),

		// logger
		a.WithLogger(log),

		// cancel by Ctrl+C
		a.WithSignalCancel(),
	)

	return &App{
		app: app,
	}, nil
}

func (app *App) Run() error {
	if app == nil {
		return fmt.Errorf("%w: app: run", errdefs.ErrNilObjectCall)
	}
	return app.app.Run()
}
