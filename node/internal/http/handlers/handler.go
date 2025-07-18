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

		// check content type
		if err := h.ensureJSONContentType(r); err != nil {
			errproc.ResponseErr(w, http.StatusUnsupportedMediaType, err.Error())
			errproc.LogRequestErr(r.Context(), log, err)
			return
		}

		var requestContent api.StartRequest
		if err := json.NewDecoder(r.Body).Decode(&requestContent); err != nil {
			errproc.ResponseErr(w, http.StatusBadRequest, err.Error())
			errproc.LogRequestErr(r.Context(), log, err)
			return
		}

		if err := h.validator.Struct(requestContent); err != nil {
			errproc.ResponseErr(w, http.StatusBadRequest, err.Error())
			errproc.LogRequestErr(r.Context(), log, err)
			return
		}

		// convert to service params
		params := converters.StartParamsFromAPI(requestContent)

		// process
		result, err := h.service.Start(r.Context(), params)
		if err != nil {
			errproc.ResponseErr(w, http.StatusInternalServerError, "")
			errproc.LogResponseErr(r.Context(), log, err)
			return
		}

		// convert to response
		response := converters.StartResultToAPI(*result)

		// write result, don't log results
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		json.NewEncoder(w).Encode(response)
	}
}

func (h *Handler) parseJSONRequest[T any](
	w http.ResponseWriter,
	r *http.Request,
	log *zap.Logger,
) (*T, error) {
    // 1. Проверка Content-Type
    if err := h.ensureJSONContentType(r); err != nil {
        return nil, fmt.Errorf("content type validation failed: %w", err)
    }

    // 2. Декодирование JSON
    var requestContent T
    if err := json.NewDecoder(r.Body).Decode(&requestContent); err != nil {
        return nil, fmt.Errorf("json decoding failed: %w", err)
    }

    // 3. Валидация структуры (если валидатор доступен)
    if h.validator != nil {
        if err := h.validator.Struct(requestContent); err != nil {
            return nil, fmt.Errorf("request validation failed: %w", err)
        }
    }

    return &requestContent, nil
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
