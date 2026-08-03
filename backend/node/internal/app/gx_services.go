package app

import (
	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/node/internal/service"
	"go.uber.org/fx"
)

var Services = gx.Module("service",
	fx.Provide(service.New),
)
