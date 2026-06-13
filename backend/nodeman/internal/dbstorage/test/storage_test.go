package dbstoragetest

import (
	"context"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/auth/password"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/stats/poolstats"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
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
	err := s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		for _, node := range []*models.Node{&node1, &node2, &node3} {
			if err := uowctx.NewNode(ctx, node); err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		return uowctx.DeleteNode(ctx, node2.ID)
	})
	require.NoError(t, err)

	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		return uowctx.SetTargetNodeStatus(ctx, node1.ID, models.NodeStatusStopped)
	})
	require.NoError(t, err)

	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		return uowctx.SetCurrentNodeStatus(ctx, node3.ID, models.NodeStatusRunning)
	})
	require.NoError(t, err)

	var nodesList []models.Node
	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		nodesList, err = uowctx.ListNodes(ctx)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(nodesList))
	require.Equal(t, models.NodeStatusStopped, nodesList[0].TargetStatus)
	require.Equal(t, models.NodeStatusUnknown, nodesList[0].CurrentStatus)
	require.Equal(t, models.NodeStatusUnknown, nodesList[1].TargetStatus)
	require.Equal(t, models.NodeStatusRunning, nodesList[1].CurrentStatus)

	var existedNode *models.Node
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		existedNode, err = uowctx.GetNode(ctx, node1.ID)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, node1.ID, existedNode.ID)

	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		_, err = uowctx.GetNode(ctx, node2.ID)
		return err
	})
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
	err := s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		for _, user := range []*models.User{&user1, &user2, &user3} {
			if err := uowctx.NewUser(ctx, user); err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		return uowctx.DeleteUser(ctx, user2.Profile.ID)
	})
	require.NoError(t, err)

	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		return uowctx.SetTargetUserStatus(ctx, user1.Profile.ID, models.UserStatusDisabled)
	})
	require.NoError(t, err)

	var usersList []models.User
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		usersList, err = uowctx.ListUsers(ctx)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(usersList))
	require.Equal(t, models.UserStatusDisabled, usersList[0].TargetStatus)
	require.Equal(t, models.UserStatusUnknown, usersList[1].TargetStatus)

	var existedUser *models.UserView
	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		existedUser, err = uowctx.GetUserView(ctx, user1.Profile.ID, user1.Profile.Name)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, user1.Profile.ID, existedUser.User.Profile.ID)

	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		_, err = uowctx.GetUserView(ctx, user2.Profile.ID, user2.Profile.Name)
		return err
	})
	require.ErrorIs(t, err, errdefs.ErrNotFound)

	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		_, err = uowctx.GetUserView(ctx, user1.Profile.ID, "fake name")
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

	err := s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		for _, user := range []*models.User{&user1, &user2, &user3} {
			if err := uowctx.NewUser(ctx, user); err != nil {
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

	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		return uowctx.NewNode(ctx, &node1)
	})
	require.NoError(t, err)

	var pendingSyncs []models.UserSyncStatus
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		pendingSyncs, err = uowctx.FindPendingSyncs(ctx, node1.ID)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(pendingSyncs))
	require.Equal(t, user1.Profile.ID, pendingSyncs[0].User.Profile.ID)
	require.Equal(t, models.UserStatusDisabled, pendingSyncs[0].CurrentStatus)
	require.Equal(t, models.UserStatusEnabled, pendingSyncs[0].User.TargetStatus)

	syncsPath := []models.UserStatusPatch{
		{UserID: user1.Profile.ID, Status: models.UserStatusEnabled},
		{UserID: user2.Profile.ID, Status: models.UserStatusEnabled},
	}
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		return uowctx.UpdateNodeUsers(ctx, node1.ID, syncsPath)
	})
	require.NoError(t, err)

	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		pendingSyncs, err = uowctx.FindPendingSyncs(ctx, node1.ID)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(pendingSyncs))
	require.Equal(t, user2.Profile.ID, pendingSyncs[0].User.Profile.ID)
	require.Equal(t, models.UserStatusEnabled, pendingSyncs[0].CurrentStatus)
	require.Equal(t, models.UserStatusDisabled, pendingSyncs[0].User.TargetStatus)

	syncsPath = []models.UserStatusPatch{
		{UserID: user2.Profile.ID, Status: models.UserStatusEnabled},
	}
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		return uowctx.SetNodeUsers(ctx, node1.ID, syncsPath)
	})
	require.NoError(t, err)
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		pendingSyncs, err = uowctx.FindPendingSyncs(ctx, node1.ID)
		return err
	})
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

	err := s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		for _, user := range []*models.User{&user1, &user2} {
			if err := uowctx.NewUser(ctx, user); err != nil {
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

	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		for _, node := range []*models.Node{&node1, &node2} {
			if err := uowctx.NewNode(ctx, node); err != nil {
				return err
			}
		}
		return nil
	})
	require.NoError(t, err)

	// update stats
	err = s.StatsStorage().DoUoW(ctx, func(uowctx poolstats.UoWContext) error {
		if err := uowctx.UpdateNodeStats(ctx, node1.ID, models.NodeStats{
			Users: []models.UserStats{
				{ID: user1.Profile.ID, Upload: 1, Download: 2},
				{ID: user2.Profile.ID, Upload: 3, Download: 4},
			},
		}); err != nil {
			return err
		}
		if err := uowctx.UpdateNodeStats(ctx, node2.ID, models.NodeStats{
			Users: []models.UserStats{
				{ID: user1.Profile.ID, Upload: 5, Download: 6},
			},
		}); err != nil {
			return err
		}
		if err := uowctx.UpdateDailyStats(ctx, time.Now().Add(-60*24*time.Hour)); err != nil {
			return err
		}
		if err := uowctx.UpdateNodeStats(ctx, node2.ID, models.NodeStats{
			Users: []models.UserStats{
				{ID: user2.Profile.ID, Upload: 11, Download: 12},
			},
		}); err != nil {
			return err
		}
		return nil
	})
	require.NoError(t, err)

	var usersList []models.UserView
	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		usersList, err = uowctx.ListUserViews(ctx)
		return err
	})
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

	var userView *models.UserView
	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		userView, err = uowctx.GetUserView(ctx, user1.Profile.ID, user1.Profile.Name)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, int64(6), userView.Traffic.Total.Upload)
	require.Equal(t, int64(8), userView.Traffic.Total.Download)
	require.Equal(t, int64(0), userView.Traffic.LastMonth.Upload)
	require.Equal(t, int64(0), userView.Traffic.LastMonth.Download)

	err = s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		userView, err = uowctx.GetUserView(ctx, user2.Profile.ID, user2.Profile.Name)
		return err
	})
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
	err := s.PasswordStorage().DoUoW(ctx, func(uowctx password.UoWContext) error {
		return uowctx.SetAuth(ctx, auth)
	})
	require.NoError(t, err)

	var dbauth *models.Auth
	err = s.PasswordStorage().DoUoW(ctx, func(uowctx password.UoWContext) error {
		dbauth, err = uowctx.GetAuth(ctx)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, auth, *dbauth)
}
