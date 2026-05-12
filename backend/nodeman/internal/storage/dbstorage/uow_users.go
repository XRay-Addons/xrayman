package dbstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

// service.UoWContext::UsersStorage impl
func (uow *uowctx) NewUser(ctx context.Context, user *models.User) error {
	query := queryReplacer.Replace(`
		INSERT INTO {users} (
			{display_name},
			{user_name},
			{vless_uuid},
			{user_target_status}
		) VALUES ($1, $2, $3, $4)
		RETURNING {user_id}
	`)

	err := uow.tx.QueryRowContext(ctx, query,
		user.Profile.DisplayName,
		user.Profile.Name,
		user.Profile.VlessUUID,
		user.TargetStatus,
	).Scan(&user.Profile.ID)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) GetUser(ctx context.Context, id models.UserID) (*models.User, bool, error) {
	query := queryReplacer.Replace(`
		SELECT
			{user_id},
			{display_name},
			{user_name},
			{vless_uuid},
			{user_target_status}
		FROM {users}
		WHERE {user_id} = $1
		  AND {deleted_at} IS NULL
	`)

	row := uow.tx.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.Profile.ID,
		&user.Profile.DisplayName,
		&user.Profile.Name,
		&user.Profile.VlessUUID,
		&user.TargetStatus,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, xerr.WrapWithStack(err)
	}

	return &user, true, nil
}

func (uow *uowctx) ListUsers(ctx context.Context) (users []models.User, err error) {
	query := queryReplacer.Replace(`
		SELECT
			{user_id},
			{display_name},
			{user_name},
			{vless_uuid},
			{user_target_status}
		FROM {users}
		WHERE {deleted_at} IS NULL
		ORDER BY {user_id} ASC
	`)

	rows, err := uow.tx.QueryContext(ctx, query)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	defer func() {
		err = xerr.Join(err, xerr.WrapWithStack(rows.Close()))
	}()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Profile.ID,
			&user.Profile.DisplayName,
			&user.Profile.Name,
			&user.Profile.VlessUUID,
			&user.TargetStatus,
		)
		if err != nil {
			return nil, xerr.WrapWithStack(err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return users, nil
}

func (uow *uowctx) SetTargetUserStatus(ctx context.Context, id models.UserID, status models.UserStatus) error {
	query := queryReplacer.Replace(`
		UPDATE {users}
		SET
			{user_target_status} = $1,
			{updated_at} = now()
		WHERE {user_id} = $2
		  AND {deleted_at} IS NULL
	`)

	_, err := uow.tx.ExecContext(ctx, query, status, id)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) DeleteUser(ctx context.Context, id models.UserID) error {
	query := queryReplacer.Replace(`
		UPDATE {users}
		SET {deleted_at} = now()
		WHERE {user_id} = $1
		  AND {deleted_at} IS NULL
	`)

	_, err := uow.tx.ExecContext(ctx, query, id)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}
