package app

import (
	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/security"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/jwt"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
)

var j = gx.Provide(
	gx.Annotate(
		func(s string, iss string) (*jwt.JWT, error) {
			return jwt.New(s, jwt.WithIssuer(iss))
		},
		gx.ParamTags(`name:"jwt-secret"`, `name:"jwt-issuer"`),
		gx.As(new(auth.JWT)),
		gx.As(new(security.JWT)),
	),
)

var Security = gx.Module("security",
	j,
)
