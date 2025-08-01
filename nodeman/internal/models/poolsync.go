package models

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
