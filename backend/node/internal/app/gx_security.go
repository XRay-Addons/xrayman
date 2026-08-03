package app

import (
	"crypto/tls"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/node/internal/http/security"
	"github.com/XRay-Addons/xrayman/node/internal/infra/auth/jwt"
	"github.com/XRay-Addons/xrayman/node/internal/infra/secrets"
	"github.com/XRay-Addons/xrayman/node/internal/infra/tlscfg"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"go.uber.org/fx"
)

var Security = gx.Module("security",
	gx.ProvideAnnotated(
		jwt.New,
		gx.As(new(security.JWT)),
	),
	gx.ProvideAnnotated(
		func(d string) (*secrets.Secrets, error) {
			return secrets.Init(d)
		},
		gx.ParamTags(`name:"persistent-dir"`),
	),
	fx.Provide(
		func(s *secrets.Secrets) models.AccessKey {
			return s.AccessKey
		},
	),
	fx.Provide(
		func(k models.AccessKey) models.AccessSecret {
			return k.AccessSecret
		},
	),
	fx.Provide(
		func(s *secrets.Secrets) (*tls.Config, error) {
			return tlscfg.Load(s.Cert, s.Key)
		},
	),
)
