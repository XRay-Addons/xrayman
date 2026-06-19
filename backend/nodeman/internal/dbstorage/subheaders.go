package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/convert"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) NewSubHeader(ctx context.Context,
	header *models.Header,
) error {
	// pre-convert
	req := queries.NewSubHeaderParams{
		HeaderKey:   header.Key,
		HeaderValue: header.Value,
	}

	// request
	resp, err := doAny(ctx, s, func(ctx context.Context, q *queries.Queries) (int64, error) {
		return q.NewSubHeader(ctx, req)
	})
	if err != nil {
		return err
	}

	header.ID = models.HeaderID(resp)
	return nil
}

func (s *Storage) DeleteSubHeader(ctx context.Context,
	id models.HeaderID,
) error {
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.DeleteSubHeader(ctx, int64(id))
	})
}

func (s *Storage) ListSubHeaders(ctx context.Context) (
	[]models.Header, error,
) {
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries) ([]queries.ListSubHeadersRow, error,
	) {
		return q.ListSubHeaders(ctx)
	})
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return convert.ListSubHeadersResp(resp), nil
}
