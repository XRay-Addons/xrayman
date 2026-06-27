package nodes

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type TxFn = func(context.Context) error

type Storage interface {
	// add new node to storage, assign NodeID to node
	NewNode(ctx context.Context, node *models.Node) error
	// get all nodes
	ListNodes(ctx context.Context) ([]models.Node, error)
	// change node target status
	SetTargetNodeStatus(ctx context.Context, id models.NodeID,
		status models.NodeStatus) error
	// delete node
	DeleteNode(ctx context.Context,
		id models.NodeID) error
	// call multiple operations as tx
	DoTx(ctx context.Context, fn TxFn) error
}
