package dbstorage

import (
	"context"
	"encoding/json"

	"github.com/XRay-Addons/xrayman/common/xerr"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) EnsureDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error {
	// pre-convert
	raw, err := json.Marshal(cfg)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	// request
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.EnsureDynamicConfig(ctx, raw)
	})
}

func (s *Storage) GetDynamicConfig(ctx context.Context) (*models.DynamicConfig, error) {
	// request
	raw, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) (json.RawMessage, error) {
		return q.GetDynamicConfig(ctx)
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	var cfg models.DynamicConfig
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return &cfg, nil
}

func (s *Storage) SetDynamicConfig(ctx context.Context, cfg models.DynamicConfig) error {
	// pre-convert
	raw, err := json.Marshal(cfg)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	// request
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.SetDynamicConfig(ctx, raw)
	})
}
