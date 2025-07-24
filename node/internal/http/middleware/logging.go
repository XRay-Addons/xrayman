package middleware

import (
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// create middleware for requests and responses logging
func Logger(log *zap.Logger) Middleware {
	if log == nil {
		log = zap.NewNop()
	}

	return func(next http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				log.Log(getLogLevel(ww.Status()), "request",
					zap.String(constants.RequestIDLogTag, chimw.GetReqID(r.Context())),
					zap.String("uri", r.URL.Path),
					zap.String("method", r.Method),
					zap.Int("status", ww.Status()),
					zap.Duration("duration", time.Since(start)),
					zap.Int("size", ww.BytesWritten()))
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(logFn)
	}
}

func getLogLevel(statusCode int) zapcore.Level {
	switch {
	case statusCode < http.StatusBadRequest:
		return zapcore.InfoLevel
	case statusCode < http.StatusInternalServerError:
		return zapcore.WarnLevel
	default:
		return zapcore.ErrorLevel
	}
}
