package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"go.uber.org/zap"
)

// user description
type UserID int

type User struct {
	ID        UserID
	Name      string
	VlessUUID string
}

type UserStatus int

const (
	UserDisabled UserStatus = iota + 1
	UserEnabled
)

func (s UserStatus) String() string {
	switch s {
	case UserDisabled:
		return "Off"
	case UserEnabled:
		return "On"
	default:
		return "Unknown"
	}
}

// node description

//go:generate stringer -type=NodeState
type NodeState int

const (
	NodeUnavailable NodeState = iota + 1
	NodeStopped
	NodeRunning
)

func (s NodeState) String() string {
	switch s {
	case NodeUnavailable:
		return "Unavailable"
	case NodeStopped:
		return "Stopped"
	case NodeRunning:
		return "Running"
	default:
		return "Unknown"
	}
}

type UserState struct {
	User   User
	Status UserStatus
}

type NodeConfig struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mocks/node_api.go -package=mocks
type NodeAPI interface {
	Start(ctx context.Context, users []User) (*NodeConfig, error)
	Stop(ctx context.Context) error
	Status(ctx context.Context) (NodeState, error)
	EditUsers(ctx context.Context, users []UserState) error
	Close(ctx context.Context) error
}

type NodeUpdater interface {
	SetConfig(cfg *NodeConfig)
	SetState(state NodeState)
	SetUsers(users []UserState)
	Apply(ctx context.Context) error
}

type NodeStorage interface {
	GetNodeState(ctx context.Context) (actual, required NodeState, err error)
	GetPendingUsers(ctx context.Context) ([]UserState, error)
	GetAllUsers(ctx context.Context) ([]UserState, error)

	GetUpdater(ctx context.Context) (NodeUpdater, error)
}

type NodeController struct {
	nodeAPI      NodeAPI
	stateStorage NodeStorage
	log          *zap.Logger
}

func New(api NodeAPI, storage NodeStorage, log *zap.Logger) (*NodeController, error) {
	if api == nil {
		return nil, fmt.Errorf("node init: api: %w", errdefs.ErrNilArgPassed)
	}
	if storage == nil {
		return nil, fmt.Errorf("node init: storage: %w", errdefs.ErrNilArgPassed)
	}
	if log == nil {
		return nil, fmt.Errorf("node init: logger: %w", errdefs.ErrNilArgPassed)
	}

	return &NodeController{
		nodeAPI:      api,
		stateStorage: storage,
		log:          log,
	}, nil

}

func (c *NodeController) Close(ctx context.Context) error {
	if c == nil || c.nodeAPI == nil {
		return nil
	}
	if err := c.nodeAPI.Close(ctx); err != nil {
		return fmt.Errorf("node: close: %w", err)
	}
	return nil
}

func (c *NodeController) SyncNodeStatus(ctx context.Context) (err error) {
	if c == nil || c.nodeAPI == nil || c.stateStorage == nil {
		return fmt.Errorf("node: sync: %w", errdefs.ErrNilObjectCall)
	}
	prevState, requiredState, err := c.getStoredState(ctx)
	if err != nil {
		return fmt.Errorf("node: sync: %w", err)
	}
	c.log.Sugar().Infof("sync node: prev: %v, required: %v", prevState, requiredState)

	// check status if required
	currentState := prevState
	if c.statusCheckRequired(currentState, requiredState) {
		if currentState, err = c.getNodeState(ctx); err != nil {
			return fmt.Errorf("node: sync: check status: %w", err)
		}
		c.log.Sugar().Infof("sync node: curr: %v, required: %v", currentState, requiredState)
	}

	//if c.statusChangingRequired(currentState, requiredState) {

	//}

	// update stored node state on exit
	stateUpdater, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		return fmt.Errorf("node: sync: %w", err)
	}
	defer func() {
		if updateErr := stateUpdater.Apply(ctx); updateErr != nil {
			err = errors.Join(err, fmt.Errorf("node: sync: update: %w", updateErr))
			return
		}
		c.log.Info("sync node: update state OK")
	}()
	stateUpdater.SetState(currentState)

	c.log.Sugar().Infof("sync node: change state %v -> %v", currentState, requiredState)

	// start, stop or edit node users
	switch {
	case c.nodeStartRequired(currentState, requiredState):
		if err = c.prepareToChangeState(ctx); err != nil {
			return fmt.Errorf("node: sync: %w", err)
		}
		c.log.Info("sync node: start")
		err = c.startNode(ctx, stateUpdater)
	case c.nodeStopRequired(currentState, requiredState):
		if err = c.prepareToChangeState(ctx); err != nil {
			return fmt.Errorf("node: sync: %w", err)
		}
		c.log.Info("sync node: stop")
		err = c.stopNode(ctx, stateUpdater)
	case c.editUsersRequired(currentState, requiredState):
		c.log.Info("sync node: edit users")
		err = c.editNodeUsers(ctx, stateUpdater)
	}

	if err != nil {
		return fmt.Errorf("node: sync: %w", err)
	}
	c.log.Info("sync node: state changed OK")

	return nil
}

