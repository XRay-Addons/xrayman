package errors

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func LogRequestError(log *zap.Logger, req *http.Request, err error) {
	// get id from context
	reqID := chiMiddleware.GetReqID(req.Context())
	log.Error("request error", zap.String("requestID", reqID), zap.Error(err))
}
