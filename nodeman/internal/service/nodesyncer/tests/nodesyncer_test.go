package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodesyncer"
	"github.com/stretchr/testify/require"
)

func TestNodeSyncer(t *testing.T) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100

	// create node based on mocks
	client := NewClientMock()
	storage := NewStorageMock(nUsers)
	node, err := nodesyncer.New(storage, client)
	require.NoError(t, err)

	for range nRuns {
		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			node.SyncState(context.TODO())
		}

		// check state is ok. only node required to be running matters
		if storage.TargetState != models.NodeStatusRunning {
			continue
		}

		require.Equal(t, storage.CurrentState, storage.TargetState, "stored node state check")
		require.Equal(t, storage.CurrentState, client.Status, "node state check")

		for _, u := range storage.Users {
			require.Equal(t, u.TargetStatus, u.CurrentStatus, fmt.Sprintf("user %d check", u.User.ID))
		}
	}
}

func TestNode_UnstableStorage(t *testing.T) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100
	var instability float32 = 0.25

	// create node based on mocks
	client := NewClientMock()
	storage := NewUnstableStorageMock(nUsers)
	node, err := nodesyncer.New(storage, client)
	require.NoError(t, err)

	for range nRuns {
		storage.Instability = instability

		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			node.SyncState(context.TODO())
		}

		// disable instability for one check to fix state
		storage.Instability = 0.
		node.SyncState(context.TODO())

		baseStorage := storage.BaseStorage
		// check state is ok. only node required to be running matters
		if baseStorage.TargetState != models.NodeStatusRunning {
			continue
		}

		require.Equal(t, baseStorage.CurrentState, baseStorage.TargetState, "stored node state check")
		require.Equal(t, baseStorage.CurrentState, client.Status, "node state check")

		for _, u := range baseStorage.Users {
			require.Equal(t, u.TargetStatus, u.CurrentStatus, fmt.Sprintf("user %d check", u.User.ID))
		}
	}
}

func TestNode_UnstableStorage_UnstableNode(t *testing.T) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100
	var instability float32 = 0.25

	// create node based on mocks
	client := NewUnstableClientMock()
	storage := NewUnstableStorageMock(nUsers)
	node, err := nodesyncer.New(storage, client)
	require.NoError(t, err)

	for range nRuns {
		storage.Instability = instability
		client.Instability = instability

		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			node.SyncState(context.TODO())
		}

		// disable instability for one check to fix state
		storage.Instability = 0.
		client.Instability = 0.
		node.SyncState(context.TODO())

		baseStorage := storage.BaseStorage
		baseAPI := client.BaseAPI
		// check state is ok. only node required to be running matters
		if baseStorage.TargetState != models.NodeStatusRunning {
			continue
		}

		require.Equal(t, baseStorage.CurrentState, baseStorage.TargetState, "stored node state check")
		require.Equal(t, baseStorage.CurrentState, baseAPI.Status, "node state check")

		for _, u := range baseStorage.Users {
			require.Equal(t, u.TargetStatus, u.CurrentStatus, fmt.Sprintf("user %d check", u.User.ID))
		}
	}
}
