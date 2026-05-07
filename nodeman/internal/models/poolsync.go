package models

import (
	"github.com/XRay-Addons/xrayman/common/xerr"
)

type NodeSyncResult struct {
	ID       NodeID
	Endpoint string
	Err      error
}

type PoolSyncStatus int

type PoolSyncResult struct {
	Nodes []NodeSyncResult
}

func (r *PoolSyncResult) GetNodeErr(id NodeID) error {
	for _, node := range r.Nodes {
		if node.ID == id {
			return node.Err
		}
	}
	return xerr.New("node not found", xerr.Withf("node id: %v", id))
}
