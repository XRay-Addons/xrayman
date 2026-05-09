package httperr

import (
	"fmt"
	"net/http"
	"strings"

	api "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
)

var (
	ErrInternalServerError = new(http.StatusInternalServerError,
		"Internal server error")
	ErrAuthToken = new(http.StatusUnauthorized,
		"Invalid auth token", "try another one")
	ErrUnknown = new(http.StatusInternalServerError,
		"unknown error", "we really don't know")
)

// error impl containing api.ErrorStatusCode
// to return it as error from middleware or handlers
// and later process in handler.NewError
type HttpErr api.ErrorStatusCode

func new(statusCode int, message string, details ...string) *HttpErr {
	he := api.Error{Message: message}

	if d := strings.Join(details, ""); len(d) > 0 {
		he.Details.SetTo(d)
	}

	return &HttpErr{
		StatusCode: statusCode,
		Response:   he,
	}
}

func (e HttpErr) Error() string {
	return fmt.Sprintf("http error: %s: %s",
		http.StatusText(e.StatusCode), e.Response.Message)
}
