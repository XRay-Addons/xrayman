package httperrdefs

import (
	"net/http"
	"strings"

	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/openapi-gen"
)

var (
	ErrInvaildPayload = new(http.StatusBadRequest,
		"invalid request payload", "try to send valid")
	ErrInternalServerError = new(http.StatusInternalServerError,
		"internal server error")
	ErrAuthToken = new(http.StatusUnauthorized,
		"invalid authorization", "try another one")
	ErrUnknown = new(http.StatusInternalServerError,
		"unknown error", "we really don't know")
	ErrNotFound = new(http.StatusNotFound,
		"somebody not found", "try another")
	ErrAccessDenied = new(http.StatusForbidden,
		"Access denied", "denied deined")
	ErrTemporaryUnavailable = new(http.StatusServiceUnavailable,
		"Temporary unavailable", "please try later")
	ErrConnection = new(http.StatusExpectationFailed,
		"Connection issues", "try better connection")
)

func new(statusCode int, message string, details ...string) *api.ErrorStatusCode {
	he := api.Error{Message: message}

	if d := strings.Join(details, ""); len(d) > 0 {
		he.Details.SetTo(d)
	}

	return &api.ErrorStatusCode{
		StatusCode: statusCode,
		Response:   he,
	}
}
