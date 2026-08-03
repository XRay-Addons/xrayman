package app

import (
	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"go.uber.org/zap"
)

var startupLog = gx.Invoke(
	func(k models.AccessKey, log *zap.Logger) {
		log.Warn("node access", zap.String("key", k.String()))
		log.Warn("node access", zap.String("key", k.String()))
		log.Warn("node access", zap.String("key", k.String()))
	},
)

var Startup = gx.Module("startup",
	startupLog,
)
