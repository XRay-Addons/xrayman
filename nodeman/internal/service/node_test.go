package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func NewTestLogger(t *testing.T) *zap.Logger {
	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.TimeKey = ""
	encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encCfg)

	writer := zapcore.AddSync(zaptest.NewTestingWriter(t))
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		zap.DebugLevel,
	)

	logger := zap.New(core)
	t.Cleanup(func() { _ = logger.Sync() })
	return logger
}

func TestNode(t *testing.T) {

	log := NewTestLogger(t)

	nUsers := 10
	api := NewNodeAPIEmulator(log)
	storage := NewNodeStorageEmulator(nUsers, log)
	storage.requiredState = NodeRunning

	node, err := New(api, storage, log)
	require.NoError(t, err)

	for range 1000 {
		// enable storage external modifications and faults
		api.unstable = true
		storage.unstable = true

		// 100 times apply external modification and sync after that
		// sync also contains modifications, so theoretically node state
		// and storage state are unsynced after that
		for range 2 {

			log.Info("new iteration:")
			log.Sugar().Infof("storage: actual: %v required %v pending %d; node: %v",
				storage.actualState, storage.requiredState, len(storage.pendingUsers),
				api.status)

			storage.applyExternalModifications()
			if err := node.SyncNodeStatus(context.TODO()); err != nil {
				log.Error("Sync error", zap.Error(err))
			} else {
				log.Info("Sync done")
			}
			log.Info("")
		}

		// last sync in stable state
		api.unstable = false
		storage.unstable = false

		node.SyncNodeStatus(context.TODO())

		log.Info("----------------- End of test iteration --------------------")
		checkNodeState(t, storage, api, log)
	}
}

func checkNodeState(t *testing.T,
	storage *NodeStorageEmulator,
	node *NodeAPIEmulator,
	log *zap.Logger,
) {
	log.Info("    Check node state:")
	log.Info("  Storage:")
	log.Sugar().Infof("actual: %v, required: %v",
		storage.actualState, storage.requiredState)
	for _, u := range storage.users {
		logStr := fmt.Sprintf("user %d: %v", u.User.ID, u.Status)
		if p, exists := storage.pendingUsers[u.User]; exists {
			logStr += fmt.Sprintf(" pending: %v", p)
		}
		log.Info(logStr)
	}
	log.Info("\n")

	log.Info("  API:")
	log.Sugar().Infof("status: %v", node.status)
	log.Sugar().Infof("users: ")
	for u, _ := range node.users {
		log.Sugar().Infof("user: %v", u)
	}
	log.Info("\n")
	defer func() { log.Info("---- end of test check ----\n") }()

	if storage.requiredState == NodeStopped {
		// node is required to be stopped. it't ok if its
		// check storage
		require.NotEqual(t, NodeRunning, storage.actualState, "actual storage state")
		switch storage.actualState {
		case NodeStopped:
			require.Equal(t, 0, len(storage.pendingUsers), "pending users")
			require.NotEqual(t, NodeRunning, node.status, "node status")
			require.Equal(t, 0, len(node.users), "node users count")
		case NodeRunning:
			require.Fail(t, "node actual state should not be 'running'")
		}
	} else {
		// check storage
		require.Equal(t, NodeRunning, storage.actualState, "actual storage state")
		require.Equal(t, 0, len(storage.pendingUsers), "pending users")
		// check node
		require.Equal(t, NodeRunning, node.status, "node status")
		storageUsers := 0
		for _, su := range storage.users {
			if su.Status == UserEnabled {
				storageUsers++
			}
		}
		require.Equal(t, storageUsers, len(node.users), "node users count")
	}
}
