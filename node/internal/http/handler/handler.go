package handler

import (
	"fmt"

	mw "github.com/XRay-Addons/xrayman/node/internal/http/middleware"
	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
	"go.uber.org/zap"
)

func New(s Service, log *zap.Logger) (*api.Server, error) {
	h, err := NewHandlerImpl(s)
	if err != nil {
		return nil, fmt.Errorf("handler init: %w", err)
	}
	srv, err := api.NewServer(h,
		api.WithMiddleware(mw.Transparent),
	)
	if err != nil {
		return nil, fmt.Errorf("handler init: %w", err)
	}
	return srv, nil
}
