package errproc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// write error and its reason to log and http.Response.
// nil errors allowed
func Write(ctx context.Context, err error, w http.ResponseWriter, log *zap.Logger) {
	if err == nil {
		return
	}

	// if error contains http error - write it to response,
	// if not - write as Unknown error
	respErr := ErrUnknownError
	_ = errors.As(err, &respErr)
	writeResponse(w, respErr)

	// write error to log. log level depends on status code -
	// client's errors as warnings, server - as errors
	logError(ctx, err, getLogLevel(respErr.Code()), log)
}

func writeResponse(w http.ResponseWriter, err *Response) {
	w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
	w.WriteHeader(err.Code())
	respContent := struct {
		Err     string `json:"error"`
		Details string `json:"details"`
	}{
		Err:     http.StatusText(err.Code()),
		Details: err.Error(),
	}
	_ = json.NewEncoder(w).Encode(respContent)
}

func logError(ctx context.Context, err error, lvl zapcore.Level, log *zap.Logger) {
	log.Log(lvl, "request",
		zap.String(constants.RequestIDLogTag, chimw.GetReqID(ctx)),
		zap.Error(err),
	)
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
