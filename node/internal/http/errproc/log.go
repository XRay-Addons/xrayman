package errproc

import (
	"context"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func LogRequestErr(ctx context.Context, log *zap.Logger, err error) {
	reqID := chimw.GetReqID(ctx)
	log.Warn("request error", zap.String(constants.RequestIDLogTag, reqID), zap.Error(err))
}

func LogResponseErr(ctx context.Context, log *zap.Logger, err error) {
	reqID := chimw.GetReqID(ctx)
	log.Error("response error", zap.String(constants.RequestIDLogTag, reqID), zap.Error(err))
}
