package node

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/service/models"
)

type UserState struct {
	User   models.User
	Status models.UserStatus
}

type OutOfSyncUser struct {
	User     models.User
	Required models.UserStatus
	Actual   models.UserStatus
}

type UserStatusUpdate struct {
	ID     models.UserID
	Actual models.UserStatus
}

type StorageUpdate struct {
	ActualStatus models.UserStatus // !default is invalid value
	Properties   models.NodeProperties
	UserUpdates  []UserStatusUpdate
}

type Sy struct {
	Required models.NodeStatus
	Actual   models.NodeStatus
}

// storage Unit of work for transactional writing
type StorageWriteUOW interface {
	SetActualStatus(s models.NodeStatus)
	SetNodeProperties(p models.NodeProperties)
	SetActualUserStates(u []UserStatusUpdate)
	Do(ctx context.Context) error
}

// storage interface for node
type Storage interface {
	// get node status
	GetNodeState(ctx context.Context) (required, actual models.NodeStatus, err error)
	// get all users related to node
	ListUsers(ctx context.Context) ([]UserState, error)
	// get users with actual state != required state
	ListOutOfSyncUsers(ctx context.Context) ([]OutOfSyncUser, error)
	// tx apply actual state update
	GetWriteUOW() StorageWriteUOW
}

// node api interface for node
type API interface {
	Start(ctx context.Context, users []models.User) (*models.NodeProperties, error)
	Stop(ctx context.Context) error
	GetStatus(ctx context.Context) (models.NodeStatus, error)
	EditUsers(ctx context.Context, users []UserState) error
}
