package app

import (
	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/node/internal/infra/performance"
	"github.com/XRay-Addons/xrayman/node/internal/service"
)

var perf = gx.ProvideAnnotated(
	performance.New,
	gx.As(new(service.Performance)),
)

var Performance = gx.Module("performance",
	perf,
)