func (c *NodeController) getStoredState(ctx context.Context) (actual, required NodeState, err error) {
	actual, required, err = c.stateStorage.GetNodeState(ctx)
	if err != nil {
		err = fmt.Errorf("get stored state: %w", err)
	}
	return
}

func (c *NodeController) statusCheckRequired(actual, required NodeState) bool {
	switch {
	case required == NodeStopped:
		// no check needed before stopping node
		return false
	case actual == NodeStopped:
		// no check required before starting node
		return false
	case actual == NodeUnavailable:
		// check required to ensure unavailability is temporary or not
		return false
	case actual == NodeRunning:
		// check required to ensure node is still running
		return true
	default:
		// no way we are here
		return false
	}
}

func (c *NodeController) statusChangingRequired(actual, required NodeState) bool {
	switch {
	case required == NodeRunning && actual != NodeRunning:
		return true
	case required == NodeStopped && actual == NodeRunning:
		return true
	default:
		return false
	}
}

func (c *NodeController) prepareToChangeState(ctx context.Context) error {
	updater, err := c.stateStorage.GetUpdater(ctx)
	if err != nil {
		return fmt.Errorf("prepare to change state: %w", err)
	}
	updater.SetState(NodeUnavailable)

	if err := updater.Apply(ctx); err != nil {
		return fmt.Errorf("prepare to change state: %w", err)
	}
	return nil
}

func (c *NodeController) getNodeState(ctx context.Context) (NodeState, error) {
	state, err := c.nodeAPI.Status(ctx)
	if err != nil {
		return state, fmt.Errorf("get node status: %w", err)
	}
	return state, nil
}

func (c *NodeController) nodeStartRequired(actual, required NodeState) bool {
	// start required if node is stopped or unavailable
	return required == NodeRunning && actual != NodeRunning
}

func (c *NodeController) nodeStopRequired(actual, required NodeState) bool {
	// stop only available nodes, unavailable skip,
	// use it only when node starting required
	return required == NodeStopped && actual == NodeRunning
}

func (c *NodeController) editUsersRequired(actual, required NodeState) bool {
	// edit users only on running nodes
	return required == NodeRunning && actual == NodeRunning
}

func (c *NodeController) startNode(ctx context.Context, updater NodeUpdater) error {
	// get list of all users
	allUsers, err := c.stateStorage.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("start node: get all users: %w", err)
	}

	// select only enabled users,
	enabledUsers := make([]User, 0, len(allUsers))
	for _, u := range allUsers {
		if u.Status != UserEnabled {
			continue
		}
		enabledUsers = append(enabledUsers, u.User)
	}

	// start node
	updater.SetState(NodeUnavailable)
	cfg, err := c.nodeAPI.Start(ctx, enabledUsers)
	if err != nil {
		return fmt.Errorf("start node: %w", err)
	}

	// update node state in storage
	updater.SetState(NodeRunning)
	updater.SetConfig(cfg)
	updater.SetUsers(allUsers)

	return nil
}

func (c *NodeController) stopNode(ctx context.Context, updater NodeUpdater) error {
	if err := c.nodeAPI.Stop(ctx); err != nil {
		return fmt.Errorf("stop node: %w", err)
	}
	c.log.Info("sync node: api stopped")
	updater.SetState(NodeStopped)

	// all requests to make something from this node is not pending anymore
	pendingUsers, err := c.stateStorage.GetPendingUsers(ctx)
	if err != nil {
		return fmt.Errorf("stop node: %w", err)
	}
	updater.SetUsers(pendingUsers)

	return nil
}

func (c *NodeController) editNodeUsers(ctx context.Context, updater NodeUpdater) error {
	pendingUsers, err := c.stateStorage.GetPendingUsers(ctx)
	if err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}
	if err := c.nodeAPI.EditUsers(ctx, pendingUsers); err != nil {
		return fmt.Errorf("edit node users: %w", err)
	}
	updater.SetUsers(pendingUsers)
	return nil
}
