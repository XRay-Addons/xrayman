package dbstorage

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/convert"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

const month = 30 * 24 * time.Hour

func (s *Storage) NewUser(ctx context.Context, user *models.User) error {
	// pre-convert
	arg := convert.NewUserReq(user)

	// request
	userID, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) (int64, error) {
		return q.NewUser(ctx, *arg)
	})
	if err != nil {
		return err
	}

	// post-convert
	user.Profile.ID = models.UserID(userID)

	return nil
}

func (s *Storage) GetUserView(ctx context.Context,
	id models.UserID, name string,
) (*models.UserView, error) {
	// pre-convert
	from := time.Now().Add(-month)
	req := queries.GetUserViewParams{
		FromDay:  from,
		UserID:   int64(id),
		UserName: name,
	}

	// request
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) (queries.GetUserViewRow, error) {
		return q.GetUserView(ctx, req)
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.GetUserViewResp(&resp), nil
}

func (s *Storage) ListUsers(ctx context.Context) (
	[]models.User, error,
) {
	// request
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) ([]queries.ListUsersRow, error) {
		return q.ListUsers(ctx)
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.ListUsersResp(resp), nil
}

func (s *Storage) ListUserViews(ctx context.Context) ([]models.UserView, error) {
	// request
	from := time.Now().Add(-month)
	resp, err := doAny(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) ([]queries.ListUserViewsRow, error) {
		return q.ListUserViews(ctx, from)
	})
	if err != nil {
		return nil, err
	}

	// post-convert
	return convert.ListUserViewsResp(resp), nil
}

func (s *Storage) SetTargetUserStatus(ctx context.Context,
	id models.UserID, status models.UserStatus,
) error {
	return doVoid(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) error {
		return q.SetTargetUserStatus(ctx,
			queries.SetTargetUserStatusParams{
				UserTargetStatus: int16(status),
				UserID:           int64(id),
			})
	})
}

func (s *Storage) DeleteUser(ctx context.Context,
	id models.UserID,
) error {
	return doVoid(ctx, s, func(ctx context.Context,
		q *queries.Queries,
	) error {
		return q.DeleteUser(ctx, int64(id))
	})
}
