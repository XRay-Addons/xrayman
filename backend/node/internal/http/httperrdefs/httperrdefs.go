package httperrdefs

import (
	"net/http"
	"strings"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
)

var (
	ErrInternalServerError = new(http.StatusInternalServerError,
		"Internal server error")
	ErrAuthToken = new(http.StatusUnauthorized,
		"Invalid auth token", "try another one")
	ErrAccessDenied = new(http.StatusForbidden,
		"Access denied", "denied deined")
	ErrTemporaryUnavailable = new(http.StatusServiceUnavailable,
		"Temporary unavailable", "please try later")
	ErrConnection = new(http.StatusExpectationFailed,
		"Connection issues", "try better connection")
	ErrUnknown = new(http.StatusInternalServerError,
		"unknown error", "we really don't know")
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
