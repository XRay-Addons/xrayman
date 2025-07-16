package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/http/errors"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// create middleware for requests and responses logging
func Logger(log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := middleware.GetReqID(r.Context())

			var responseInfo responseInfo
			next.ServeHTTP(withResponseInfo(w, &responseInfo), r)

			duration := time.Since(start)

			if responseInfo.underlyingWriterErr == nil {
				errors.LogRequestError(log, r, responseInfo.underlyingWriterErr)
			}

			log.Info("request",
				zap.String("id", requestID),
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Int("status", responseInfo.status),
				zap.Duration("duration", duration),
				zap.Int("size", responseInfo.size))
		}

		return http.HandlerFunc(logFn)
	}
}

type responseInfo struct {
	status              int
	size                int
	underlyingWriterErr error
}

func withResponseInfo(w http.ResponseWriter, info *responseInfo) http.ResponseWriter {
	return &responseMiddleware{
		ResponseWriter: w,
		info:           info,
	}
}

type responseMiddleware struct {
	http.ResponseWriter
	info *responseInfo
}

var _ http.ResponseWriter = (*responseMiddleware)(nil)

func (m *responseMiddleware) Write(data []byte) (int, error) {
	// part of http.ResponseWriter's interface contract:
	// it writes StatusOK on Write if it was not called before.
	// due to lack of go's «inheritance» we can not intercept
	// such WriteHeader call from m.writer.Write (it calls directly)
	// m.writer.WriteHeader, not m.WriteHeader, so call it manually.
	if m.info.status == 0 {
		m.WriteHeader(http.StatusOK)
	}

	// chi and go/http ignores it, but we not.
	// it's not covered entire scope of errors,
	// only error of .Write method of underlying
	// response writer's method Write. in most cases
	// we can not do anything with such errors, it can
	// be related to connection troubles or go/http lib something
	size, err := m.ResponseWriter.Write(data)
	m.info.size += size
	if err != nil {
		err = fmt.Errorf("write response: %w", err)
	}
	return size, err
}

func (m *responseMiddleware) WriteHeader(status int) {
	m.ResponseWriter.WriteHeader(status)
	m.info.status = status
}
