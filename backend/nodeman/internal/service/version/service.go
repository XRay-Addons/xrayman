package version

import (
	"context"

	"github.com/XRay-Addons/xrayman/nodeman/internal/http/handler"
	v "github.com/XRay-Addons/xrayman/nodeman/internal/version"
)

type Service struct {
}

var _ handler.VersionService = (*Service)(nil)

func New() *Service {
	return &Service{}
}

func (s *Service) GetVersion(ctx context.Context) (string, error) {
	return v.Version, nil
}
