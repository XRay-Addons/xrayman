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

type UsersHandler struct {
	service Service
	errors.ErrorsWriter
}

func NewUsersHandler(s Service, log *zap.Logger) (*UsersHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &UsersHandler{
		service:      s,
		ErrorsWriter: errors.NewErrorsWriter(log),
	}, nil
}

func (h *UsersHandler) AddUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(constants.ContentType) != constants.ContentTypeJSON {
			h.WriteError(w, errors.ErrUnsupportedContentType)
			return
		}
		var request models.AddUsersRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
			return
		}
		if err := h.service.AddUsers(r.Context(), request.Users); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UsersHandler) DelUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(constants.ContentType) != constants.ContentTypeJSON {
			h.WriteError(w, errors.ErrUnsupportedContentType)
			return
		}
		var request models.DelUsersRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
			return
		}
		if err := h.service.DelUsers(r.Context(), request.Users); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
