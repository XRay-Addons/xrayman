package app

import (
	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/node/internal/config"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/clientcfg"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/servercfg"
	"github.com/XRay-Addons/xrayman/node/internal/service"
)

var Config = gx.Module("config",
	gx.ProvideNamed(
		func(cfg *config.Config) string {
			return cfg.PersistentDir
		},
		"persistent-dir",
	),
	gx.ProvideNamed(
		func(cfg *config.Config) string {
			return cfg.XRayServer()
		},
		"xray-server",
	),
	gx.ProvideNamed(
		func(cfg *config.Config) string {
			return cfg.XRayClient()
		},
		"xray-client",
	),
	gx.ProvideNamed(
		func(cfg *config.Config) string {
			return cfg.XRayDataDir
		},
		"xray-data-dir",
	),
	gx.ProvideNamed(
		func(cfg *config.Config) string {
			return cfg.Endpoint
		},
		"endpoint",
	),
	gx.ProvideAnnotated(
		servercfg.New,
		gx.As(new(service.ServerCfg)),
		gx.As(gx.Self()),
		gx.ParamTags(`name:"xray-server"`),
	),
	gx.ProvideAnnotated(
		clientcfg.New,
		gx.As(new(service.ClientConfig)),
		gx.ParamTags(`name:"xray-client"`),
	),
)
