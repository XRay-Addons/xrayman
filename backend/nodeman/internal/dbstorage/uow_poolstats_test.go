package dbstorage

/*func TestUsersStats(t *testing.T) {
	logger := zaptest.NewLogger(t)

	s, _, _, cleanup := setupTestDB(t, logger)
	defer cleanup()
	logger.Info("new test db inited")

	ctx := context.Background()

	// create couple of users
	var user1, user2, user3 models.User
	err := s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		require.NoError(t, uowctx.NewUser(ctx, &user1))
		require.NoError(t, uowctx.NewUser(ctx, &user2))
		require.NoError(t, uowctx.NewUser(ctx, &user3))
		return nil
	})
	require.NoError(t, err)

	// create couple of nodes
	var node1, node2, node3 models.Node
	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		require.NoError(t, uowctx.NewNode(ctx, &node1))
		require.NoError(t, uowctx.NewNode(ctx, &node2))
		require.NoError(t, uowctx.NewNode(ctx, &node3))
		return nil
	})
	require.NoError(t, err)

	// add statistics for diffirent days
	err = s.StatsStorage().DoUoW(ctx, func(uowctx poolstats.UoWContext) error {
		// update node 1 (x2)
		upd11 := []models.UserStats{
			{ID: user1.Profile.ID, Upload: 1, Download: 2},
			{ID: user2.Profile.ID, Upload: 3, Download: 4},
		}
		require.NoError(t, uowctx.UpdateNodeStats(ctx, node1.ID, models.NodeStats{Users: upd11}))
		upd21 := []models.UserStats{
			{ID: user1.Profile.ID, Upload: 5, Download: 6},
			{ID: user3.Profile.ID, Upload: 7, Download: 8},
		}
		require.NoError(t, uowctx.UpdateNodeStats(ctx, node1.ID, models.NodeStats{Users: upd21}))

		// update node 2 (x2)
		upd12 := []models.UserStats{
			{ID: user2.Profile.ID, Upload: 9, Download: 10},
			{ID: user3.Profile.ID, Upload: 11, Download: 12},
		}
		require.NoError(t, uowctx.UpdateNodeStats(ctx, node2.ID, models.NodeStats{Users: upd12}))
		upd22 := []models.UserStats{
			{ID: user1.Profile.ID, Upload: 13, Download: 14},
		}
		require.NoError(t, uowctx.UpdateNodeStats(ctx, node2.ID, models.NodeStats{Users: upd22}))

		// update node 3 (x2)
		upd13 := []models.UserStats{
			{ID: user2.Profile.ID, Upload: 15, Download: 16},
		}
		require.NoError(t, uowctx.UpdateNodeStats(ctx, node3.ID, models.NodeStats{Users: upd13}))
		upd23 := []models.UserStats{
			{ID: user3.Profile.ID, Upload: 17, Download: 18},
		}
		require.NoError(t, uowctx.UpdateNodeStats(ctx, node3.ID, models.NodeStats{Users: upd23}))

		return nil
	})

	var poolUsers []models.User
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		poolUsers, err = uowctx.ListUsers(ctx)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, len(poolUsers), 3)
	require.Equal(t, poolUsers[0].Profile.ID, user1.Profile.ID)
	require.Equal(t, poolUsers[1].Profile.ID, user2.Profile.ID)
	require.Equal(t, poolUsers[2].Profile.ID, user3.Profile.ID)

	var userViews []models.UserView
	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) (err error) {
		userViews, err = uowctx.ListUserViews(ctx)
		return
	})
	require.NoError(t, err)
	require.Equal(t, len(userViews), 3)
	require.Equal(t, userViews[0].User.Profile.ID, user1.Profile.ID)
	require.Equal(t, int64(19), userViews[0].Traffic.Total.Upload)
	require.Equal(t, int64(22), userViews[0].Traffic.Total.Download)
	require.Equal(t, int64(19), userViews[0].Traffic.LastMonth.Upload)
	require.Equal(t, int64(22), userViews[0].Traffic.LastMonth.Download)

	require.Equal(t, userViews[1].User.Profile.ID, user2.Profile.ID)
	require.Equal(t, int64(27), userViews[1].Traffic.Total.Upload)
	require.Equal(t, int64(30), userViews[1].Traffic.Total.Download)
	require.Equal(t, int64(27), userViews[1].Traffic.LastMonth.Upload)
	require.Equal(t, int64(30), userViews[1].Traffic.LastMonth.Download)

	require.Equal(t, userViews[2].User.Profile.ID, user3.Profile.ID)
	require.Equal(t, int64(35), userViews[2].Traffic.Total.Upload)
	require.Equal(t, int64(38), userViews[2].Traffic.Total.Download)
	require.Equal(t, int64(35), userViews[2].Traffic.LastMonth.Upload)
	require.Equal(t, int64(38), userViews[2].Traffic.LastMonth.Download)
}

func TestUsersStats_HighLoad(t *testing.T) {
	realLog := zaptest.NewLogger(t)

	for _, nNodes := range []int{1, 2, 4, 8, 16, 32, 64, 128} {
		for _, nUsers := range []int{128, 512, 2048, 8192, 32768, 131072} {
			fmt.Println(nNodes, nUsers)
			nUsersUpdate := nUsers / 2
			nUpdates := 10

			logger := zap.NewNop()
			//logger := zaptest.NewLogger(t)
			logger.Info("new test log inited")

			s, _, _, cleanup := setupTestDB(t, logger)
			defer cleanup()
			logger.Info("new test db inited")

			ctx := context.Background()

			// create nodes
			createNodesBegin := time.Now()
			nodesList := make([]models.Node, nNodes, nNodes)
			err := s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
				for i, node := range nodesList {
					require.NoError(t, uowctx.NewNode(ctx, &node))
					nodesList[i] = node
				}
				return nil
			})
			require.NoError(t, err)
			addNodesTime := time.Now().Sub(createNodesBegin)

			// create users
			createUsersBegin := time.Now()
			usersList := make([]models.User, nUsers, nUsers)
			err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
				for i, user := range usersList {
					require.NoError(t, uowctx.NewUser(ctx, &user))
					usersList[i] = user
				}
				return nil
			})
			require.NoError(t, err)
			addUsersTime := time.Now().Sub(createUsersBegin)

			// updates - nUpdates times for random nUsersUpdate
			updateStatsBegin := time.Now()
			var shuffleTime time.Duration
			for _ = range nUpdates {
				err = s.StatsStorage().DoUoW(ctx, func(uowctx poolstats.UoWContext) error {
					shuffleStart := time.Now()
					update := make([]models.UserStats, nUsersUpdate, nUsersUpdate)

					rand.Shuffle(len(usersList), func(i, j int) {
						usersList[i], usersList[j] = usersList[j], usersList[i]
					})

					for i := range nUsersUpdate {
						update[i] = models.UserStats{
							ID:       usersList[i].Profile.ID,
							Download: 2,
							Upload:   1,
						}
					}
					sort.Slice(update, func(i, j int) bool {
						return update[i].ID < update[j].ID
					})
					shuffleTime += time.Now().Sub(shuffleStart)

					nodeID := nodesList[rand.IntN(len(nodesList))].ID
					require.NoError(t, uowctx.UpdateNodeStats(ctx, nodeID, models.NodeStats{Users: update}))
					return nil
				})
				require.NoError(t, err)
			}
			updateTime := (time.Now().Sub(updateStatsBegin) - shuffleTime) / time.Duration(nUpdates)

			// get result and check it
			var userViewsList []models.UserView
			listUsersBegin := time.Now()
			err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
				userViewsList, err = uowctx.ListUserViews(ctx)
				require.NoError(t, err)
				return nil
			})
			require.NoError(t, err)
			listUsersTime := time.Now().Sub(listUsersBegin)

			var totalUpload, totalDownload int64
			for _, u := range userViewsList {
				totalUpload += u.Traffic.Total.Upload
				totalDownload += u.Traffic.Total.Download
			}
			require.Equal(t, int64(nUpdates*nUsersUpdate*2), totalDownload)
			require.Equal(t, int64(nUpdates*nUsersUpdate), totalUpload)
			require.NoError(t, err)

			realLog.Info("test stats:",
				zap.Int("nNodes", nNodes),
				zap.Int("nUsers", nUsers),
				zap.Duration("addUsersTime", addUsersTime),
				zap.Duration("addNodesTime", addNodesTime),
				zap.Duration("updateTime", updateTime),
				zap.Duration("listUsersTime", listUsersTime))
		}
	}
}

func TestUsersStats_Explain(t *testing.T) {
	realLog := zaptest.NewLogger(t)

	nNodes := 100
	nUsers := 100000
	nUsersUpdate := nUsers / 2
	nUpdates := 10

	logger := zap.NewNop()
	//logger := zaptest.NewLogger(t)
	logger.Info("new test log inited")

	s, _, _, cleanup := setupTestDB(t, logger)
	defer cleanup()
	logger.Info("new test db inited")

	ctx := context.Background()

	// create nodes
	createNodesBegin := time.Now()
	nodesList := make([]models.Node, nNodes, nNodes) // 100000
	err := s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		for i, node := range nodesList {
			require.NoError(t, uowctx.NewNode(ctx, &node))
			nodesList[i] = node
		}
		return nil
	})
	require.NoError(t, err)
	addNodesTime := time.Now().Sub(createNodesBegin)

	// create users
	createUsersBegin := time.Now()
	usersList := make([]models.User, nUsers, nUsers) // 100000
	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		for i, user := range usersList {
			require.NoError(t, uowctx.NewUser(ctx, &user))
			usersList[i] = user
		}
		return nil
	})
	require.NoError(t, err)
	addUsersTime := time.Now().Sub(createUsersBegin)

	// updates - nUpdates times for random nUsersUpdate
	updateStatsBegin := time.Now()
	var shuffleTime time.Duration
	for i := range nUpdates {
		if i == 0 {
			s.explainLog = realLog
		} else {
			s.explainLog = nil
		}
		err = s.StatsStorage().DoUoW(ctx, func(uowctx poolstats.UoWContext) error {
			shuffleStart := time.Now()
			update := make([]models.UserStats, nUsersUpdate, nUsersUpdate)

			rand.Shuffle(len(usersList), func(i, j int) {
				usersList[i], usersList[j] = usersList[j], usersList[i]
			})

			for i := range nUsersUpdate {
				update[i] = models.UserStats{
					ID:       usersList[i].Profile.ID,
					Download: 2,
					Upload:   1,
				}
			}
			sort.Slice(update, func(i, j int) bool {
				return update[i].ID < update[j].ID
			})
			shuffleTime += time.Now().Sub(shuffleStart)

			nodeID := nodesList[rand.IntN(len(nodesList))].ID
			require.NoError(t, uowctx.UpdateNodeStats(ctx, nodeID, models.NodeStats{Users: update}))
			return nil
		})
		require.NoError(t, err)
	}
	updateTime := (time.Now().Sub(updateStatsBegin) - shuffleTime) / time.Duration(nUpdates)

	// get result and check it
	listUsersBegin := time.Now()
	s.explainLog = realLog
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		_, err = uowctx.ListUsers(ctx)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
	listUsersTime := time.Now().Sub(listUsersBegin)

	// get result and check it
	var userViewsList []models.UserView
	listUsersViewBegin := time.Now()
	s.explainLog = realLog
	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		userViewsList, err = uowctx.ListUserViews(ctx)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
	listUsersViewTime := time.Now().Sub(listUsersViewBegin)

	var totalUpload, totalDownload int64
	for _, u := range userViewsList {
		totalUpload += u.Traffic.Total.Upload
		totalDownload += u.Traffic.Total.Download
	}
	require.Equal(t, int64(nUpdates*nUsersUpdate*2), totalDownload)
	require.Equal(t, int64(nUpdates*nUsersUpdate), totalUpload)
	require.NoError(t, err)

	realLog.Info("test stats:",
		zap.Int("nNodes", nNodes),
		zap.Int("nUsers", nUsers),
		zap.Duration("addUsersTime", addUsersTime),
		zap.Duration("addNodesTime", addNodesTime),
		zap.Duration("updateTime", updateTime),
		zap.Duration("listUsersTime", listUsersTime),
		zap.Duration("listUsersViewTime", listUsersViewTime))
}*/
