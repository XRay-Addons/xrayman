package xrayapi

import (
	"context"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	statsService "github.com/xtls/xray-core/app/stats/command"
)

func Ping(ctx context.Context, ssClient statsService.StatsServiceClient) error {
	if ssClient == nil {
		return fmt.Errorf("%w: stats service not exists", errdefs.ErrIPE)
	}

	_, err := ssClient.GetSysStats(ctx, &statsService.SysStatsRequest{})
	if err != nil {
		return fmt.Errorf("%w: ping call: %v", errdefs.ErrXRay, err)
	}

	return nil
}
