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

func (r *PoolOpResult) JointErr() error {
	errs := make([]error, 0, len(r.Nodes))
	for _, node := range r.Nodes {
		if node.Err != nil {
			errs = append(errs, node.Err)
		}
	}
	return xerr.Join(errs...)
}
