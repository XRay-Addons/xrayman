// AI generated tests
package dbstorage

/*import (
	"context"
	"fmt"
	"testing"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/nodes"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNodes(t *testing.T) {
	logger := zaptest.NewLogger(t)

	s, _, db, cleanup := setupTestDB(t, logger)
	defer cleanup()
	logger.Info("new test db inited")
	ctx := context.Background()
	runningNode := models.Node{TargetStatus: models.NodeStatusRunning}
	stoppedNode := models.Node{TargetStatus: models.NodeStatusStopped}

	// add two nodes, check status
	err := s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		if err := uowctx.NewNode(ctx, &runningNode); err != nil {
			return err
		}
		if err := uowctx.NewNode(ctx, &stoppedNode); err != nil {
			return err
		}
		return nil
	})
	logger.Info(fmt.Sprintf("NewNode explain:\n%v", db.expl.List()))
	db.expl.Reset()
	require.NoError(t, err)
	require.Equal(t, int(runningNode.ID), 1)
	require.Equal(t, int(stoppedNode.ID), 2)

	// list nodes
	var nodesList []models.Node
	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		nodesList, err = uowctx.ListNodes(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	logger.Info(fmt.Sprintf("ListNodes explain:\n%v", db.expl.List()))
	db.expl.Reset()
	require.NoError(t, err)
	require.Equal(t, len(nodesList), 2)
	require.Equal(t, nodesList[0], runningNode)
	require.Equal(t, nodesList[1], stoppedNode)

	// set node status
	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		err = uowctx.SetTargetNodeStatus(ctx, 1, models.NodeStatusUnknown)
		if err != nil {
			return err
		}
		return nil
	})
	logger.Info(fmt.Sprintf("SetTargetNodeStatus explain:\n%v", db.expl.List()))
	db.expl.Reset()

	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		nodesList, err = uowctx.ListNodes(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	logger.Info(fmt.Sprintf("ListNodes explain:\n%v", db.expl.List()))
	db.expl.Reset()
	require.NoError(t, err)
	require.Equal(t, nodesList[0].TargetStatus, models.NodeStatusUnknown)

	// delete node
	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		err = uowctx.DeleteNode(ctx, 1)
		if err != nil {
			return err
		}
		return nil
	})
	logger.Info(fmt.Sprintf("DeleteNode explain:\n%v", db.expl.List()))
	db.expl.Reset()

	err = s.NodesStorage().DoUoW(ctx, func(uowctx nodes.UoWContext) error {
		nodesList, err = uowctx.ListNodes(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	logger.Info(fmt.Sprintf("ListNodes explain:\n%v", db.expl.List()))
	db.expl.Reset()

	require.NoError(t, err)
	require.Equal(t, len(nodesList), 1)
	require.Equal(t, nodesList[0].TargetStatus, models.NodeStatusStopped)

	db.expl.List()
}*/
