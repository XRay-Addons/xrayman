package api

import (
	"net/http"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/http/handler/ogenserver"
)

func NewHandler(h ogenserver.Handler, s ogenserver.SecurityHandler) (http.Handler, error) {
	if h == nil {
		return nil, errdefs.NilArg("api.Handler")
	}

	apiHandler, err := ogenserver.NewServer(h, s)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return apiHandler, nil
}
