package app

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/servercfg"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayapi"
	xs "github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayservice"
	"github.com/XRay-Addons/xrayman/node/internal/service"
	"go.uber.org/zap"
)

type XRayServiceParams struct {
	gx.In
	Lc      gx.Lifecycle
	DataDir string `name:"xray-data-dir"`
	Log     *zap.Logger
}

var xrayService = gx.ProvideAnnotated(
	func(p XRayServiceParams) (*xs.XRayService, error) {
		s, err := xs.New(p.DataDir, xs.WithLogger(p.Log))
		if err != nil {
			return nil, err
		}
		p.Lc.AppendCloser(gx.Closer{
			Name: "xray service",
			OnClose: func(ctx context.Context) error {
				return s.Close(ctx)
			},
		})
		return s, err
	},
	gx.As(new(service.XRayService)),
)

type XRayApiParams struct {
	gx.In
	ServerCfg *servercfg.Config
	Lc        gx.Lifecycle
	Log       *zap.Logger
}

var xrayApi = gx.ProvideAnnotated(
	func(p XRayApiParams) (*xrayapi.XRayApi, error) {
		sc := p.ServerCfg
		xrayAPI, err := xrayapi.New(sc.GetApiURL(),
			sc.GetInbounds(), xrayapi.WithLogger(p.Log))
		if err != nil {
			return nil, err
		}
		p.Lc.AppendCloser(gx.Closer{
			Name: "xray api",
			OnClose: func(ctx context.Context) error {
				return xrayAPI.Close(ctx)
			},
		})
		return xrayAPI, nil
	},
	gx.As(new(service.XRayAPI)),
	gx.As(gx.Self()),
)

var XRay = gx.Module("xray",
	xrayService,
	xrayApi,
)
