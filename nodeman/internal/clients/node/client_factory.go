package client

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service"
)

type SecurityFactory interface {
	GetSecurity(secret []byte) (api.SecuritySource, error)
}

type ClientFactory struct {
	httpClient *http.Client
	security   SecurityFactory
}

var _ service.NodeClientFactory = (*ClientFactory)(nil)

func NewClientFactory(httpClient *http.Client, security SecurityFactory) (
	*ClientFactory, error,
) {
	if security == nil {
		return nil, fmt.Errorf("node client factory: init: %w", errdefs.ErrNilArgPassed)
	}

	return &ClientFactory{
		httpClient: httpClient,
		security:   security,
	}, nil
}

func (f *ClientFactory) Get(ctx context.Context,
	endpoint string, secret []byte,
) (service.NodeClient, error) {
	if f == nil || f.security == nil {
		return nil, fmt.Errorf("node client factory: get: %w", errdefs.ErrNilObjectCall)
	}
	var security api.SecuritySource
	if f.security != nil {
		var err error
		if security, err = f.security.GetSecurity(secret); err != nil {
			return nil, fmt.Errorf("node client factory: get: %w", err)
		}
	}
	return NewNodeClient(endpoint, security, f.httpClient)
}
