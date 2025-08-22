package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/node"
	"github.com/stretchr/testify/require"
)

func TestNodeSyncer(t *testing.T) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100

	// create node based on mocks
	client := NewClientMock()
	storage := NewStorageMock(nUsers)
	uow, err := storage.NewUoW()
	require.NoError(t, err)
	syncer := node.NewSyncer()

	for range nRuns {
		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			syncer.SyncNodeState(context.TODO(), client, uow)
		}

		// check state is ok. only node required to be running matters
		if storage.TargetStatus != models.NodeStatusRunning {
			continue
		}

		require.Equal(t, storage.CurrentStatus, storage.TargetStatus, "stored node state check")
		require.Equal(t, storage.CurrentStatus, client.Status, "node state check")

		for i, u := range storage.Users {
			require.Equal(t, u.TargetStatus, storage.CurrentUserStatus[i],
				fmt.Sprintf("user %s check", u.Profile.Name))
		}
	}
}

func TestNodeSyncer_UnstableStorage(t *testing.T) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100
	var instability float32 = 0.25

	// create node based on mocks
	client := NewClientMock()
	storage := NewUnstableStorageMock(nUsers)
	uow, err := storage.NewUoW()
	require.NoError(t, err)
	syncer := node.NewSyncer()

	for range nRuns {
		storage.Instability = instability

		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			syncer.SyncNodeState(context.TODO(), client, uow)
		}

		// disable instability for one check to fix state
		storage.Instability = 0.
		syncer.SyncNodeState(context.TODO(), client, uow)

		baseStorage := storage.BaseStorage
		// check state is ok. only node required to be running matters
		if baseStorage.TargetStatus != models.NodeStatusRunning {
			continue
		}

		require.Equal(t, baseStorage.CurrentStatus, baseStorage.TargetStatus,
			"stored node state check")
		require.Equal(t, baseStorage.CurrentStatus, client.Status,
			"node state check")

		for i, u := range baseStorage.Users {
			require.Equal(t, u.TargetStatus, baseStorage.CurrentUserStatus[i],
				fmt.Sprintf("user %s check", u.Profile.Name))
		}
	}
}

func TestNodeSyncer_UnstableStorage_UnstableNode(t *testing.T) {
	nUsers := 10
	nRuns := 100
	nRunOps := 100
	var instability float32 = 0.25

	// create node based on mocks
	client := NewUnstableClientMock()
	storage := NewUnstableStorageMock(nUsers)
	uow, err := storage.NewUoW()
	require.NoError(t, err)
	syncer := node.NewSyncer()

	for range nRuns {
		storage.Instability = instability
		client.Instability = instability

		for range nRunOps {
			// apply random operation, then sync
			storage.RandomExternalOperation()
			syncer.SyncNodeState(context.TODO(), client, uow)
		}

		// disable instability for one check to fix state
		storage.Instability = 0.
		client.Instability = 0.
		syncer.SyncNodeState(context.TODO(), client, uow)

		baseStorage := storage.BaseStorage
		baseClient := client.BaseClient
		// check state is ok. only node required to be running matters
		if baseStorage.TargetStatus != models.NodeStatusRunning {
			continue
		}

		require.Equal(t, baseStorage.CurrentStatus, baseStorage.TargetStatus,
			"stored node state check")
		require.Equal(t, baseStorage.CurrentStatus, baseClient.Status,
			"node state check")

		for i, u := range baseStorage.Users {
			require.Equal(t, u.TargetStatus, baseStorage.CurrentUserStatus[i],
				fmt.Sprintf("user %s check", u.Profile.Name))
		}
	}
}
