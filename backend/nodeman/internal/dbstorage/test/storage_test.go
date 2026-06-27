package dbstoragetest

import (
	"context"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////
// Storage test
func TestStorage_Nodes(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	s, _ := setupTestDB(t, logger)
	logger.Info("new test db inited")

	node1 := models.Node{
		CurrentStatus: models.NodeStatusUnknown,
		TargetStatus:  models.NodeStatusRunning,
	}
	node2 := models.Node{
		CurrentStatus: models.NodeStatusRunning,
		TargetStatus:  models.NodeStatusUnknown,
	}
	node3 := models.Node{
		CurrentStatus: models.NodeStatusStopped,
		TargetStatus:  models.NodeStatusUnknown,
	}
	err := s.DoTx(ctx, func(ctx context.Context) error {
		for _, node := range []*models.Node{&node1, &node2, &node3} {
			if err := s.NewNode(ctx, node); err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	err = s.DeleteNode(ctx, node2.ID)
	require.NoError(t, err)

	err = s.SetTargetNodeStatus(ctx, node1.ID, models.NodeStatusStopped)
	require.NoError(t, err)

	err = s.SetCurrentNodeStatus(ctx, node3.ID, models.NodeStatusRunning)
	require.NoError(t, err)

	nodesList, err := s.ListNodes(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(nodesList))
	require.Equal(t, models.NodeStatusStopped, nodesList[0].TargetStatus)
	require.Equal(t, models.NodeStatusUnknown, nodesList[0].CurrentStatus)
	require.Equal(t, models.NodeStatusUnknown, nodesList[1].TargetStatus)
	require.Equal(t, models.NodeStatusRunning, nodesList[1].CurrentStatus)

	existedNode, err := s.GetNode(ctx, node1.ID)
	require.NoError(t, err)
	require.Equal(t, node1.ID, existedNode.ID)

	_, err = s.GetNode(ctx, node2.ID)
	require.ErrorIs(t, err, errdefs.ErrNotFound)
}

func TestStorage_Users(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	s, _ := setupTestDB(t, logger)
	logger.Info("new test db inited")

	user1 := models.User{
		TargetStatus: models.UserStatusEnabled,
	}
	user2 := models.User{
		TargetStatus: models.UserStatusDisabled,
	}
	user3 := models.User{
		TargetStatus: models.UserStatusUnknown,
	}
	err := s.DoTx(ctx, func(ctx context.Context) error {
		for _, user := range []*models.User{&user1, &user2, &user3} {
			if err := s.NewUser(ctx, user); err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	err = s.DeleteUser(ctx, user2.Profile.ID)
	require.NoError(t, err)

	err = s.SetTargetUserStatus(ctx, user1.Profile.ID, models.UserStatusDisabled)
	require.NoError(t, err)

	usersList, err := s.ListUsers(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(usersList))
	require.Equal(t, models.UserStatusDisabled, usersList[0].TargetStatus)
	require.Equal(t, models.UserStatusUnknown, usersList[1].TargetStatus)

	existedUser, err := s.GetUserView(ctx, user1.Profile.ID, user1.Profile.Name)
	require.NoError(t, err)
	require.Equal(t, user1.Profile.ID, existedUser.User.Profile.ID)

	_, err = s.GetUserView(ctx, user2.Profile.ID, user2.Profile.Name)
	require.ErrorIs(t, err, errdefs.ErrNotFound)

	err = s.DoTx(ctx, func(ctx context.Context) error {
		_, err = s.GetUserView(ctx, user1.Profile.ID, "fake name")
		return err
	})
	require.ErrorIs(t, err, errdefs.ErrNotFound)
}

func TestStorage_UserNodes(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	s, _ := setupTestDB(t, logger)
	logger.Info("new test db inited")

	user1 := models.User{
		TargetStatus: models.UserStatusEnabled,
	}
	user2 := models.User{
		TargetStatus: models.UserStatusDisabled,
	}
	user3 := models.User{
		TargetStatus: models.UserStatusDisabled,
	}

	err := s.DoTx(ctx, func(ctx context.Context) error {
		for _, user := range []*models.User{&user1, &user2, &user3} {
			if err := s.NewUser(ctx, user); err != nil {
				return err
			}
		}
		return nil
	})

	require.NoError(t, err)
	node1 := models.Node{
		CurrentStatus: models.NodeStatusRunning,
		TargetStatus:  models.NodeStatusRunning,
	}

	err = s.NewNode(ctx, &node1)
	require.NoError(t, err)

	pendingSyncs, err := s.FindPendingSyncs(ctx, node1.ID)
	require.NoError(t, err)
	require.Equal(t, 1, len(pendingSyncs))
	require.Equal(t, user1.Profile.ID, pendingSyncs[0].User.Profile.ID)
	require.Equal(t, models.UserStatusDisabled, pendingSyncs[0].CurrentStatus)
	require.Equal(t, models.UserStatusEnabled, pendingSyncs[0].User.TargetStatus)

	syncsPath := []models.UserStatusPatch{
		{UserID: user1.Profile.ID, Status: models.UserStatusEnabled},
		{UserID: user2.Profile.ID, Status: models.UserStatusEnabled},
	}
	err = s.UpdateNodeUsers(ctx, node1.ID, syncsPath)
	require.NoError(t, err)

	pendingSyncs, err = s.FindPendingSyncs(ctx, node1.ID)
	require.NoError(t, err)
	require.Equal(t, 1, len(pendingSyncs))
	require.Equal(t, user2.Profile.ID, pendingSyncs[0].User.Profile.ID)
	require.Equal(t, models.UserStatusEnabled, pendingSyncs[0].CurrentStatus)
	require.Equal(t, models.UserStatusDisabled, pendingSyncs[0].User.TargetStatus)

	syncsPath = []models.UserStatusPatch{
		{UserID: user2.Profile.ID, Status: models.UserStatusEnabled},
	}
	err = s.SetNodeUsers(ctx, node1.ID, syncsPath)
	require.NoError(t, err)
	pendingSyncs, err = s.FindPendingSyncs(ctx, node1.ID)
	require.NoError(t, err)
	require.Equal(t, 2, len(pendingSyncs))
}

func TestStorage_Stats(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	s, _ := setupTestDB(t, logger)
	logger.Info("new test db inited")

	user1 := models.User{
		TargetStatus: models.UserStatusEnabled,
	}
	user2 := models.User{
		TargetStatus: models.UserStatusEnabled,
	}

	err := s.DoTx(ctx, func(ctx context.Context) error {
		for _, user := range []*models.User{&user1, &user2} {
			if err := s.NewUser(ctx, user); err != nil {
				return err
			}
		}
		return nil
	})

	require.NoError(t, err)
	node1 := models.Node{
		CurrentStatus: models.NodeStatusRunning,
		TargetStatus:  models.NodeStatusRunning,
	}
	node2 := models.Node{
		CurrentStatus: models.NodeStatusRunning,
		TargetStatus:  models.NodeStatusRunning,
	}

	err = s.DoTx(ctx, func(ctx context.Context) error {
		for _, node := range []*models.Node{&node1, &node2} {
			if err := s.NewNode(ctx, node); err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	// update stats
	err = s.DoTx(ctx, func(ctx context.Context) error {
		if err := s.UpdateNodeStats(ctx, node1.ID, models.NodeStats{
			Users: []models.UserStats{
				{ID: user1.Profile.ID, Uplink: 1, Downlink: 2},
				{ID: user2.Profile.ID, Uplink: 3, Downlink: 4},
			},
		}); err != nil {
			return err
		}
		if err := s.UpdateNodeStats(ctx, node2.ID, models.NodeStats{
			Users: []models.UserStats{
				{ID: user1.Profile.ID, Uplink: 5, Downlink: 6},
			},
		}); err != nil {
			return err
		}
		if err := s.UpdateDailyStats(ctx, time.Now().Add(-60*24*time.Hour)); err != nil {
			return err
		}
		if err := s.UpdateNodeStats(ctx, node2.ID, models.NodeStats{
			Users: []models.UserStats{
				{ID: user2.Profile.ID, Uplink: 11, Downlink: 12},
			},
		}); err != nil {
			return err
		}
		return nil
	})
	require.NoError(t, err)

	usersList, err := s.ListUserViews(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(usersList))
	require.Equal(t, int64(6), usersList[0].Traffic.Total.Upload)
	require.Equal(t, int64(8), usersList[0].Traffic.Total.Download)
	require.Equal(t, int64(14), usersList[1].Traffic.Total.Upload)
	require.Equal(t, int64(16), usersList[1].Traffic.Total.Download)

	require.Equal(t, int64(0), usersList[0].Traffic.LastMonth.Upload)
	require.Equal(t, int64(0), usersList[0].Traffic.LastMonth.Download)
	require.Equal(t, int64(11), usersList[1].Traffic.LastMonth.Upload)
	require.Equal(t, int64(12), usersList[1].Traffic.LastMonth.Download)

	userView, err := s.GetUserView(ctx, user1.Profile.ID, user1.Profile.Name)
	require.NoError(t, err)
	require.Equal(t, int64(6), userView.Traffic.Total.Upload)
	require.Equal(t, int64(8), userView.Traffic.Total.Download)
	require.Equal(t, int64(0), userView.Traffic.LastMonth.Upload)
	require.Equal(t, int64(0), userView.Traffic.LastMonth.Download)

	userView, err = s.GetUserView(ctx, user2.Profile.ID, user2.Profile.Name)
	require.NoError(t, err)
	require.Equal(t, int64(14), userView.Traffic.Total.Upload)
	require.Equal(t, int64(16), userView.Traffic.Total.Download)
	require.Equal(t, int64(11), userView.Traffic.LastMonth.Upload)
	require.Equal(t, int64(12), userView.Traffic.LastMonth.Download)
}

func TestStorage_Password(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	s, _ := setupTestDB(t, logger)
	logger.Info("new test db inited")

	auth := models.Auth{
		PasswordHash: []byte("hash"),
	}
	err := s.SetAuth(ctx, &auth)
	require.NoError(t, err)

	dbauth, err := s.GetAuth(ctx)
	require.NoError(t, err)
	require.Equal(t, auth, *dbauth)
}

func TestStorage_DynConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	s, _ := setupTestDB(t, logger)
	logger.Info("new test db inited")

	cfg := models.DynamicConfig{
		UserPage:     "user page",
		UsersMessage: "users message",
		TgPage:       "tg page",
	}
	err := s.SetDynamicConfig(ctx, cfg)
	require.NoError(t, err)

	readCfg, err := s.GetDynamicConfig(ctx)
	require.NoError(t, err)
	require.Equal(t, cfg, *readCfg)

	cfg2 := models.DynamicConfig{
		UserPage:     "user page2",
		UsersMessage: "users message2",
		TgPage:       "tg page2",
	}
	err = s.EnsureDynamicConfig(ctx, cfg2)
	require.NoError(t, err)

	readCfg, err = s.GetDynamicConfig(ctx)
	require.NoError(t, err)
	require.Equal(t, cfg, *readCfg)
}

func TestStorage_Tx(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	s, _ := setupTestDB(t, logger)
	logger.Info("new test db inited")

	node1 := models.Node{
		CurrentStatus: models.NodeStatusRunning,
		TargetStatus:  models.NodeStatusRunning,
	}
	node2 := models.Node{
		CurrentStatus: models.NodeStatusRunning,
		TargetStatus:  models.NodeStatusRunning,
	}
	node3 := models.Node{
		CurrentStatus: models.NodeStatusRunning,
		TargetStatus:  models.NodeStatusRunning,
	}

	err := s.DoTx(ctx, func(ctx context.Context) error {
		// no error
		if err := s.NewNode(ctx, &node1); err != nil {
			return err
		}

		// nested tx - add node and than error op - so rollback
		err := s.DoTx(ctx, func(ctx context.Context) error {
			// no error
			if err := s.NewNode(ctx, &node2); err != nil {
				return err
			}
			// with error
			if _, err := s.GetNode(ctx, models.NodeID(-1)); err != nil {
				return err
			}
			return nil
		})
		require.Error(t, err)

		// nnnested tx - add node and than error op - so rollback
		err = s.DoTx(ctx, func(ctx context.Context) error {
			// no error
			if err := s.NewNode(ctx, &node2); err != nil {
				return err
			}
			return s.DoTx(ctx, func(ctx context.Context) error {
				// no error
				if err := s.NewNode(ctx, &node3); err != nil {
					return err
				}
				// with error
				if _, err := s.GetNode(ctx, models.NodeID(-1)); err != nil {
					return err
				}
				return nil
			})
		})
		require.Error(t, err)

		err = s.DoTx(ctx, func(ctx context.Context) error {
			// success tttx
			_ = s.DoTx(ctx, func(ctx context.Context) error {
				if err := s.NewNode(ctx, &node3); err != nil {
					return err
				}
				return nil
			})
			// failed tttx
			_ = s.DoTx(ctx, func(ctx context.Context) error {
				if err := s.NewNode(ctx, &node3); err != nil {
					return err
				}
				// with error
				if _, err := s.GetNode(ctx, models.NodeID(-1)); err != nil {
					return err
				}
				return nil
			})
			// success tttx
			_ = s.DoTx(ctx, func(ctx context.Context) error {
				if err := s.NewNode(ctx, &node3); err != nil {
					return err
				}
				return nil
			})
			return nil
		})
		require.NoError(t, err)

		return nil
	})
	require.NoError(t, err)

	nodes, err := s.ListNodes(ctx)
	require.NoError(t, err)
	require.Equal(t, 4, len(nodes))
}
