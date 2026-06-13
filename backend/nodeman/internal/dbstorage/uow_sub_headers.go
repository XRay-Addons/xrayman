package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (uow *uowctx) NewSubHeader(ctx context.Context,
	header *models.Header,
) error {
	resp, err := uow.q.NewSubHeader(ctx, queries.NewSubHeaderParams{
		HeaderKey:   header.Key,
		HeaderValue: header.Value,
	})
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	header.ID = models.HeaderID(resp)

	return nil
}

func (uow *uowctx) DeleteSubHeader(ctx context.Context,
	id models.HeaderID,
) error {
	if err := uow.q.DeleteSubHeader(ctx, int64(id)); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

func (uow *uowctx) ListSubHeaders(ctx context.Context) (
	[]models.Header, error,
) {
	resp, err := uow.q.ListSubHeaders(ctx)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	headers, err := ConvertArray[queries.ListSubHeadersRow, models.Header](resp,
		With(func(from *queries.ListSubHeadersRow, to *models.Header) {
			to.ID = models.HeaderID(from.HeaderID)
			to.Key = from.HeaderKey
			to.Value = from.HeaderValue
		}),
	)
	if err != nil {
		return nil, err
	}

	return headers, nil
}
