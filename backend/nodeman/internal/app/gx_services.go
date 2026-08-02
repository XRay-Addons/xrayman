package app

import (
	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/settings"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
)

var Services = gx.Module("services",
	gx.ProvideAnnotated(
		nodes.New,
		gx.ParamTags(``, ``, `name:"service-sync-timeout"`, ``),
		gx.As(new(handler.NodesService)),
	),
	gx.ProvideAnnotated(
		users.New,
		gx.ParamTags(``, ``, `name:"service-sync-timeout"`, ``),
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
	),
	gx.ProvideAnnotated(
		auth.New,
		gx.As(new(handler.AuthService)),
	),
)
