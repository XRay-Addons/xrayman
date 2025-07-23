package middleware

import (
	omw "github.com/ogen-go/ogen/middleware"
)

// sample middleware, TODO: remove
func Transparent(req omw.Request, next omw.Next) (omw.Response, error) {
	return next(req)
}
