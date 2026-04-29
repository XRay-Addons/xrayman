package syncman

/*func TestSyncService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := zaptest.NewLogger(t)
	defer log.Sync()

	mockSyncer := mocks.NewMockPoolSyncer(ctrl)
	expected := []models.NodeSyncResult{
		{1, "node1.endpoint", nil},
		{2, "node2.endpoint", nil},
	}
	mockSyncer.
		EXPECT().
		SyncPoolState(gomock.Any()).
		DoAndReturn(func(ctx context.Context) ([]models.NodeSyncResult, error) {
			time.Sleep(1 * time.Second)
			return expected, nil
		}).
		Times(2)

	poolMonitor, err := New(mockSyncer,
		WithLogger(log),
		WithSyncInterval(2*time.Second))
	require.NoError(t, err)
	defer poolMonitor.Close()

	time.Sleep(500 * time.Millisecond)
	_, err = poolMonitor.SyncNodesPool(context.TODO())
	require.NoError(t, err)
}*/
