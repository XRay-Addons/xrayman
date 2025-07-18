package handlers

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/converters"
	"github.com/XRay-Addons/xrayman/node/internal/http/errproc"
	"github.com/XRay-Addons/xrayman/node/internal/http/router"
	"github.com/XRay-Addons/xrayman/node/pkg/api"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
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

		var err error
		defer func() { errproc.Write(r.Context(), err, w, log) }()

		// parse request content
		requestContent, err := parseJSONRequest[api.StartRequest](r, h.validator)
		if err != nil {
			return
		}

		// convert to service params
		params := converters.StartParamsFromAPI(requestContent)

		// process
		result, err := h.service.Start(r.Context(), params)
		if err != nil {
			return
		}

		// convert to response
		response := converters.StartResultToAPI(result)

		// write result, don't log results
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		json.NewEncoder(w).Encode(response)
	}
}

func (h *Handlers) Stop(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error
		defer func() { errproc.Write(r.Context(), err, w, log) }()

		// request has no content
		requestContent := &api.StopRequest{}

		// convert to service params
		params := converters.StopParamsFromAPI(requestContent)

		// process
		result, err := h.service.Stop(r.Context(), params)
		if err != nil {
			return
		}

		// convert to response
		response := converters.StopResultToAPI(result)

		// write result, don't log results
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		json.NewEncoder(w).Encode(response)
	}
}

func (h *Handlers) Status(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error
		defer func() { errproc.Write(r.Context(), err, w, log) }()

		// request has no content
		requestContent := &api.StatusRequest{}

		// convert to service params
		params := converters.StatusParamsFromAPI(requestContent)

		// process
		result, err := h.service.Status(r.Context(), params)
		if err != nil {
			return
		}

		// convert to response
		response := converters.StatusResultToAPI(result)

		// write result, don't log results
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		json.NewEncoder(w).Encode(response)
	}
}

func (h *Handlers) EditUsers(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error
		defer func() { errproc.Write(r.Context(), err, w, log) }()

		// parse request content
		requestContent, err := parseJSONRequest[api.EditUsersRequest](r, h.validator)
		if err != nil {
			return
		}

		// convert to service params
		params := converters.EditUsersParamsFromAPI(requestContent)

		// process
		result, err := h.service.EditUsers(r.Context(), params)
		if err != nil {
			return
		}

		// convert to response
		response := converters.EditUsersResultToAPI(result)

		// write result, don't log results
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		json.NewEncoder(w).Encode(response)
	}
}

func parseJSONRequest[T any](r *http.Request, v validator.Validate) (*T, error) {
	// check content type
	mt, _, err := mime.ParseMediaType(r.Header.Get(constants.ContentType))
	if err != nil {
		err = fmt.Errorf("parse json request: media type: %w", err)
		return nil, errproc.NewError(errproc.ErrContentType, err)
	}
	if mt != constants.ContentTypeJSON {
		err = fmt.Errorf("parse json request: media type: %s", mt)
		return nil, errproc.NewError(errproc.ErrContentType, err)
	}

	// parse request
	var requestContent T
	if err := json.NewDecoder(r.Body).Decode(&requestContent); err != nil {
		err = fmt.Errorf("parse json request: %w", err)
		return nil, errproc.NewError(errproc.ErrContentParsing, err)
	}

	// validate
	if err := v.Struct(requestContent); err != nil {
		err = fmt.Errorf("validate json request: %w", err)
		return nil, errproc.NewError(errproc.ErrContentValidation, err)
	}

	// finally, return something
	return &requestContent, nil
}
