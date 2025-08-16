package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Service struct {
	syncman SyncMan
	uow     UoW
}

var _ handler.Service = (*Service)(nil)

func New(syncman SyncMan, uow UoW) (*Service, error) {
	if syncman == nil {
		return nil, fmt.Errorf("service init: syncman: %w", errdefs.ErrNilArgPassed)
	}
	if uow == nil {
		return nil, fmt.Errorf("service init: uow: %w", errdefs.ErrNilArgPassed)
	}
	return &Service{
		syncman: syncman,
		uow:     uow,
	}, nil
}

func (s *Service) NewNode(ctx context.Context, p models.NewNodeParams) (*models.NewNodeResult, error) {
	if s == nil {
		return nil, fmt.Errorf("service: start: %w", errdefs.ErrNilObjectCall)
	}
	var node models.Node
	node.Config.ConnectionInfo.Endpoint = p.Endpoint
	node.Config.ConnectionInfo.AccessKey = p.AccessKey

	node.CurrentStatus = models.NodeStatusStopped
	node.TargetStatus = models.NodeStatusRunning
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewNode(ctx, &node)
		return
	}); err != nil {
		return nil, fmt.Errorf("service: new node: %w", err)
	}

	_ = s.syncNode(ctx, node.ID)

	return &models.NewNodeResult{
		ID:       node.ID,
		Endpoint: node.Config.ConnectionInfo.Endpoint,
	}, nil
}

func (s *Service) StartNode(ctx context.Context, p models.StartNodeParams) (*models.StartNodeResult, error) {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusRunning); err != nil {
		return nil, fmt.Errorf("service: start node: %w", err)
	}
	return &models.StartNodeResult{}, nil
}

func (s *Service) StopNode(ctx context.Context, p models.StopNodeParams) (*models.StopNodeResult, error) {
	if err := s.setNodeStatus(ctx, p.ID, models.NodeStatusStopped); err != nil {
		return nil, fmt.Errorf("service: stop node: %w", err)
	}
	return &models.StopNodeResult{}, nil
}

func (s *Service) ListNodes(ctx context.Context, p models.ListNodeParams) (*models.ListNodeResult, error) {
	if s == nil {
		return nil, fmt.Errorf("service: list nodes: %w", errdefs.ErrNilObjectCall)
	}
	var nodes []models.Node
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		nodes, err = uowctx.ListNodes(ctx)
		return
	}); err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	return &models.ListNodeResult{
		Nodes: nodes,
	}, nil
}

func (s *Service) NewUser(ctx context.Context, p models.NewUserParams) (*models.NewUserResult, error) {
	if s == nil {
		return nil, fmt.Errorf("service: new user: %w", errdefs.ErrNilObjectCall)
	}
	vlessUUID, err := generateVlessUUID()
	if err != nil {
		return nil, fmt.Errorf("service: new user: %w", err)
	}
	slugName := makeSlugName(p.Name)

	var user models.User
	user.Profile.Name = p.Name
	user.Profile.SlugName = slugName
	user.Profile.VlessUUID = vlessUUID
	user.TargetStatus = models.UserStatusEnabled

	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewUser(ctx, &user)
		return
	}); err != nil {
		return nil, fmt.Errorf("service: new node: %w", err)
	}

	_ = s.syncAllNodes(ctx)

	return &models.NewUserResult{
		ID:          user.ID,
		Name:        user.Profile.Name,
		UserPageURL: makeUserPageURL(user),
	}, nil
}

func (s *Service) EnableUser(ctx context.Context, p models.EnableUserParams) (*models.EnableUserResult, error) {
	if err := s.setUserStatus(ctx, p.ID, models.UserStatusEnabled); err != nil {
		return nil, fmt.Errorf("service: enable user: %w", err)
	}
	return &models.EnableUserResult{}, nil
}

func (s *Service) DisableUser(ctx context.Context, p models.DisableUserParams) (*models.DisableUserResult, error) {
	if err := s.setUserStatus(ctx, p.ID, models.UserStatusDisabled); err != nil {
		return nil, fmt.Errorf("service: disable user: %w", err)
	}
	return &models.DisableUserResult{}, nil
}

func (s *Service) ListUsers(ctx context.Context, p models.ListUserParams) (*models.ListUsersResult, error) {
	if s == nil {
		return nil, fmt.Errorf("service: list users: %w", errdefs.ErrNilObjectCall)
	}
	var users []models.User
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		users, err = uowctx.ListUsers(ctx)
		return
	}); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return &models.ListUsersResult{
		Users: users,
	}, nil
}

func (s *Service) setNodeStatus(ctx context.Context, id models.NodeID, status models.NodeStatus) error {
	if s == nil {
		return fmt.Errorf("set node status: %w", errdefs.ErrNilObjectCall)
	}
	// set target node state to storage
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetTargetNodeStatus(ctx, id, status)
		return
	}); err != nil {
		return fmt.Errorf("set node status: %w", err)
	}

	_ = s.syncNode(ctx, id)
	return nil
}

func (s *Service) syncNode(ctx context.Context, id models.NodeID) error {
	syncResults, err := s.syncman.SyncNodesPool(ctx)
	if err != nil {
		return fmt.Errorf("service: sync node: %w", err)
	}
	for _, syncRes := range syncResults {
		if syncRes.ID != id {
			continue
		}
		if syncRes.Err == nil {
			return nil
		}
		return fmt.Errorf("service: sync node: %w", syncRes.Err)
	}
	return fmt.Errorf("service: sync node: node not found: %w", errdefs.ErrIPE)
}

func (s *Service) setUserStatus(ctx context.Context, id models.UserID, status models.UserStatus) error {
	if s == nil {
		return fmt.Errorf("set user status: %w", errdefs.ErrNilObjectCall)
	}
	// set target user state to storage
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetTargetUserStatus(ctx, id, status)
		return
	}); err != nil {
		return fmt.Errorf("set user status: %w", err)
	}

	// sync nodes. errors is not a problem, it will updates in background
	_ = s.syncAllNodes(ctx)
	return nil
}

// sync all nodes, return nil if at least one node synced ok
func (s *Service) syncAllNodes(ctx context.Context) error {
	syncResults, err := s.syncman.SyncNodesPool(ctx)
	if err != nil {
		return fmt.Errorf("service: sync all nodes: %w", err)
	}
	if len(syncResults) == 0 {
		return nil
	}
	var errs []error
	for _, syncRes := range syncResults {
		if syncRes.Err == nil {
			return nil
		}
		errs = append(errs, syncRes.Err)
	}
	return fmt.Errorf("servuce: sync all nodes: %w", errors.Join(errs...))
}
