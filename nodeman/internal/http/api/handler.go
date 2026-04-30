package api

import (
	"net/http"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	genapi "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

func NewHandler(h genapi.Handler) (http.Handler, error) {
	if h == nil {
		return nil, errdefs.NewNilArg("api.Handler")
	}

	apiHandler, err := genapi.NewServer(h)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	logged := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiHandler.ServeHTTP(w, r)
	})

	return logged, nil
}
