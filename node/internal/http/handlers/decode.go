package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/go-playground/validator"
)

func decode(w http.ResponseWriter, r *http.Request, content any) bool {
	if r.Header.Get(constants.ContentType) != constants.ContentTypeJSON {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return false
	}
	if err := json.NewDecoder(r.Body).Decode(content); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if err := validator.New().Struct(content); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}
