package auth

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	storage Storage
	jwt     JWT
}

const adminTokenSubject = "admin"

func New(storage Storage, jwt JWT) (*Service, error) {
	if storage == nil {
		return nil, errdefs.NilArg("storage")
	}
	if jwt == nil {
		return nil, errdefs.NilArg("jwt")
	}
	return &Service{
		storage: storage,
		jwt:     jwt,
	}, nil
}

func (s *Service) Auth(ctx context.Context, p models.AuthParams) (*models.AuthResult, error) {
	if s == nil {
		return nil, errdefs.NilCall()
	}
	auth, err := s.storage.GetAuth(ctx)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(auth.PasswordHash,
		[]byte(p.Password),
	); err != nil {
		return nil, errdefs.AccessDenied()
	}
	token, err := s.jwt.GenerateToken(adminTokenSubject)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *Service) Update(ctx context.Context, password string) error {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	if err := s.storage.SetAuth(ctx, &models.Auth{
		PasswordHash: pwdHash,
	}); err != nil {
		return err
	}
	return nil
}
