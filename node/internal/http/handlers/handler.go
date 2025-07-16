package handlers

import (
	"fmt"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/errors"
	"github.com/XRay-Addons/xrayman/node/internal/http/router"
	apirequests "github.com/XRay-Addons/xrayman/node/pkg/api/requests"
	"go.uber.org/zap"
)

type Handlers struct {
	service Service
}

var _ router.Handlers = (*Handlers)(nil)

func New(service Service) (*Handlers, error) {
	if service == nil {
		return nil, fmt.Errorf("%w: handlers init: service", errdefs.ErrNilArgPassed)
	}
	return &Handlers{
		service: service,
	}, nil
}

func (h *Handlers) Start(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// convert from request
		var request apirequests.StartRequest
		if !decode(w, r, request) {
			return
		}
		users := fromStartRequest(request)

		// process
		nodeProps, err := h.service.Start(r.Context(), users)
		if err != nil {
			errors.LogRequestError(log, r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// convert to response
		response := toStartResponse(*nodeProps)
		encode(w, r, response, log)
	}
}

func (h *Handlers) Stop(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *Handlers) Status(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *Handlers) EditUsers(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
