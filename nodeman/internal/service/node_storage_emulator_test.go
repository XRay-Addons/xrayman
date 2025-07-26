package service

import (
	"context"
	"fmt"
	"math/rand/v2"

	"go.uber.org/zap"
)

type storageUser struct {
	user           User
	actualState    UserStatus
	requiredStatus UserStatus
}

type NodeStorageEmulator struct {
	actualState   NodeStatus
	requiredState NodeStatus

	users []storageUser

	rand                     *rand.Rand
	unavailableProb          float32
	externalModificationProb float32

	unstable bool

	log *zap.Logger
}

func NewNodeStorageEmulator(nUsers int, log *zap.Logger) *NodeStorageEmulator {
	users := make([]storageUser, 0)
	for i := range nUsers {
		state := UserEnabled
		if i%2 == 0 {
			state = UserDisabled
		}
		users = append(users, storageUser{
			User{ID: UserID(i)},
			UserStatusUnknown,
			state,
		})
	}

	return &NodeStorageEmulator{
		actualState:              NodeStatusUnknown,
		requiredState:            NodeStopped,
		users:                    users,
		rand:                     rand.New(rand.NewPCG(0, 0)),
		unavailableProb:          0.25,
		externalModificationProb: 0.25,
		unstable:                 true,
		log:                      log,
	}
}

func (n *NodeStorageEmulator) GetAllUsers(ctx context.Context) ([]UserState, error) {
	if err := n.unavaliableIncident(); err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	n.applyExternalModifications()

	newSlice := make([]UserState, 0, len(n.users))
	for _, u := range n.users {
		newSlice = append(newSlice, UserState{
			u.user,
			u.requiredStatus,
		})
	}

	return newSlice, nil
}

func (n *NodeStorageEmulator) GetNodeState(ctx context.Context) (actual, required NodeStatus, err error) {
	if err := n.unavaliableIncident(); err != nil {
		return NodeStatusUnknown, NodeStatusUnknown, fmt.Errorf("get all users: %w", err)
	}
	n.applyExternalModifications()

	return n.actualState, n.requiredState, nil
}

// GetPendingUsers implements NodeStorage.
func (n *NodeStorageEmulator) GetOutOfSyncUsers(ctx context.Context) ([]UserState, error) {
	if err := n.unavaliableIncident(); err != nil {
		return nil, fmt.Errorf("get pending users: %w", err)
	}
	n.applyExternalModifications()

	pendingUsers := make([]UserState, 0, len(n.users))
	for _, u := range n.users {
		if u.actualState != u.requiredStatus {
			pendingUsers = append(pendingUsers, UserState{
				u.user,
				u.requiredStatus,
			})
		}
	}
	return pendingUsers, nil
}

type Updater struct {
	parent *NodeStorageEmulator
	cfg    *NodeConfig
	state  NodeStatus
	users  []UserState
}

func (u *Updater) Apply(ctx context.Context) error {
	return u.parent.apply(u)
}

func (u *Updater) SetConfig(cfg *NodeConfig) {
	u.cfg = cfg
}

func (u *Updater) SetStatus(state NodeStatus) {
	u.state = state
}

func (u *Updater) SetUsers(users []UserState) {
	u.users = users
}

var _ NodeUpdater = (*Updater)(nil)

func (n *NodeStorageEmulator) GetUpdater(ctx context.Context) (NodeUpdater, error) {
	if err := n.unavaliableIncident(); err != nil {
		return nil, fmt.Errorf("node get updater: %w", err)
	}
	n.applyExternalModifications()

	return &Updater{parent: n}, nil
}

func (n *NodeStorageEmulator) apply(upd *Updater) error {
	if err := n.unavaliableIncident(); err != nil {
		return fmt.Errorf("node apply update: %w", err)
	}
	if upd.state > 0 {
		n.actualState = upd.state
	}

	for _, u := range upd.users {
		for i, uu := range n.users {
			if uu.user == u.User {
				updatedU := n.users[i]
				updatedU.actualState = u.Status
				n.users[i] = updatedU
			}
		}
	}

	n.log.Sugar().Infof("node storage: status updated: %v", upd.state)

	return nil
}

func (n *NodeStorageEmulator) unavaliableIncident() error {
	if !n.unstable {
		return nil
	}
	if n.rand.Float32() > n.unavailableProb {
		return nil
	}
	n.log.Warn("node storage: incident: storage unavailable")
	return fmt.Errorf("storage unavailable")
}

func (n *NodeStorageEmulator) applyExternalModifications() {
	if !n.unstable {
		return
	}
	if n.rand.Float32() > n.unavailableProb {
		return
	}

	if n.rand.IntN(2) == 1 {
		n.requiredState = (NodeStopped + NodeRunning) - n.requiredState
		n.log.Sugar().Warn("node storage: external state switch: ", n.requiredState)
	}

	editIdx := n.rand.IntN(len(n.users))
	editState := n.users[editIdx].requiredStatus
	if n.rand.IntN(2) == 1 {
		editState = (UserDisabled + UserEnabled) - editState
		n.log.Sugar().Warnf("node storage: external state switch: user %d -> %v",
			n.users[editIdx].user.ID, editState)
		n.users[editIdx].requiredStatus = editState
	}
}

var _ NodeStorage = (*NodeStorageEmulator)(nil)
