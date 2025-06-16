package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/errors"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"go.uber.org/zap"
)

type StatusHandler struct {
	service Service
	errors.ErrorsWriter
}

func NewStatusHandler(s Service, log *zap.Logger) (*StatusHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &StatusHandler{
		service:      s,
		ErrorsWriter: errors.NewErrorsWriter(log),
	}, nil
}

func (h *StatusHandler) StartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(constants.ContentType) != constants.ContentTypeJSON {
			h.WriteError(w, errors.ErrUnsupportedContentType)
			return
		}
		var request models.StartNodeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
			return
		}
		node, err := h.service.Start(r.Context(), request.Users)
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
		response := models.StartNodeResponse{
			NodeProperties: *node,
		}
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
	}
}

func (h *StatusHandler) StopHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.service.Stop(r.Context()); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *StatusHandler) StatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := h.service.Status(r.Context())
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}

		if _, err := w.Write([]byte(status)); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
	}
}
