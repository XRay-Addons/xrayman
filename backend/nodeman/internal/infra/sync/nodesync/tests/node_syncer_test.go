package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/stretchr/testify/require"
)

func TestNodeSync(t *testing.T) {

	t.Run("memory stable", func(t *testing.T) {
		s := NewMemoryMockStorage()
		testNodeSync(t, s)
	})

	//t.Run("db stable", func(t *testing.T) {
	//	s, _ := setupTestDB(t, zap.NewNop())
	//	testNodeSync(t, s)
	//})

	t.Run("memory unstable storage", func(t *testing.T) {
		s := NewMemoryMockStorage()
		testNodeSync_UnstableStorage(t, s)
	})

	//t.Run("db unstable storage", func(t *testing.T) {
	//	s, _ := setupTestDB(t, zap.NewNop())
	//	testNodeSync_UnstableStorage(t, s)
	//})

	t.Run("memory unstable storage unstable node", func(t *testing.T) {
		s := NewMemoryMockStorage()
		testNodeSync_UnstableStorage_UnstableNode(t, s)
	})

	//t.Run("db unstable storage unstable node", func(t *testing.T) {
	//	s, _ := setupTestDB(t, zap.NewNop())
	//	testNodeSync_UnstableStorage_UnstableNode(t, s)
	//})
}

func testNodeSync(t *testing.T, s StableStorage) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100
	var instability float32 = 0.25

	// create node based on mocks
	client := NewClientMock()
	storage, err := WithInstability(s, nUsers, 0.)
	require.NoError(t, err)

	for range nRuns {
		for range nRunOps {
			// apply random operation, then sync
			storage.Instability = instability
			_ = storage.RandomExternalOperation(context.TODO())
			storage.Instability = 0.
			err = nodesync.SyncState(context.TODO(), client, storage)
			require.NoError(t, err, "node sync error")
		}

		checkFullConsistency(t, client, storage)
	}
}

func testNodeSync_UnstableStorage(t *testing.T, s StableStorage) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100
	var instability float32 = 0.25

	// create node based on mocks
	client := NewClientMock()
	storage, err := WithInstability(s, nUsers, 0.)
	require.NoError(t, err)

	for range nRuns {
		storage.Instability = instability

		for range nRunOps {
			// apply random operation, then sync
			_ = storage.RandomExternalOperation(context.TODO())     // #nosec
			_ = nodesync.SyncState(context.TODO(), client, storage) // #nosec
		}

		// disable instability for one check to fix state
		storage.Instability = 0.
		err := nodesync.SyncState(context.TODO(), client, storage)
		require.NoError(t, err)
		checkFullConsistency(t, client, storage)
	}
}

func testNodeSync_UnstableStorage_UnstableNode(t *testing.T, s StableStorage) {
	nUsers := 10
	nRuns := 1000
	nRunOps := 10
	var instability float32 = 0.75

	// create node based on mocks
	client := NewUnstableClientMock()
	storage, err := WithInstability(s, nUsers, 0.)
	require.NoError(t, err)

	for range nRuns {
		storage.Instability = instability
		client.Instability = instability

		for range nRunOps {
			// apply random operation, then sync
			_ = storage.RandomExternalOperation(context.TODO())     // #nosec
			_ = nodesync.SyncState(context.TODO(), client, storage) // #nosec
		}

		// disable storage instability for one check to fix state
		storage.Instability = 0.
		err := nodesync.SyncState(context.TODO(), client, storage)
		if err != nil {
			checkStorageConsistency(t, client.BaseClient, storage)
		} else {
			checkFullConsistency(t, client.BaseClient, storage)
		}

		client.Instability = 0.
		err = nodesync.SyncState(context.TODO(), client, storage)
		require.NoError(t, err)

		checkFullConsistency(t, client.BaseClient, storage)
	}
}

func checkFullConsistency(t *testing.T, c *ClientMock, s *UnstableStorage) {

	users, err := s.s.ListUsers(context.TODO())
	require.NoError(t, err)
	target, current, err := s.GetNodeStatus(context.TODO())
	require.NoError(t, err)
	rev, err := s.s.GetNodeRev(context.TODO(), s.nodeID)
	require.NoError(t, err)
	fmt.Println("node rev", rev)
	pending, err := s.s.FindPendingSyncs(context.TODO(), s.nodeID, rev)
	require.NoError(t, err)
	for _, p := range pending {
		fmt.Println("pending user rev", p.Revision)
	}
	fmt.Println("target", target, "current", current)

	// check state is ok. only node required to be running matters
	if target != models.NodeStatusRunning {
		return
	}

	require.Equal(t, 0, len(pending),
		"pending syncs exist")
	require.Equal(t, current, target,
		"stored node state check")
	require.Equal(t, current, c.Status,
		"node state check")

	for _, su := range users {
		_, ok := c.Users[su.Profile]
		switch su.TargetStatus {
		case models.UserStatusEnabled:
			require.True(t, ok,
				"user %s (%d) check", su.Profile.Name, su.Profile.ID)
		case models.UserStatusDisabled:
			require.False(t, ok,
				"user %s (%d) check", su.Profile.Name, su.Profile.ID)
		default:
			require.Fail(t, "undefiend node user status")
		}
	}
}

func checkStorageConsistency(t *testing.T, c *ClientMock, s *UnstableStorage) {
	users, err := s.s.ListUsers(context.TODO())
	require.NoError(t, err)
	target, current, err := s.GetNodeStatus(context.TODO())
	require.NoError(t, err)
	rev, err := s.s.GetNodeRev(context.TODO(), s.nodeID)
	require.NoError(t, err)
	pending, err := s.s.FindPendingSyncs(context.TODO(), s.nodeID, rev)
	require.NoError(t, err)

	if current != models.NodeStatusRunning {
		return
	}
	require.Equal(t, current, target,
		"stored node state check")
	require.Equal(t, current, c.Status,
		"node state check")

	for _, u := range users {
		userPending := false
		for _, p := range pending {
			if p.User.Profile.ID == u.Profile.ID {
				userPending = true
				break
			}
		}

		if u.TargetStatus == models.UserStatusEnabled {
			_, ok := c.Users[u.Profile]
			require.True(t, userPending || ok, "user %s (%d) check", u.Profile.Name, u.Profile.ID)
		}
	}
}
