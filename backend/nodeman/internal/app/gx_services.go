package app

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/settings"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/version"
	"go.uber.org/zap"
)

type NodesServiceParams struct {
	gx.In
	Lc          gx.Lifecycle
	PoolSyncer  nodes.Syncer
	Storage     nodes.Storage
	SyncTimeout time.Duration `name:"service-sync-timeout"`
	Log         *zap.Logger
}

type UsersServiceParams struct {
	gx.In
	Lc          gx.Lifecycle
	PoolSyncer  users.Syncer
	Storage     users.Storage
	SyncTimeout time.Duration `name:"service-sync-timeout"`
	Log         *zap.Logger
}

var Services = gx.Module("services",
	gx.ProvideAnnotated(
		func(p NodesServiceParams) (*nodes.Service, error) {
			ns, err := nodes.New(p.PoolSyncer, p.Storage, p.SyncTimeout, p.Log)
			if err != nil {
				return nil, err
			}
			p.Lc.AppendCloser(gx.Closer{
				Name: "NodesService",
				OnClose: func(context.Context) error {
					ns.Close()
					return nil
				},
			})
			return ns, nil
		},
		gx.As(new(handler.NodesService)),
	),
	gx.ProvideAnnotated(
		func(p UsersServiceParams) (*users.Service, error) {
			us, err := users.New(p.PoolSyncer, p.Storage, p.SyncTimeout, p.Log)
			if err != nil {
				return nil, err
			}
			p.Lc.AppendCloser(gx.Closer{
				Name: "UsersService",
				OnClose: func(context.Context) error {
					us.Close()
					return nil
				},
			})
			return us, nil
		},
		gx.As(new(handler.UsersService)),
	),
	gx.ProvideAnnotated(
		subscr.New,
		gx.As(new(handler.SubscrService)),
		gx.As(gx.Self()),
	),
	gx.ProvideAnnotated(
		settings.New,
		gx.As(new(handler.SettingsService)),
		gx.As(gx.Self()),
	),
	gx.ProvideAnnotated(
		auth.New,
		gx.As(new(handler.AuthService)),
		gx.As(gx.Self()),
	),
	gx.ProvideAnnotated(
		version.New,
		gx.As(new(handler.VersionService)),
	),
)
