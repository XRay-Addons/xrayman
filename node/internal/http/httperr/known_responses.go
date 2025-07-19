package httperr

import "net/http"

// known error responses
var (
	ErrAuth                = NewResponse(http.StatusUnauthorized, "invalid auth JWT")
	ErrContentEncryption   = NewResponse(http.StatusUnauthorized, "invalid content JWE")
	ErrNonZeroContentLen   = NewResponse(http.StatusBadRequest, "non zero content length")
	ErrContentType         = NewResponse(http.StatusUnsupportedMediaType, "invalid content type")
	ErrContentParsing      = NewResponse(http.StatusBadRequest, "content parsing error")
	ErrContentValidation   = NewResponse(http.StatusBadRequest, "content validation error")
	ErrInternalServerError = NewResponse(http.StatusInternalServerError, "internal server error")
	ErrUnknownError        = NewResponse(http.StatusInternalServerError, "unknown error")
)
