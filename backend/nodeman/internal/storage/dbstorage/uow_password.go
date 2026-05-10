package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

const adminID = 0

// auth.UoWContext impl
func (uow *uowctx) GetAuth(ctx context.Context) (*models.Auth, error) {
	query := queryReplacer.Replace(`
		SELECT
			{admin_id},
			{password_hash}
		FROM {admin_auth}
		WHERE {admin_id} = $1
		  AND {deleted_at} IS NULL
	`)

	row := uow.tx.QueryRowContext(ctx, query, adminID)

	var authId int
	var auth models.Auth
	err := row.Scan(
		&authId,
		&auth.PasswordHash,
	)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return &auth, nil
}

func (uow *uowctx) SetAuth(ctx context.Context, a models.Auth) error {
	query := queryReplacer.Replace(`
		INSERT INTO {admin_auth} (
			{admin_id},
			{password_hash},
			{updated_at}
		) VALUES ($1, $2, now())
		ON CONFLICT ({admin_id})
		DO UPDATE
		SET
			{password_hash} = EXCLUDED.{password_hash},
			{updated_at} = now()
	`)

	_, err := uow.tx.ExecContext(ctx, query, adminID, a.PasswordHash)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}
