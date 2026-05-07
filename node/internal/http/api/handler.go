package api

import (
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	genapi "github.com/XRay-Addons/xrayman/node/pkg/api/http/gen"
)

func NewHandler(h genapi.Handler, s genapi.SecurityHandler) (http.Handler, error) {
	if h == nil {
		return nil, errdefs.NewNilArg("api.Handler")
	}

	apiHandler, err := genapi.NewServer(h, s)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	// ??? WTF TODO
	// logged := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	apiHandler.ServeHTTP(w, r)
	//})

	return apiHandler, nil
}
