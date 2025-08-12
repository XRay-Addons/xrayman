package app

import (
	"context"
	"fmt"
	"net/http"

	client "github.com/XRay-Addons/xrayman/nodeman/internal/clients/node"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/router"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/server"
	a "github.com/XRay-Addons/xrayman/nodeman/internal/infra/app"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/httpclient"
	"github.com/XRay-Addons/xrayman/nodeman/internal/node"
	"github.com/XRay-Addons/xrayman/nodeman/internal/pool"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/memstorage"
	"github.com/XRay-Addons/xrayman/nodeman/internal/syncman"

	"go.uber.org/zap"
)

type App struct {
	app *a.App
}

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	if log == nil {
		return nil, fmt.Errorf("%w: app init: logger", errdefs.ErrNilArgPassed)
	}

	var storage *memstorage.Storage

	var httpClient *httpclient.ClientFactory
	var poolClient *client.PoolClient

	var nodeSyncer *node.Syncer
	var poolSyncer *pool.Syncer
	var syncMan *syncman.Manager

	var s *service.Service

	var h *handler.Handler
	var r http.Handler
	var httpServer *server.HttpServer

	app := a.New(
		// http client
		a.WithComponent("http client",
			func() (err error) {
				httpClient = httpclient.NewClientFactory()
				return
			}, nil,
			//func(ctx context.Context) error {
			//	httpClient.CloseIdleConnections()
			//	return nil
			//},
		),
		// pool client
		a.WithComponent("pool client",
			func() (err error) {
				poolClient, err = client.NewPoolClient(client.WithHTTPClient(httpClient))
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

		// node syncer
		a.WithComponent("node syncer",
			func() error {
				nodeSyncer = node.NewSyncer()
				return nil
			}, nil,
		),
		// pool syncer
		a.WithComponent("pool syncer",
			func() (err error) {
				poolSyncer, err = pool.NewSyncer(storage.PoolUoW(), poolClient, nodeSyncer)
				return
			}, nil,
		),
		// sync manager
		a.WithComponent("sync man",
			func() (err error) {
				syncMan, err = syncman.New(poolSyncer, syncman.WithLog(log))
				return
			},
			func(ctx context.Context) error {
				if err := syncMan.Close(); err != nil {
					return fmt.Errorf("app close: sync service: %w", err)
				}
				return nil
			},
		),

		// service
		a.WithComponent("service",
			func() (err error) {
				s, err = service.New(syncMan, storage.ServiceUoW())
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
