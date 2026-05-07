package auth

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

type Service struct {
	pwd Password
	jwt JWT
}

const adminTokenSubject = "admin"

func New(pwd Password, jwt JWT) (*Service, error) {
	if pwd == nil {
		return nil, errdefs.NilArg("pwd")
	}
	if jwt == nil {
		return nil, errdefs.NilArg("jwt")
	}
	return &Service{
		pwd: pwd,
		jwt: jwt,
	}, nil
}

func (s *Service) Auth(ctx context.Context, p models.AuthParams) (*models.AuthResult, error) {
	if s == nil {
		return nil, errdefs.NilCall()
	}
	err := s.pwd.Verify(ctx, p.Password)
	if err != nil {
		return nil, err
	}
	token, err := s.jwt.GenerateToken(adminTokenSubject)
	if err != nil {
		return nil, err
	}
	return token, nil
}
