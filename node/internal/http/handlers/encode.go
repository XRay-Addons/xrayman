package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/errors"
	"go.uber.org/zap"
)

func encode(w http.ResponseWriter, r *http.Request, content any, log *zap.Logger) bool {
	w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
	if err := json.NewEncoder(w).Encode(content); err != nil {
		errors.LogRequestError(log, r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}
