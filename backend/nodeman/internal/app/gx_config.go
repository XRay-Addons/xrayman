package app

import (
	"time"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
)

var Config = gx.Module("config", gx.Provide(
	config.Init,
	gx.Named(
		func(cfg *config.Config) time.Duration {
			return cfg.NodeCallTimeout
		},
		"node-call-timeout",
	),
	gx.Named(
		func(cfg *config.Config) time.Duration {
			return cfg.NodeCallTimeout + cfg.StorageCallTimeout
		},
		"service-sync-timeout",
	),
	gx.Named(
		func(cfg *config.Config) string {
			return cfg.JwtSecret
		},
		"jwt-secret",
	),
	gx.Named(
		func(cfg *config.Config) string {
			return cfg.AdminPassword
		},
		"admin-password",
	),
	gx.Named(
		func() string {
			return JWTIssuer
		},
		"jwt-issuer",
	),
))
