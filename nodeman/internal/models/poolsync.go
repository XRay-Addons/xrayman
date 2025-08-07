package models

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

type NodeSyncResult struct {
	ID       NodeID
	Endpoint string
	Err      error
}

type PoolSyncStatus int

const (
	PoolSyncFailed PoolSyncStatus = iota + 1
	PoolSyncPartially
	PoolSyncOk
)

type PoolSyncResult struct {
	Status PoolSyncStatus
	Nodes  []NodeSyncResult
}

func (r *PoolSyncResult) GetNodeErr(id NodeID) error {
	for _, node := range r.Nodes {
		if node.ID == id {
			return node.Err
		}
	}
	return fmt.Errorf("pool sync result: get node err: node not exists: %w", errdefs.ErrIPE)
}
