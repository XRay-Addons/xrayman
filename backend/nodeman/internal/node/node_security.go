package node

import (
	"context"
	"time"

	jwtools "github.com/XRay-Addons/xrayman/common/http/jwt"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/node/ogenclient"
)

type NodeSecurity struct {
	secret     models.AccessSecret
	issuer     string
	expiration time.Duration
}

func (s *NodeSecurity) BearerAuth(ctx context.Context,
	op ogenclient.OperationName,
) (ogenclient.BearerAuth, error) {
	token, err := jwtools.GenerateToken(s.secret[:],
		jwtools.WithIssuer(s.issuer),
		jwtools.WithTTL(s.expiration))
	if err != nil {
		return ogenclient.BearerAuth{}, xerr.WrapWithStack(err)
	}

	return ogenclient.BearerAuth{
		Token: token,
	}, nil
}
