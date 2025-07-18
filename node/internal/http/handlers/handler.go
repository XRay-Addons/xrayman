package handlers

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/converters"
	"github.com/XRay-Addons/xrayman/node/internal/http/errwriter"
	"github.com/XRay-Addons/xrayman/node/internal/http/router"
	"github.com/XRay-Addons/xrayman/node/pkg/api"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Handlers struct {
	service   Service
	validator validator.Validate
}

var _ router.Handlers = (*Handlers)(nil)

// such beautiful world without templates

func New(service Service) (*Handlers, error) {
	if service == nil {
		return nil, fmt.Errorf("%w: handlers init: service", errdefs.ErrNilArgPassed)
	}
	return &Handlers{
		service:   service,
		validator: *validator.New(),
	}, nil
}

func (h *Handlers) Start(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		errWriter := errwriter.New(w, r, log)

		// check content type
		if err := h.ensureJSONContentType(r); err != nil {
			errWriter(err, http.StatusUnsupportedMediaType)
			return
		}

		var requestContent api.StartRequest
		if err := json.NewDecoder(r.Body).Decode(&requestContent); err != nil {
			errWriter(err, http.StatusBadRequest)
			return
		}

		if err := h.validator.Struct(requestContent); err != nil {
			errWriter(err, http.StatusBadRequest)
			return
		}

		// convert to service params
		params := converters.StartParamsFromAPI(requestContent)

		// process
		result, err := h.service.Start(r.Context(), params)
		if err != nil {
			errWriter(err, http.StatusInternalServerError)
			return
		}

		// convert to response
		response := converters.StartResultToAPI(*result)

		// write result
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			errwriter.LogRequestErr(r.Context(), zapcore.WarnLevel, err, log)
			return
		}
	}
}

func (h *Handlers) ensureJSONContentType(r *http.Request) error {
	ct := r.Header.Get(constants.ContentType)
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil || mt != constants.ContentTypeJSON {
		return errdefs.ErrJSONContentType
	}
	return nil
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
