package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
)

const adminID = 0

// auth.UoWContext impl
func (uow *uowctx) GetAdmin(ctx context.Context) (*auth.Auth, error) {
	query := queryReplacer.Replace(`
		SELECT
			{admin_id},
			{password_hash},
		FROM {admin_auth}
		WHERE {admin_id} = $1
		  AND {deleted_at} IS NULL
	`)

	row := uow.tx.QueryRowContext(ctx, query, adminID)

	var authId int
	var auth auth.Auth
	err := row.Scan(
		&authId,
		&auth.PasswordHash,
	)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	return &auth, nil
}

func (uow *uowctx) SetAdmin(ctx context.Context, a *auth.Auth) error {
	query := queryReplacer.Replace(`
		INSERT INTO {admin_auth} (
			{admin_id},
			{password_hash},
			{updated_at}
		) VALUES ($1, $2, now())
		ON CONFLICT ({id})
		DO UPDATE
		SET
			{password_hash} = EXCLUDED.{password_hash},
			{updated_at} = now()
	`)

	_, err := uow.tx.ExecContext(ctx, query, adminID, a.PasswordHash)
	if err != nil {
		return errdefs.WrapWithStack(err)
	}

	return nil
}
