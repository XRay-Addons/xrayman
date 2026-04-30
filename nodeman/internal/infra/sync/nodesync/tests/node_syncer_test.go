package tests

import (
	"context"
	"testing"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/nodesync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/stretchr/testify/require"
)

func Testnodesync(t *testing.T) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100

	// create node based on mocks
	client := NewClientMock()
	storage := NewStorageMock(nUsers)

	for range nRuns {
		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			_ = nodesync.SyncState(context.TODO(), client, storage)
		}

		checkFullConsistency(t, client, storage)
	}
}

func Testnodesync_UnstableStorage(t *testing.T) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100
	var instability float32 = 0.25

	// create node based on mocks
	client := NewClientMock()
	storage := NewUnstableStorageMock(nUsers)

	for range nRuns {
		storage.Instability = instability

		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			_ = nodesync.SyncState(context.TODO(), client, storage) // #nosec
		}

		// disable instability for one check to fix state
		storage.Instability = 0.
		err := nodesync.SyncState(context.TODO(), client, storage)
		require.NoError(t, err)

		checkFullConsistency(t, client, storage.BaseStorage)
	}
}

func Testnodesync_UnstableStorage_UnstableNode(t *testing.T) {
	nUsers := 10
	nRuns := 1000
	nRunOps := 10
	var instability float32 = 0.75

	// create node based on mocks
	client := NewUnstableClientMock()
	storage := NewUnstableStorageMock(nUsers)

	for range nRuns {
		storage.Instability = instability
		client.Instability = instability

		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			_ = nodesync.SyncState(context.TODO(), client, storage) // #nosec
		}

		// disable storage instability for one check to fix state
		storage.Instability = 0.
		err := nodesync.SyncState(context.TODO(), client, storage)
		if err != nil {
			checkStorageConsistency(t, client.BaseClient, storage.BaseStorage)
		} else {
			checkFullConsistency(t, client.BaseClient, storage.BaseStorage)
		}

		client.Instability = 0.
		err = nodesync.SyncState(context.TODO(), client, storage)
		require.NoError(t, err)

		checkFullConsistency(t, client.BaseClient, storage.BaseStorage)
	}
}

func checkFullConsistency(t *testing.T, c *ClientMock, s *StorageMock) {
	// check state is ok. only node required to be running matters
	if s.TargetStatus != models.NodeStatusRunning {
		return
	}

	require.Equal(t, s.CurrentStatus, s.TargetStatus,
		"stored node state check")
	require.Equal(t, s.CurrentStatus, c.Status,
		"node state check")

	for i, u := range s.Users {
		require.Equal(t, u.TargetStatus, s.CurrentUserStatus[i],
			"user %s check", u.Profile.Name)
	}
}

func checkStorageConsistency(t *testing.T, c *ClientMock, s *StorageMock) {
	if s.CurrentStatus != models.NodeStatusRunning {
		return
	}
	require.Equal(t, s.CurrentStatus, s.TargetStatus,
		"stored node state check")
	require.Equal(t, s.CurrentStatus, c.Status,
		"node state check")

	for i, u := range s.Users {
		if s.CurrentUserStatus[i] == models.UserStatusEnabled {
			_, ok := c.Users[u.Profile]
			require.True(t, ok, "user %s check", u.Profile.Name)
		}
	}
}
