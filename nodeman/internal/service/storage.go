package service

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/uow"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodesStorage interface {
	// add new node to storage, assign NodeID to node
	NewNode(ctx context.Context, node *models.Node) error
	// get all nodes
	ListNodes(ctx context.Context) ([]models.Node, error)
	// change node target status
	SetTargetNodeStatus(ctx context.Context, id models.NodeID, status models.NodeStatus) error
}

type UsersStorage interface {
	// add new user to storage, assign UserID to user
	NewUser(ctx context.Context, user *models.User) error
	// get user by id
	GetUser(ctx context.Context, id models.UserID) (*models.User, error)
	// get all users
	ListUsers(ctx context.Context) ([]models.User, error)
	// change user target status
	SetTargetUserStatus(ctx context.Context, id models.UserID, status models.UserStatus) error
}

type UserNodesStorage interface {
	// get nodes where user is active
	GetUserNodes(ctx context.Context, id models.UserID) ([]models.Node, error)
}

type UoWContext interface {
	NodesStorage
	UsersStorage
	UserNodesStorage
}

type UoWFn = uow.Fn[UoWContext]

type UoW = uow.UoW[UoWContext]
