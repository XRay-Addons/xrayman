package poolmonitor

import (
	"context"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/poolmonitor/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestPoolMontitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := zaptest.NewLogger(t)
	defer log.Sync()

	mockSyncer := mocks.NewMockPoolSyncer(ctrl)
	expected := &models.PoolSyncResult{
		Status: models.PoolSyncOk,
		Nodes: []models.NodeSyncResult{
			{1, "node1.endpoint", nil},
			{2, "node2.endpoint", nil},
		},
	}
	mockSyncer.
		EXPECT().
		SyncNodesPool(gomock.Any()).
		DoAndReturn(func(ctx context.Context) (*models.PoolSyncResult, error) {
			time.Sleep(1 * time.Second)
			return expected, nil
		}).
		Times(2)

	poolMonitor, err := New(mockSyncer,
		WithLog(log),
		WithSyncInterval(2*time.Second))
	require.NoError(t, err)
	defer poolMonitor.Close()

	time.Sleep(500 * time.Millisecond)
	_, err = poolMonitor.Sync(context.TODO())
	require.NoError(t, err)
}
