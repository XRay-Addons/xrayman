package auth

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	storage Storage
}

func New(s Storage) (*Service, error) {
	if s == nil {
		return nil, errdefs.NewNilArg("s")
	}
	return &Service{
		storage: s,
	}, nil
}

func (s *Service) AuthAdmin(ctx context.Context, password string) error {
	if s == nil {
		return errdefs.NewNilCall()
	}
	var admin *Auth
	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		admin, err = uowctx.GetAdmin(ctx)
		return
	}); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(admin.PasswordHash, []byte(password)); err != nil {
		return errdefs.WrapWithStack(err)
	}
	return nil
}

func (s *Service) SetAdmin(ctx context.Context, password string) error {
		if s == nil {
		return errdefs.NewNilCall()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return  errdefs.WrapWithStack(err)
	}

	if err := s.storage.DoUoW(ctx, func(uowctx UoWContext) (err error) {
		err = uowctx.SetAdmin(ctx, &Auth{PasswordHash: hash})
		return
	}); err != nil {
		return err
	}

	return nil
}