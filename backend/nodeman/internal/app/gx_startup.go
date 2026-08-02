package app

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/config"
	"go.uber.org/zap"
)

var startupLog = gx.Invoke(
	func(cfg *config.Config, log *zap.Logger) {
		log.Warn(fmt.Sprintf("api available on %s via %s",
			cfg.ApiServicePath, cfg.ApiServiceUrl))
		log.Warn(fmt.Sprintf("user page available on %s via %s",
			cfg.UserSpaPath, cfg.UserSpaUrl))
		log.Warn(fmt.Sprintf("admin page available on %s via %s",
			cfg.AdminSpaPath, cfg.AdminSpaUrl))
	},
)

var Startup = gx.Module("startup",
	startupLog,
)
