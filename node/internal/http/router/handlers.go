package router

import (
	"net/http"

	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/handlers_mock.go -package=mocks . Handlers
type Handlers interface {
	Start(log *zap.Logger) http.HandlerFunc
	Stop(log *zap.Logger) http.HandlerFunc
	Status(log *zap.Logger) http.HandlerFunc
	EditUsers(log *zap.Logger) http.HandlerFunc
}
