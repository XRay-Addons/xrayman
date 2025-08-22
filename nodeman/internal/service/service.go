package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"text/template"

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
		return nil, errdefs.NewNilArg("syncman")
	}
	if uow == nil {
		return nil, errdefs.NewNilArg("uow")
	}
	return &Service{
		syncman: syncman,
		uow:     uow,
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
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewNode(ctx, &node)
		return
	}); err != nil {
		return nil, err
	}

	_ = s.syncNode(ctx, node.ID)

	return &models.NewNodeResult{
		ID:       node.ID,
		Endpoint: node.Config.ConnectionInfo.Endpoint,
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
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		nodes, err = uowctx.ListNodes(ctx)
		return
	}); err != nil {
		return nil, err
	}
	return &models.ListNodeResult{
		Nodes: nodes,
	}, nil
}

func (s *Service) NewUser(ctx context.Context, p models.NewUserParams) (*models.NewUserResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	vlessUUID, err := generateVlessUUID()
	if err != nil {
		return nil, err
	}
	name := makeSlugName(p.VisibleName)

	var user models.User
	user.Profile.VisibleName = p.VisibleName
	user.Profile.Name = name
	user.Profile.VlessUUID = vlessUUID
	user.TargetStatus = models.UserStatusEnabled

	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.NewUser(ctx, &user)
		return
	}); err != nil {
		return nil, err
	}

	_ = s.syncAllNodes(ctx)

	return &models.NewUserResult{
		ID:          user.ID,
		VisibleName: user.Profile.VisibleName,
		UserPageURL: makeUserPageURL(user),
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

func (s *Service) ListUsers(ctx context.Context, p models.ListUserParams) (*models.ListUsersResult, error) {
	if s == nil {
		return nil, errdefs.NewNilCall()
	}
	var users []models.User
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		users, err = uowctx.ListUsers(ctx)
		return
	}); err != nil {
		return nil, err
	}
	return &models.ListUsersResult{
		Users: users,
	}, nil
}

func (s *Service) GetUserSub(ctx context.Context, p models.GetUserSubParams) (
	_ *models.GetUserSubResult, exists bool, _ error,
) {
	if s == nil {
		return nil, false, errdefs.NewNilCall()
	}

	// validate user
	var user *models.User
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		user, err = uowctx.GetUser(ctx, p.ID)
		return
	}); err != nil {
		return nil, false, err
	}
	if p.Name != user.Profile.Name {
		return nil, false, nil
	}

	// get active nodes for user
	var userNodes []models.Node
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		userNodes, err = uowctx.GetUserNodes(ctx, user.ID)
		return
	}); err != nil {
		return nil, false, err
	}

	// create subscription content
	subscriptions := make([]json.RawMessage, 0, len(userNodes))
	for _, node := range userNodes {
		tmpl, err := template.New("client").Parse(node.Config.ClientConfig.Template)
		if err != nil {
			return nil, false, errdefs.WrapWithStack(err)
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{
			node.Config.ClientConfig.UserNameField:  user.Profile.Name,
			node.Config.ClientConfig.VlessUUIDField: user.Profile.VlessUUID,
		})
		if err != nil {
			return nil, false, errdefs.WrapWithStack(err)
		}
		// parse as array or as single subscription
		var arr []json.RawMessage
		if err := json.Unmarshal(buf.Bytes(), &arr); err == nil {
			subscriptions = append(subscriptions, arr...)
			continue
		}
		subscriptions = append(subscriptions, json.RawMessage(buf.Bytes()))
	}

	return &subscriptions, true, nil
}

func (s *Service) setNodeStatus(ctx context.Context, id models.NodeID, status models.NodeStatus) error {
	if s == nil {
		return errdefs.NewNilCall()
	}
	// set target node state to storage
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetTargetNodeStatus(ctx, id, status)
		return
	}); err != nil {
		return err
	}

	_ = s.syncNode(ctx, id)
	return nil
}

func (s *Service) syncNode(ctx context.Context, id models.NodeID) error {
	syncResults, err := s.syncman.SyncNodesPool(ctx)
	if err != nil {
		return err
	}
	for _, syncRes := range syncResults {
		if syncRes.ID != id {
			continue
		}
		if syncRes.Err == nil {
			return nil
		}
		return syncRes.Err
	}
	return errdefs.New("sync node not found", errdefs.Withf("node id: %v", id))
}

func (s *Service) setUserStatus(ctx context.Context, id models.UserID, status models.UserStatus) error {
	if s == nil {
		return errdefs.NewNilCall()
	}
	// set target user state to storage
	if err := s.uow.Do(ctx, func(uowctx UoWContext) (err error) {
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
	syncResults, err := s.syncman.SyncNodesPool(ctx)
	if err != nil {
		return err
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
	return errors.Join(errs...)
}
