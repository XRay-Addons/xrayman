package api

import (
	"net/http"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
	genapi "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
)

func NewHandler(h genapi.Handler, s genapi.SecurityHandler) (http.Handler, error) {
	if h == nil {
		return nil, errdefs.NilArg("api.Handler")
	}

	apiHandler, err := genapi.NewServer(h, s)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	// ??? WTF TODO
	// logged := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	apiHandler.ServeHTTP(w, r)
	//})

	return apiHandler, nil
}
