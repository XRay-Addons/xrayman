package nodesync

import (
	"context"
	"time"

	jwtools "github.com/XRay-Addons/xrayman/common/http/jwt"
	"github.com/XRay-Addons/xrayman/common/xerr"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type NodeSecurity struct {
	secret     models.AccessSecret
	issuer     string
	expiration time.Duration
}

func (s *NodeSecurity) BearerAuth(ctx context.Context,
	op api.OperationName,
) (api.BearerAuth, error) {
	token, err := jwtools.GenerateToken(s.secret[:],
		jwtools.WithIssuer(s.issuer),
		jwtools.WithTTL(s.expiration))
	if err != nil {
		return api.BearerAuth{}, xerr.WrapWithStack(err)
	}

	return api.BearerAuth{
		Token: token,
	}, nil
}
