package users

import (
	"context"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/supervisor"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	storage    Storage
	poolSyncer Syncer

	syncTimeout time.Duration
	sv          *supervisor.Supervisor

	logger *zap.Logger
}

var _ handler.UsersService = (*Service)(nil)

func New(poolSyncer Syncer,
	storage Storage,
	syncTimeout time.Duration,
	logger *zap.Logger,
) (*Service, error) {
	if poolSyncer == nil {
		return nil, errdefs.NilArg("poolSyncer")
	}
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}
	if logger == nil {
		return nil, errdefs.NilArg("logger")
	}

	return &Service{
		storage:     storage,
		poolSyncer:  poolSyncer,
		syncTimeout: syncTimeout,
		sv:          supervisor.New(),
		logger:      logger,
	}, nil
}

func (s *Service) Close() {
	if s == nil || s.sv == nil {
		return
	}
	s.sv.Close()
}

func (s *Service) requestNodesSync() {
	s.sv.Go(func(ctx context.Context) {
		if err := s.syncAllNodes(ctx); err != nil {
			s.logger.Warn("nodes sync request", zap.Error(err))
		}
	}, s.syncTimeout)
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

	// sync all nodes on user add to return valid configuration
	// to new user and avoid situation when nodes are temporary
	// unavailable, user successfully created and get empty
	// subscription config just after that.
	if err := s.storage.DoTx(ctx, func(context.Context) error {
		if err := s.storage.NewUser(ctx, &user); err != nil {
			return err
		}
		if err := s.syncAllNodes(ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

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

	if err := s.storage.DoTx(ctx, func(ctx context.Context) error {
		if err := s.storage.SetTargetUserStatus(ctx,
			p.ID, models.UserStatusDisabled,
		); err != nil {
			return err
		}
		if err := s.storage.DeleteUser(ctx, p.ID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	s.requestNodesSync()

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
	s.requestNodesSync()

	return nil
}

// sync all nodes, return nil if at least one node synced ok
func (s *Service) syncAllNodes(ctx context.Context) error {
	syncResults, err := s.poolSyncer.SyncPoolState(ctx)
	if err != nil {
		return err
	}
	return syncResults.GetEntireErr()
}
