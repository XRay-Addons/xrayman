package users

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Service struct {
	storage    Storage
	poolSyncer Syncer
}

var _ handler.UsersService = (*Service)(nil)

func New(poolSyncer Syncer,
	storage Storage,
) (*Service, error) {
	if poolSyncer == nil {
		return nil, errdefs.NilArg("poolSyncer")
	}
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}

	return &Service{
		storage:    storage,
		poolSyncer: poolSyncer,
	}, nil
}

func (s *Service) NewUser(ctx context.Context, p models.NewUserParams) (
	*models.User, error,
) {
	if s == nil {
		return nil, errdefs.NilCall()
	}
	vlessUUID, err := generateVlessUUID()
	if err != nil {
		return nil, err
	}
	name := makeSlugName(p.DisplayName)

	var user models.User
	user.Profile.DisplayName = p.DisplayName
	user.Profile.Name = name
	user.Profile.VlessUUID = vlessUUID
	user.TargetStatus = models.UserStatusEnabled

	if err := s.storage.NewUser(ctx, &user); err != nil {
		return nil, err
	}

	_ = s.syncAllNodes(ctx)

	return &user, nil
}

func (s *Service) GetUserView(ctx context.Context, p models.GetUserParams) (
	*models.UserView, error,
) {
	if s == nil {
		return nil, errdefs.NilCall()
	}

	// find user with given id
	userView, err := s.storage.GetUserView(ctx, p.ID, p.Name)
	if err != nil {
		return nil, err
	}

	return userView, nil
}

func (s *Service) ListUsers(ctx context.Context, p models.ListUserParams) (
	*models.ListUsersResult, error,
) {
	if s == nil {
		return nil, errdefs.NilCall()
	}
	users, err := s.storage.ListUserViews(ctx)
	if err != nil {
		return nil, err
	}
	return &models.ListUsersResult{
		Users: users,
	}, nil
}

func (s *Service) EnableUser(ctx context.Context, p models.EnableUserParams) (
	*models.EnableUserResult, error,
) {
	if err := s.setUserStatus(ctx, p.ID, models.UserStatusEnabled); err != nil {
		return nil, err
	}
	return &models.EnableUserResult{}, nil
}

func (s *Service) DisableUser(ctx context.Context, p models.DisableUserParams) (
	*models.DisableUserResult, error,
) {
	if err := s.setUserStatus(ctx, p.ID, models.UserStatusDisabled); err != nil {
		return nil, err
	}
	return &models.DisableUserResult{}, nil
}

func (s *Service) DeleteUser(ctx context.Context, p models.DeleteUserParams) (
	*models.DeleteUserResult, error,
) {
	if s == nil {
		return nil, errdefs.NilCall()
	}
	// disable user before deleting
	if err := s.setUserStatus(ctx, p.ID, models.UserStatusDisabled); err != nil {
		return nil, err
	}

	if err := s.storage.DeleteUser(ctx, p.ID); err != nil {
		return nil, err
	}

	_ = s.syncAllNodes(ctx)

	return &models.DeleteUserResult{}, nil
}

func (s *Service) setUserStatus(ctx context.Context,
	id models.UserID, status models.UserStatus,
) error {
	if s == nil {
		return errdefs.NilCall()
	}
	// set target user state to storage
	if err := s.storage.SetTargetUserStatus(ctx, id, status); err != nil {
		return err
	}

	// sync nodes. errors is not a problem, it will updates in background
	_ = s.syncAllNodes(ctx)
	return nil
}

// sync all nodes, return nil if at least one node synced ok
func (s *Service) syncAllNodes(ctx context.Context) error {
	syncResults, err := s.poolSyncer.SyncPoolState(ctx)
	if err != nil {
		return err
	}
	if len(syncResults.Nodes) == 0 {
		return nil
	}
	var errs []error
	for _, syncRes := range syncResults.Nodes {
		if syncRes.Err == nil {
			return nil
		}
		errs = append(errs, syncRes.Err)
	}
	return xerr.Join(errs...)
}
