package middleware

import (
	"net/http"
	"strings"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/errors"
	"go.uber.org/zap"
)

func Compression(log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		compression := func(w http.ResponseWriter, r *http.Request) {
			// handle zipped request - replace request body to unzipped
			if isCompressedRequest(r) {
				cr, err := newGZipReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					errors.LogRequestError(log, r, err)
					return
				}
				r.Body = cr
				// don't move this code block out of this function!
				defer func() {
					if err := cr.Close(); err != nil {
						errors.LogRequestError(log, r, err)
					}
				}()
			}

			// client supports compression - replace response writer
			if compressedResponseSupported(r) {
				cw, err := newGZipWriter(w)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					errors.LogRequestError(log, r, err)
					return
				}
				w = cw
				// don't move this code block out of this function!
				defer func() {
					if err := cw.Close(); err != nil {
						errors.LogRequestError(log, r, err)
					}
				}()
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(compression)
	}
}

func isCompressedRequest(r *http.Request) bool {
	return lookupHeaderComponent(
		r.Header.Values(constants.ContentEncoding),
		constants.GZipEncoding)
}

func compressedResponseSupported(r *http.Request) bool {
	return lookupHeaderComponent(
		r.Header.Values(constants.AcceptEncoding),
		constants.GZipEncoding)
}

func lookupHeaderComponent(header []string, target string) bool {
	for _, values := range header {
		for _, value := range strings.Split(values, ",") {
			if strings.TrimSpace(value) == target {
				return true
			}
		}
	}
	return false
}
