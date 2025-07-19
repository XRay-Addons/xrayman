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

// such beautiful world without OOP

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

		var requestContent api.StartRequest
		if err = parseJSONRequest(r, h.validator, &requestContent); err != nil {
			return
		}
		params := converters.StartParamsFromAPI(requestContent)

		result, err := h.service.Start(r.Context(), params)
		if err != nil {
			return
		}

		response := converters.StartResultToAPI(*result)
		writeJSONResponse(w, response)
	}
}

func (h *Handlers) Stop(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error
		defer func() { errproc.Write(r.Context(), err, w, log) }()

		requestContent := api.StopRequest{}
		if err := checkEmptyRequest(r); err != nil {
			return
		}
		params := converters.StopParamsFromAPI(requestContent)

		result, err := h.service.Stop(r.Context(), params)
		if err != nil {
			return
		}

		response := converters.StopResultToAPI(*result)
		writeJSONResponse(w, response)
	}
}

func (h *Handlers) Status(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error
		defer func() { errproc.Write(r.Context(), err, w, log) }()

		requestContent := api.StatusRequest{}
		if err := checkEmptyRequest(r); err != nil {
			return
		}
		params := converters.StatusParamsFromAPI(requestContent)

		result, err := h.service.Status(r.Context(), params)
		if err != nil {
			return
		}

		response := converters.StatusResultToAPI(*result)
		writeJSONResponse(w, response)
	}
}

func (h *Handlers) EditUsers(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error
		defer func() { errproc.Write(r.Context(), err, w, log) }()

		var requestContent api.EditUsersRequest
		if err = parseJSONRequest(r, h.validator, &requestContent); err != nil {
			return
		}
		params := converters.EditUsersParamsFromAPI(requestContent)

		result, err := h.service.EditUsers(r.Context(), params)
		if err != nil {
			return
		}

		response := converters.EditUsersResultToAPI(*result)
		writeJSONResponse(w, response)
	}
}

func checkEmptyRequest(r *http.Request) error {
	if r.ContentLength > 0 {
		err := fmt.Errorf("non-empty request content length")
		return errproc.NewError(errproc.ErrNonZeroContentLen, err)
	}
	buf := make([]byte, 1)
	if n, _ := r.Body.Read(buf); n > 0 {
		err := fmt.Errorf("non-empty request content length")
		return errproc.NewError(errproc.ErrNonZeroContentLen, err)
	}
	return nil
}

func parseJSONRequest[T any](r *http.Request, v validator.Validate, content *T) error {
	// check content type
	mt, _, err := mime.ParseMediaType(r.Header.Get(constants.ContentType))
	if err != nil {
		err = fmt.Errorf("parse json request: media type: %w", err)
		return errproc.NewError(errproc.ErrContentType, err)
	}
	if mt != constants.ContentTypeJSON {
		err = fmt.Errorf("parse json request: media type: %s", mt)
		return errproc.NewError(errproc.ErrContentType, err)
	}

	// parse request
	if err := json.NewDecoder(r.Body).Decode(content); err != nil {
		err = fmt.Errorf("parse json request: %w", err)
		return errproc.NewError(errproc.ErrContentParsing, err)
	}

	// validate
	if err := v.Struct(*content); err != nil {
		err = fmt.Errorf("validate json request: %w", err)
		return errproc.NewError(errproc.ErrContentValidation, err)
	}

	return nil
}

func writeJSONResponse[T any](w http.ResponseWriter, r T) {
	w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
	json.NewEncoder(w).Encode(r)
}
