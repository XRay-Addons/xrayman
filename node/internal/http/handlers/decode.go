package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

func decode(
	w http.ResponseWriter,
	r *http.Request,
	content interface{},
	logger *zap.Logger,
) int {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	fmt.Println("len(body)", len(body), body)
	if len(body) > 0 && r.Header.Get(constants.ContentType) != constants.ContentTypeJSON {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return false
	}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(content); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if err := validator.New().Struct(content); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}
