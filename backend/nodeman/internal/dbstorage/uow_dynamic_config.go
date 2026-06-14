package dbstorage

import (
	"context"
	"encoding/json"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (uow *uowctx) EnsureDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	err = uow.q.EnsureDynamicConfig(ctx, raw)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) GetDynamicConfig(ctx context.Context) (*models.DynamicConfig, error) {
	raw, err := uow.q.GetDynamicConfig(ctx)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	var cfg models.DynamicConfig
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return &cfg, nil
}

func (uow *uowctx) SetDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	err = uow.q.SetDynamicConfig(ctx, raw)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}
