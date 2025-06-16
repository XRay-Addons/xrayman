package app

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/http/handlers"
	"github.com/XRay-Addons/xrayman/node/internal/service"
)

type XRayCfg interface {
	service.XRayCfg
	GetApiURL() string
}

type XRayCtl interface {
	service.XRayCtl
	Close(ctx context.Context) error
}

type XRayApi interface {
	service.XRayApi
	Close() error
}

type PerfCtl interface {
	service.PerfCtl
}

type Service interface {
	handlers.Service
}
