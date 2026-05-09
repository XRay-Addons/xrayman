package memstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (s *Storage) FindPendingSyncs(ctx context.Context,
	id models.NodeID,
) ([]models.UserSyncStatus, error) {
	var syncStatus []models.UserSyncStatus
	for userId, user := range s.users {
		if user.TargetStatus == s.syncStatus[id][userId] {
			continue
		}
		syncStatus = append(syncStatus, models.UserSyncStatus{
			User:          user,
			CurrentStatus: s.syncStatus[id][userId],
		})
	}
	return syncStatus, nil
}

func (s *Storage) UpdateNodeUsers(ctx context.Context,
	id models.NodeID, patch []models.UserStatusPatch,
) error {
	for _, p := range patch {
		s.syncStatus[id][p.UserID] = p.Status
	}
	return nil
}

func (s *Storage) SetNodeUsers(ctx context.Context,
	id models.NodeID, patch []models.UserStatusPatch,
) error {
	for userID := range s.syncStatus[id] {
		s.syncStatus[id][userID] = models.UserStatusDisabled
	}
	for _, p := range patch {
		s.syncStatus[id][p.UserID] = p.Status
	}
	return nil
}

func (s *Storage) GetUserNodes(ctx context.Context,
	id models.UserID,
) ([]models.Node, error) {
	var nodes []models.Node

	for _, node := range s.nodes {
		userNode := node.CurrentStatus == models.NodeStatusRunning &&
			node.TargetStatus == models.NodeStatusRunning &&
			s.syncStatus[node.ID][id] == models.UserStatusEnabled

		if userNode {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}
