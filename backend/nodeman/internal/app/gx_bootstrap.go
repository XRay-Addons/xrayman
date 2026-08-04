package app

import (
	"context"
	"errors"

	"github.com/XRay-Addons/xrayman/common/gx"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/settings"
	"go.uber.org/zap"
)

type MigrateParams struct {
	gx.In
	Storage *dbstorage.Storage
	Log     *zap.Logger
}

var migrateDB = gx.Invoke(
	func(lc gx.Lifecycle, s *dbstorage.Storage) {
		lc.AppendBootstrap(gx.Bootstrap{
			Name: "migrate db",
			Fn: func(ctx context.Context) error {
				return s.Migrate(ctx)
			},
			Retry: func(err error) bool {
				return errors.Is(err, errdefs.ErrTemporaryUnavailable)
			},
		})
	},
)

type PasswordParams struct {
	gx.In
	AdminPassword string `name:"admin-password"`
	Storage       *dbstorage.Storage
	Auth          *auth.Service
}

var setPassword = gx.Invoke(
	func(lc gx.Lifecycle, p PasswordParams) {
		lc.AppendBootstrap(gx.Bootstrap{
			Name: "set password",
			Fn: func(ctx context.Context) error {
				if p.AdminPassword == "" {
					return nil
				}
				return p.Auth.Update(ctx, p.AdminPassword)
			},
			Retry: func(err error) bool {
				return errors.Is(err, errdefs.ErrTemporaryUnavailable)
			},
		})
	},
)

var ensurePassword = gx.Invoke(
	func(lc gx.Lifecycle, p PasswordParams) {
		lc.AppendBootstrap(gx.Bootstrap{
			Name: "ensure password",
			Fn: func(ctx context.Context) error {
				// check empty password
				_, err := p.Auth.Auth(ctx, models.AuthParams{Password: ""})
				if err == nil {
					return xerr.New("Empty password is not acceptable")
				}
				return err
			},
			Retry: func(err error) bool {
				return errors.Is(err, errdefs.ErrTemporaryUnavailable)
			},
		})
	},
)

type SettingsParams struct {
	gx.In
	Settings *settings.Service
}

var ensureSettings = gx.Invoke(
	func(lc gx.Lifecycle, p SettingsParams) {
		lc.AppendBootstrap(gx.Bootstrap{
			Name: "ensure settings",
			Fn: func(ctx context.Context) error {
				return p.Settings.EnsureSettings(ctx)
			},
			Retry: func(err error) bool {
				return errors.Is(err, errdefs.ErrTemporaryUnavailable)
			},
		})
	},
)

var Bootstrap = gx.Module("bootstrap",
	migrateDB,
	ensureSettings,
	setPassword,
	ensurePassword,
)
