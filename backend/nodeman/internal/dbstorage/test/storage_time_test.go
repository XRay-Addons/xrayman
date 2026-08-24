package dbstoragetest

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////
// mass fill database
func fillUsers(ctx context.Context, db *sql.DB, n int) error {
	vnames := make([]string, 0, n)
	names := make([]string, 0, n)
	uuids := make([]string, 0, n)
	tstatuss := make([]int16, 0, n)
	revisions := make([]int64, 0, n)
	for i := range n {
		vnames = append(vnames, fmt.Sprintf("vname %d", i))
		names = append(names, fmt.Sprintf("name %d", i))
		uuids = append(uuids, fmt.Sprintf("uuid %d", i))
		if i%3 == 0 {
			tstatuss = append(tstatuss, int16(models.UserStatusDisabled))
		} else {
			tstatuss = append(tstatuss, int16(models.UserStatusEnabled))
		}
		revisions = append(revisions, int64(i%100))
	}

	req := `
	INSERT INTO users (
		display_name,
		user_name,
		vless_uuid,
		user_target_status,
		revision
	)
	SELECT *
	FROM unnest(
		$1::text[],
		$2::text[],
		$3::text[],
		$4::smallint[],
		$5::bigint[]
	) AS t(
		display_name,
		user_name,
		vless_uuid,
		user_target_status,
		revision
	);
	`
	_, err := db.ExecContext(ctx, req,
		pq.Array(vnames),
		pq.Array(names),
		pq.Array(uuids),
		pq.Array(tstatuss),
		pq.Array(revisions),
	)

	return err
}

func fillNodes(ctx context.Context, db *sql.DB, n int) error {
	cfgtemplates := make([]string, 0, n)
	endpoints := make([]string, 0, n)
	accesskey := make([][]byte, 0, n)
	cstatus := make([]int16, 0, n)
	tstatus := make([]int16, 0, n)
	revisions := make([]int64, 0, n)
	for i := range n {
		tmpl := models.ClientConfigTemplate{}
		tmplStr, err := tmpl.Value()
		if err != nil {
			return err
		}

		cfgtemplates = append(cfgtemplates, tmplStr)
		endpoints = append(endpoints, fmt.Sprintf("endpoint %d", i))
		accesskey = append(accesskey, make([]byte, 64))
		if i%2 == 0 {
			tstatus = append(tstatus, int16(models.NodeStatusStopped))
		} else {
			tstatus = append(tstatus, int16(models.NodeStatusRunning))
		}
		if i%3 == 0 {
			cstatus = append(cstatus, int16(models.NodeStatusStopped))
		} else if i%3 == 1 {
			cstatus = append(cstatus, int16(models.NodeStatusRunning))
		} else {
			cstatus = append(cstatus, int16(models.NodeStatusUnknown))
		}
		revisions = append(revisions, int64(i%10))
	}

	req := `
	INSERT INTO nodes (
		client_cfg_template,
		node_endpoint,
		node_access_key,
		node_current_status,
		node_target_status,
		revision
	)
	SELECT *
	FROM unnest(
		$1::text[],
		$2::text[],
		$3::bytea[],
		$4::smallint[],
		$5::smallint[],
		$6::bigint[]
	) AS t(
		client_cfg_template,
		node_endpoint,
		node_access_key,
		node_current_status,
		node_target_status,
		revision
	);
	`
	_, err := db.ExecContext(ctx, req,
		pq.Array(cfgtemplates),
		pq.Array(endpoints),
		pq.Array(accesskey),
		pq.Array(cstatus),
		pq.Array(tstatus),
		pq.Array(revisions),
	)

	return err
}

