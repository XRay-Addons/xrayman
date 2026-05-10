package memstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) NewNode(ctx context.Context, node *models.Node) error {
	node.ID = models.NodeID(len(s.nodes))
	s.nodes = append(s.nodes, *node)
	nodeUsers := make([]models.UserStatus, len(s.users))
	for i := range nodeUsers {
		nodeUsers[i] = models.UserStatusDisabled
	}
	s.syncStatus = append(s.syncStatus, nodeUsers)
	return nil
}

func (s *Storage) SetTargetNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	s.nodes[id].TargetStatus = status
	return nil
}

func (s *Storage) ListNodes(ctx context.Context) (
	[]models.Node, error,
) {
	var nodes []models.Node
	nodes = append(nodes, s.nodes...)
	return nodes, nil
}

func (s *Storage) GetNode(ctx context.Context, id models.NodeID) (
	*models.Node, bool, error,
) {
	return &s.nodes[id], true, nil
}

func (s *Storage) SetClientConfig(ctx context.Context,
	id models.NodeID, cfg models.ClientConfigTemplate,
) error {
	s.nodes[id].Config.ClientConfigTemplate = cfg
	return nil
}

func (s *Storage) SetCurrentNodeStatus(ctx context.Context,
	id models.NodeID, status models.NodeStatus,
) error {
	s.nodes[id].CurrentStatus = status
	return nil
}

func (s *Storage) DeleteNode(ctx context.Context,
	id models.NodeID,
) error {
	s.nodes[id].CurrentStatus = 0
	s.nodes[id].TargetStatus = 0
	return nil
}
