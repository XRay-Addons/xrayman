package handler

import "context"

//go:generate mockgen -source=version_service.go -destination=./mocks/mock_version_service.go -package=mocks
type VersionService interface {
	GetVersion(ctx context.Context) (string, error)
}
