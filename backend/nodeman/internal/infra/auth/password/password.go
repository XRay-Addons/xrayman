package password

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/app/bootstrap"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	storage Storage
}

var _ bootstrap.Password = (*Password)(nil)

var _ auth.Password = (*Password)(nil)

func New(s Storage) (*Password, error) {
	if s == nil {
		return nil, errdefs.NilArg("s")
	}
	return &Password{
		storage: s,
	}, nil
}

func (p *Password) Verify(ctx context.Context, password string) error {
	var auth *models.Auth
	if err := p.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		auth, err = uowctx.GetAuth(ctx)
		return
	}); err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(auth.PasswordHash, []byte(password)); err != nil {
		return errdefs.AccessDenied()
	}
	return nil
}

func (p *Password) Update(ctx context.Context, password string) error {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	if err := p.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		uowctx.SetAuth(ctx, models.Auth{PasswordHash: pwdHash})
		return
	}); err != nil {
		return err
	}

	return nil
}
