package bootstrap

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

func setupAdmin(ctx context.Context, password string, auth AuthService) error {
	if auth == nil {
		return errdefs.NewNilArg("auth")
	}
	return auth.AuthAdmin(ctx, password)
}