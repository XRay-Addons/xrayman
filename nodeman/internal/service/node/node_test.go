package node

import (
	"context"
	"fmt"
	"testing"

	"github.com/XRay-Addons/xrayman/nodeman/internal/service/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNode(t *testing.T) {
	log := zaptest.NewLogger(t)

	api := NewAPIMock(log)
	storage := NewStorageMock(10, log)
	node, err := New(storage, api)
	require.NoError(t, err)

	for range 1000 {
		// enable storage external modifications and faults
		api.Unstable = true
		storage.Unstable = true

		// 100 times apply external modification and sync after that
		// sync also contains modifications, so theoretically node state
		// and storage state are unsynced after that
		for range 1000 {
			storage.ApplyExternalModifications()
			node.SyncState(context.TODO())
		}

		// last sync in stable state
		api.Unstable = false
		storage.Unstable = false

		err := node.SyncState(context.TODO())
		require.NoError(t, err)

		if storage.requiredState == models.NodeOn {
			require.Equal(t, storage.actualState, storage.requiredState)
			require.Equal(t, storage.actualState, api.status)

			for i, u := range storage.users {
				require.Equal(t, u.Required, u.Actual, fmt.Sprintf("item %d", i))
			}
		}
	}
}
