package app

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/infra/xray/xrayapi"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/internal/service"
	"go.uber.org/zap"
)

type CheckXRayParams struct {
	gx.In
	Service *service.Service
	XRayApi *xrayapi.XRayApi
	Log     *zap.Logger
}

var checkXRay = gx.Invoke(
	func(lc gx.Lifecycle, p CheckXRayParams) {
		lc.AppendBootstrap(gx.Bootstrap{
			Name: "check xray",
			Fn: func(ctx context.Context) (err error) {
				// start service
				if _, err = p.Service.Start(ctx, models.StartParams{}); err != nil {
					return
				}
				defer func() {
					closeErr := p.Service.Stop(ctx)
					err = xerr.Join(err, closeErr)
				}()

				// ping xray
				if err = p.XRayApi.Ping(ctx); err != nil {
					return err
				}
				p.Log.Info("service ping ok!")
				return nil

			},
			Retry: func(err error) bool {
				return false
			},
		})
	},
)

var Bootstrap = gx.Module("bootstrap",
	checkXRay,
)
