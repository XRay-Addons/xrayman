package models

import (
	"github.com/XRay-Addons/xrayman/common/xerr"
)

type NodeOpResult struct {
	ID       NodeID
	Endpoint string
	Err      error
}

type PoolOpResult struct {
	Nodes []NodeOpResult
}

func (r *PoolOpResult) GetNodeErr(id NodeID) error {
	for _, node := range r.Nodes {
		if node.ID == id {
			return node.Err
		}
	}
	return xerr.Newf("node %v not found", id)
}

func (r *PoolOpResult) GetEntireErr() error {
	errs := make([]error, 0, len(r.Nodes))
	for _, node := range r.Nodes {
		if node.Err != nil {
			errs = append(errs, node.Err)
		}
	}
	return xerr.Join(errs...)
}
