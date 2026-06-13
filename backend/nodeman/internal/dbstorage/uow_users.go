package dbstorage

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	queries "github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage/sqlc/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

const month = 30 * 24 * time.Hour

func (uow *uowctx) NewUser(ctx context.Context, user *models.User) error {
	// pre-convert
	arg, err := Convert[models.User, queries.NewUserParams](user,
		With(func(from *models.User, to *queries.NewUserParams) {
			to.DisplayName = from.Profile.DisplayName
			to.UserName = from.Profile.Name
			to.UserTargetStatus = int16(from.TargetStatus)
			to.VlessUuid = from.Profile.VlessUUID
		}),
	)
	if err != nil {
		return err
	}

	// request
	userID, err := uow.q.NewUser(ctx, *arg)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	// post-convert
	user.Profile.ID = models.UserID(userID)

	return nil
}

func (uow *uowctx) GetUserView(ctx context.Context,
	id models.UserID, name string,
) (*models.UserView, error) {

	resp, err := uow.q.GetUserView(ctx, queries.GetUserViewParams{
		FromDay:  time.Now().Add(-month),
		UserID:   int64(id),
		UserName: name,
	})
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	user, err := Convert[queries.GetUserViewRow, models.UserView](&resp,
		With(func(from *queries.GetUserViewRow, to *models.UserView) {
			to.User.Profile.ID = models.UserID(from.UserID)
			to.User.Profile.Name = from.UserName
			to.User.Profile.DisplayName = from.DisplayName
			to.User.Profile.VlessUUID = from.VlessUuid
			to.User.TargetStatus = models.UserStatus(from.UserTargetStatus)
			to.Traffic.Total.Download = from.DownloadTotal
			to.Traffic.Total.Upload = from.UploadTotal
			to.Traffic.LastMonth.Download = from.DownloadLastDays
			to.Traffic.LastMonth.Upload = from.UploadLastDays
		}),
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uow *uowctx) ListUsers(ctx context.Context) (
	[]models.User, error,
) {
	resp, err := uow.q.ListUsers(ctx)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	users, err := ConvertArray[queries.ListUsersRow, models.User](resp,
		With(func(from *queries.ListUsersRow, to *models.User) {
			to.Profile.ID = models.UserID(from.UserID)
			to.Profile.Name = from.UserName
			to.Profile.DisplayName = from.DisplayName
			to.Profile.VlessUUID = from.VlessUuid
			to.TargetStatus = models.UserStatus(from.UserTargetStatus)
		}),
	)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uow *uowctx) ListUserViews(ctx context.Context) ([]models.UserView, error) {
	resp, err := uow.q.ListUserViews(ctx, time.Now().Add(-30*24*time.Hour))
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	users, err := ConvertArray[queries.ListUserViewsRow, models.UserView](resp,
		With(func(from *queries.ListUserViewsRow, to *models.UserView) {
			to.User.Profile.ID = models.UserID(from.UserID)
			to.User.Profile.Name = from.UserName
			to.User.Profile.DisplayName = from.DisplayName
			to.User.Profile.VlessUUID = from.VlessUuid
			to.User.TargetStatus = models.UserStatus(from.UserTargetStatus)
			to.Traffic.Total.Upload = from.UploadTotal
			to.Traffic.Total.Download = from.DownloadTotal
			to.Traffic.LastMonth.Download = from.DownloadLastDays
			to.Traffic.LastMonth.Upload = from.UploadLastDays
		}),
	)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uow *uowctx) SetTargetUserStatus(ctx context.Context,
	id models.UserID, status models.UserStatus,
) error {
	if err := uow.q.SetTargetUserStatus(ctx, queries.SetTargetUserStatusParams{
		UserTargetStatus: int16(status),
		UserID:           int64(id),
	}); err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) DeleteUser(ctx context.Context,
	id models.UserID,
) error {
	if err := uow.q.DeleteUser(ctx, int64(id)); err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}