// //////////////////////////////////////////////////////////////////////////
// test usernodes - 100k users, 100 nodes
func TestStorage_Time_Syncs(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	s, db := setupTestDB(t, logger)
	logger.Info("new test db inited")

	////////////////////////////////////////////////////////////////////////////
	// add nodes, and users enabled, disabled
	const nNodes = 100
	err := fillNodes(ctx, db.Raw(), nNodes)
	require.NoError(t, err)

	const nUsers = 100000
	err = fillUsers(ctx, db.Raw(), nUsers)
	require.NoError(t, err)

	////////////////////////////////////////////////////////////////////////////
	// let's go!
	expl := db.WithExplanations("FindPendingSyncs", ExplainAnalyze)
	syncs, err := s.FindPendingSyncs(ctx, models.NodeID(nNodes/2))
	require.NoError(t, err)
	metrics, err := expl.Metrics()
	require.NoError(t, err)
	metrics.Print(logger)
	// find pending syncs x nodes count <= 10 seconds
	require.Less(t, nNodes*metrics.ExecutionTime, 10*time.Second)
	logger.Info("find pending syncs", zap.Int("syncs", len(syncs)))

	err = s.SetNodeRev(ctx, models.NodeID(2), models.Revision(nUsers))
	require.NoError(t, err)
	expl = db.WithExplanations("GetUserNodes", ExplainAnalyze)
	usernodes, err := s.GetUserNodes(ctx, models.UserID(3))
	metrics, err = expl.Metrics()
	require.NoError(t, err)
	metrics.Print(logger)
	// find user nodes <= 10ms
	require.Less(t, metrics.ExecutionTime, 10*time.Millisecond)
	logger.Info("get user nodes", zap.Int("nodes", len(usernodes)))
}

// //////////////////////////////////////////////////////////////////////////
// test stats - 10k users, 10 nodes
func TestStorage_Time_Stats(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	s, db := setupTestDB(t, logger)
	logger.Info("new test db inited")

	////////////////////////////////////////////////////////////////////////////
	// add nodes, enabled, disabled
	const nNodes = 10
	err := fillNodes(ctx, db.Raw(), nNodes)
	require.NoError(t, err)
	logger.Info("nodes added")

	////////////////////////////////////////////////////////////////////////////
	// add users, enabled, disabled
	const nUsers = 10000
	err = fillUsers(ctx, db.Raw(), nUsers)
	require.NoError(t, err)
	logger.Info("users added")

	////////////////////////////////////////////////////////////////////////////
	// add stats, periodically update daily stats
	for i := range nNodes {
		err = s.DoTx(ctx, func(ctx context.Context) error {
			stats := make([]models.UserStats, nUsers, nUsers)
			for u := range nUsers {
				stats[u].ID = u
				stats[u].Downlink = int64(u - i)
				stats[u].Uplink = int64(2*u - i)
			}
			if err := s.UpdateNodeStats(ctx, i, models.NodeStats{
				Users: stats,
			}); err != nil {
				return err
			}
			if i%10 == 0 {
				var randomDayOffset time.Duration = time.Duration(nNodes-i) * 7 * 24 * time.Hour
				if err := s.UpdateDailyStats(ctx,
					time.Now().Add(-randomDayOffset),
				); err != nil {
					return err
				}
			}
			return nil
		})
		require.NoError(t, err)
	}
	logger.Info("stats added")

	////////////////////////////////////////////////////////////////////////////
	// update stats
	expl := db.WithExplanations("UpdateNodeStats", ExplainAnalyze)
	stats := make([]models.UserStats, nUsers, nUsers)
	for u := range nUsers {
		stats[u].ID = u
		stats[u].Downlink = int64(u)
		stats[u].Uplink = int64(2 * u)
	}
	err = s.UpdateNodeStats(ctx, models.NodeID(nNodes/2), models.NodeStats{
		Users: stats,
	})
	require.NoError(t, err)
	metrics, err := expl.Metrics()
	require.NoError(t, err)
	metrics.Print(logger)
	require.Less(t, nNodes*metrics.ExecutionTime, 1*time.Second)

	////////////////////////////////////////////////////////////////////////////
	// update daily sync
	expl = db.WithExplanations("UpdateDailyStats", ExplainAnalyze)
	err = s.UpdateDailyStats(ctx, time.Now())
	require.NoError(t, err)
	metrics, err = expl.Metrics()
	require.NoError(t, err)
	metrics.Print(logger)
	require.Less(t, metrics.ExecutionTime, 1*time.Second)

	////////////////////////////////////////////////////////////////////////////
	// find user
	expl = db.WithExplanations("GetUserView", ExplainAnalyze)
	_, err = s.GetUserView(ctx, 2, "name 1")
	require.NoError(t, err)
	metrics, err = expl.Metrics()
	require.NoError(t, err)
	metrics.Print(logger)
	require.Less(t, nUsers*metrics.ExecutionTime, 1*time.Second)

	////////////////////////////////////////////////////////////////////////////
	// list users
	expl = db.WithExplanations("ListUserViews", ExplainAnalyze)
	_, err = s.ListUserViews(ctx)
	require.NoError(t, err)
	metrics, err = expl.Metrics()
	require.NoError(t, err)
	metrics.Print(logger)
	require.Less(t, metrics.ExecutionTime, 1*time.Second)
}
