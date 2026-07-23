package dbstorage

import (
	"context"
	"encoding/json"

	"github.com/XRay-Addons/xrayman/common/xerr"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) EnsureSettings(ctx context.Context, settings models.Settings) error {
	// pre-convert
	raw, err := json.Marshal(settings)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	// request
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.EnsureSettings(ctx, raw)
	})
}

func (s *Storage) GetSettings(ctx context.Context) (*models.Settings, error) {
	// request
	raw, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) (json.RawMessage, error) {
		return q.GetSettings(ctx)
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	var settings models.Settings
	err = json.Unmarshal(raw, &settings)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	return &settings, nil
}

func (s *Storage) SetSettings(ctx context.Context, settings models.Settings) error {
	// pre-convert
	raw, err := json.Marshal(settings)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	// request
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.SetSettings(ctx, raw)
	})
}
