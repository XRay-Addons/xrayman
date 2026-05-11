package middleware

import (
	"context"
	"net/http"
)

const (
	HeadersTag = "middleware headers"
)

// create middleware for writing headers to response
func Headers() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		headersFn := func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			ctx := context.WithValue(r.Context(), HeadersTag, headers)
			rr := r.WithContext(ctx)
			next.ServeHTTP(w, rr)
		}

		return http.HandlerFunc(headersFn)
	}
}

// access to headers from request context. may return nil object
// if headers dont exist
func GetHeaders(ctx context.Context) http.Header {
	if h, ok := ctx.Value(HeadersTag).(http.Header); ok && h != nil {
		return h
	}
	return nil
}
