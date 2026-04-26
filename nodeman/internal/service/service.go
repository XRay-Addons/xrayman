package service

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/poolsyncer"
	"github.com/XRay-Addons/xrayman/nodeman/internal/subscrman"
)

type Service struct {
	storage    Storage
	poolSyncer poolsyncer.Syncer
	subscrMan  subscrman.SubscrMan
}

var _ handler.Service = (*Service)(nil)

func New(poolSyncer poolsyncer.Syncer,
	subscrMan subscrman.SubscrMan,
	storage Storage,
) (*Service, error) {

	if poolSyncer == nil {
		return nil, errdefs.NewNilArg("poolSyncer")
	}
	if subscrMan == nil {
		return nil, errdefs.NewNilArg("subscrMan")
	}
	if storage == nil {
		return nil, errdefs.NewNilArg("storage")
	}

	return &Service{
		storage:    storage,
		poolSyncer: poolSyncer,
		subscrMan:  subscrMan,
	}, nil
}

func (s *Service) NewNode(ctx context.Context, p models.NewNodeParams) (*models.NewNodeResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	var node models.Node
	node.Config.ConnectionInfo.Endpoint = p.Endpoint
	node.Config.ConnectionInfo.AccessKey = p.AccessKey

	node.CurrentStatus = models.NodeStatusStopped
	node.TargetStatus = models.NodeStatusRunning
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewNode(ctx, &node)
		return
	}); err != nil {
		return nil, err
	}

	_ = s.syncNode(ctx, node.ID)

	return &models.NewNodeResult{
		Node:       node,
	}, nil
}

func (s *Service) StartNode(ctx context.Context, p models.StartNodeParams) (*models.StartNodeResult, error) {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusRunning); err != nil {
		return nil, err
	}
	return &models.StartNodeResult{}, nil
}

func (s *Service) StopNode(ctx context.Context, p models.StopNodeParams) (*models.StopNodeResult, error) {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusStopped); err != nil {
		return nil, err
	}
	return &models.StopNodeResult{}, nil
}

func (s *Service) ListNodes(ctx context.Context, p models.ListNodeParams) (*models.ListNodeResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	var nodes []models.Node
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		nodes, err = uowctx.ListNodes(ctx)
		return
	}); err != nil {
		return nil, err
	}
	return &models.ListNodeResult{
		Nodes: nodes,
	}, nil
}

func (s *Service) NewUser(ctx context.Context, p models.NewUserParams) (*models.User, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
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

	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewUser(ctx, &user)
		return
	}); err != nil {
		return nil, err
	}

	_ = s.syncAllNodes(ctx)

	return &user, nil
}

func (s *Service) GetUser(ctx context.Context, p models.GetUserParams) (*models.User, bool, error) {
	if s == nil {
		return nil, false, errdefs.NewNilCall()
	}

	// find user with given id
	var user *models.User
	var exists bool
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		user, exists, err = uowctx.GetUser(ctx, p.ID)
		return
	}); err != nil {
		return nil, false, err
	}

	// check user name
	if !exists || user.Profile.Name != p.Name {
		return nil, false, nil
	}

	return user, true, nil
}

func (s *Service) ListUsers(ctx context.Context, p models.ListUserParams) (*models.ListUsersResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	var users []models.User
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		users, err = uowctx.ListUsers(ctx)
		return
	}); err != nil {
		return nil, err
	}
	return &models.ListUsersResult{
		Users: users,
	}, nil
}

func (s *Service) EnableUser(ctx context.Context, p models.EnableUserParams) (*models.EnableUserResult, error) {
	if err := s.setUserStatus(ctx, p.ID, models.UserStatusEnabled); err != nil {
		return nil, err
	}
	return &models.EnableUserResult{}, nil
}

func (s *Service) DisableUser(ctx context.Context, p models.DisableUserParams) (*models.DisableUserResult, error) {
	if err := s.setUserStatus(ctx, p.ID, models.UserStatusDisabled); err != nil {
		return nil, err
	}
	return &models.DisableUserResult{}, nil
}

func (s *Service) GetUserSub(ctx context.Context, p models.UserSubParams) (
	*models.UserSubResult, bool, error,
) {
	if s == nil {
		return nil, false, errdefs.NewNilCall()
	}

	// validate user
	user, exists, err := s.findUser(ctx, p)
	if err != nil || !exists {
		return nil, exists, err
	}

	// get user sub
	userSub, err := s.subscrMan.GetUserSub(ctx, *user)
	if err != nil {
		return nil, false, err
	}

	return userSub, true, nil
}

func (s *Service) setNodeStatus(ctx context.Context, id models.NodeID, status models.NodeStatus) error {
	if s == nil {
		return errdefs.NewNilCall()
	}
	// set target node state to storage
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetTargetNodeStatus(ctx, id, status)
		return
	}); err != nil {
		return err
	}

	_ = s.syncNode(ctx, id)
	return nil
}

func (s *Service) syncNode(ctx context.Context, id models.NodeID) error {
	syncResults, err := s.poolSyncer.SyncPoolState(ctx)
	if err != nil {
		return err
	}
	if err = syncResults.GetNodeErr(id); err != nil {
		return err
	}
	return nil
}

func (s *Service) setUserStatus(ctx context.Context, id models.UserID, status models.UserStatus) error {
	if s == nil {
		return errdefs.NewNilCall()
	}
	// set target user state to storage
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetTargetUserStatus(ctx, id, status)
		return
	}); err != nil {
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
	return errors.Join(errs...)
}

func (s *Service) findUser(ctx context.Context, p models.UserSubParams) (*models.User, bool, error) {
	// find user with given id
	var user *models.User
	var exists bool
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		user, exists, err = uowctx.GetUser(ctx, p.ID)
		return
	}); err != nil {
		return nil, false, err
	}

	// check user name
	if !exists || user.Profile.Name != p.Name {
		return nil, false, nil
	}

	return user, true, nil
}
