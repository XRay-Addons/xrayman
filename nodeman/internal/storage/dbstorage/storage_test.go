package dbstorage

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/sync/poolsync"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/subscr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/users"
	"github.com/XRay-Addons/xrayman/nodeman/internal/storage/dbstorage/sqldb"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T, l *zap.Logger) (
	s *Storage,
	postgres *embeddedpostgres.EmbeddedPostgres,
	cleanup func(),
) {
	t.Helper()

	dataDir := fmt.Sprintf("%s/pg_temp_%d", os.TempDir(), time.Now().UnixNano())

	postgres = embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Version(embeddedpostgres.V15).
			DataPath(dataDir).
			Username("test").
			Password("test").
			Database("testdb").
			Port(5434),
		//Logger(zap.NewStdLog(l).Writer()) - not working
	)

	err := postgres.Start()
	require.NoError(t, err, "failed to start embedded postgres")

	connStr := "host=localhost port=5434 user=test password=test dbname=testdb sslmode=disable"
	db, err := sqldb.New(connStr)
	require.NoError(t, err, "failed to open db")
	s, err = New(context.TODO(), db, WithMigration(), WithLogger(l))
	require.NoError(t, err, "failed to create storage")

	cleanup = func() {
		_ = sqldb.Close(db)
		_ = postgres.Stop()
		_ = os.RemoveAll(dataDir)
	}

	return
}

func TestDBStorage(t *testing.T) {
	logger := zaptest.NewLogger(t)

	s, db, cleanup := newTestDB(t, logger)
	defer cleanup()
	logger.Info("new test db inited")

	ctx := context.Background()

	// add two users - enabled and disabled
	enabledUser := models.User{TargetStatus: models.UserStatusEnabled}
	disabledUser := models.User{TargetStatus: models.UserStatusDisabled}
	err := s.UsersStorage().DoUoW(ctx, func(uowctx users.UoWContext) error {
		if err := uowctx.NewUser(ctx, &enabledUser); err != nil {
			return err
		}
		if err := uowctx.NewUser(ctx, &disabledUser); err != nil {
			return err
		}
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, int(enabledUser.Profile.ID), 1)
	require.Equal(t, int(disabledUser.Profile.ID), 2)

	// add two nodes - on and off
	runningNode := models.Node{TargetStatus: models.NodeStatusRunning}
	stoppedNode := models.Node{TargetStatus: models.NodeStatusStopped}
	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		if err := uowctx.NewNode(ctx, &runningNode); err != nil {
			return err
		}
		if err := uowctx.NewNode(ctx, &stoppedNode); err != nil {
			return err
		}
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, int(runningNode.ID), 1)
	require.Equal(t, int(stoppedNode.ID), 2)
	time.Sleep(2 * time.Second)
	logger.Info("First test passed")

	// restart db
	err = db.Stop()
	require.NoError(t, err)
	err = db.Start()
	require.NoError(t, err)

	logger.Info("DB restarted")

	// request user nodes
	var userNodes []models.Node
	err = s.SubscrStorage().DoUoW(ctx, func(uowctx subscr.UoWContext) (err error) {
		userNodes, err = uowctx.GetUserNodes(ctx, enabledUser.Profile.ID)
		return
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(userNodes))

	// request pending syncs:
	// Running Node: expected one - running node 1 and disabled user 1)
	var pendingSyncs []models.UserSyncStatus
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) (err error) {
		pendingSyncs, err = uowctx.FindPendingSyncs(ctx, runningNode.ID)
		return
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(pendingSyncs))
	require.Equal(t, models.UserSyncStatus{
		User:          enabledUser,
		CurrentStatus: models.UserStatusDisabled,
	}, pendingSyncs[0])

	// Stopped Node: expected one - running node 1 and disabled user 1)
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) (err error) {
		pendingSyncs, err = uowctx.FindPendingSyncs(ctx, stoppedNode.ID)
		return
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(pendingSyncs))
	require.Equal(t, models.UserSyncStatus{
		User:          enabledUser,
		CurrentStatus: models.UserStatusDisabled,
	}, pendingSyncs[0])

	// apply pending syncs
	syncsPatch := []models.UserStatusPatch{
		{
			UserID: enabledUser.Profile.ID,
			Status: models.UserStatusEnabled,
		},
	}
	err = s.PoolSyncStorage().DoUoW(ctx, func(uowctx poolsync.UoWContext) error {
		return uowctx.UpdateNodeUsers(ctx, runningNode.ID, syncsPatch)
	})
	require.NoError(t, err)

	// request user nodes again
	err = s.SubscrStorage().DoUoW(ctx, func(uowctx subscr.UoWContext) (err error) {
		userNodes, err = uowctx.GetUserNodes(ctx, enabledUser.Profile.ID)
		return
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(userNodes))

	logger.Info("Next test passed,closing...")
}
