package service

import (
	"context"
	"fmt"
	"math/rand/v2"

	"go.uber.org/zap"
)

// implement node emulator for tests
type NodeAPIEmulator struct {
	status NodeState
	users  map[UserID]struct{}

	rand            *rand.Rand
	failProb        float32
	turnOffProb     float32
	unavailableProb float32

	unstable bool

	log *zap.Logger
}

func NewNodeAPIEmulator(log *zap.Logger) *NodeAPIEmulator {
	return &NodeAPIEmulator{
		status:          NodeStopped,
		users:           make(map[UserID]struct{}, 0),
		rand:            rand.New(rand.NewPCG(0, 0)),
		failProb:        0.25,
		turnOffProb:     0.25,
		unavailableProb: 0.25,
		unstable:        true,
		log:             log,
	}
}

func (n *NodeAPIEmulator) Close(ctx context.Context) error {
	return nil
}

func (n *NodeAPIEmulator) EditUsers(ctx context.Context, users []UserState) error {
	if err := n.unavaliableIncident(); err != nil {
		return fmt.Errorf("edit users: %w", err)
	}
	n.turnOffIncident()

	if n.status != NodeRunning {
		return fmt.Errorf("node not running")
	}

	for _, u := range users {
		if u.Status == UserDisabled {
			delete(n.users, u.User.ID)
		} else {
			n.users[u.User.ID] = struct{}{}
		}
	}
	n.log.Info("node: users updated")

	return nil
}

func (n *NodeAPIEmulator) Start(ctx context.Context, users []User) (*NodeConfig, error) {
	if err := n.unavaliableIncident(); err != nil {
		return nil, fmt.Errorf("edit users: %w", err)
	}
	n.turnOffIncident()

	for ui := range n.users {
		delete(n.users, ui)
	}

	for _, u := range users {
		n.users[u.ID] = struct{}{}
	}
	n.status = NodeRunning
	n.log.Info("node: started")

	return &NodeConfig{}, nil
}

func (n *NodeAPIEmulator) Status(ctx context.Context) (NodeState, error) {
	if err := n.unavaliableIncident(); err != nil {
		return NodeUnavailable, nil
	}
	n.turnOffIncident()
	return n.status, nil
}

func (n *NodeAPIEmulator) Stop(ctx context.Context) error {
	if err := n.unavaliableIncident(); err != nil {
		return fmt.Errorf("node api: %w", err)
	}
	n.turnOffIncident()

	n.status = NodeStopped
	for ui := range n.users {
		delete(n.users, ui)
	}
	n.log.Info("node: stopped")
	return nil
}

func (n *NodeAPIEmulator) unavaliableIncident() error {
	if !n.unstable {
		return nil
	}
	if n.rand.Float32() > n.unavailableProb {
		return nil
	}
	n.log.Info("node: temporary unavailable")
	return fmt.Errorf("node unavailable")
}

func (n *NodeAPIEmulator) turnOffIncident() {
	if !n.unstable {
		return
	}
	if n.rand.Float32() > n.turnOffProb {
		return
	}
	n.log.Info("node: occasionally turned off")
	n.status = NodeStopped
	n.users = make(map[UserID]struct{}, 0)
}

var _ NodeAPI = (*NodeAPIEmulator)(nil)
