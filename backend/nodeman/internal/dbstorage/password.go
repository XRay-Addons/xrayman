package dbstorage

import (
	"context"

	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

const adminID = 0

func (s *Storage) GetAuth(ctx context.Context) (
	*models.Auth, error,
) {
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) (queries.GetPasswordRow, error) {
		return q.GetPassword(ctx, adminID)
	})
	if err != nil {
		return nil, err
	}

	return &models.Auth{
		PasswordHash: resp.PasswordHash,
	}, nil
}

func (s *Storage) SetAuth(ctx context.Context, a *models.Auth) error {
	req := queries.SetPasswordParams{
		AdminID:      adminID,
		PasswordHash: a.PasswordHash,
	}
	return doVoid(ctx, s, func(ctx context.Context, q *queries.Queries) error {
		return q.SetPassword(ctx, req)
	})
}
