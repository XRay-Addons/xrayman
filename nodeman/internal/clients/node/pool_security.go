package node

import (
	"fmt"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type PoolSecurity struct {
	issuer     string
	expiration time.Duration
}

func (ps *PoolSecurity) GetNodeSecurity(secret models.AccessSecret) (*NodeSecurity, error) {
	if ps == nil {
		return nil, fmt.Errorf("pool security: get node security: %w", errdefs.ErrNilObjectCall)
	}
	return &NodeSecurity{
		secret:     secret,
		issuer:     ps.issuer,
		expiration: ps.expiration,
	}, nil
}
