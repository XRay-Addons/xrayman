package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

const adminID = 0

func (uow *uowctx) GetAuth(ctx context.Context) (
	*models.Auth, error,
) {
	resp, err := uow.q.GetPassword(ctx, adminID)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	auth := models.Auth{
		PasswordHash: resp.PasswordHash,
	}

	return &auth, nil
}

func (uow *uowctx) SetAuth(ctx context.Context, a models.Auth) error {
	if err := uow.q.SetPassword(ctx, queries.SetPasswordParams{
		AdminID:      adminID,
		PasswordHash: a.PasswordHash,
	}); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}
