package security

import (
	"context"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
)

// backdoor auth token processing. for tests only!
type Backdoor struct {
}

var _ api.SecurityHandler = (*Backdoor)(nil)

func NewBackdoor() *Backdoor {
	return &Backdoor{}
}

func (b *Backdoor) HandleBearerAuth(ctx context.Context,
	operationName api.OperationName, t api.BearerAuth,
) (context.Context, error) {
	return ctx, nil
}
