package apihandler

import (
	"net/http"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	api "github.com/XRay-Addons/xrayman/nodeman/pkg/api/http/gen"
	"log"
)

func New(h api.Handler) (http.Handler, error) {
	if h == nil {
		return nil, errdefs.NewNilArg("api.Handler")
	}

	apiHandler, err := api.NewServer(h)
	if err != nil {
		return nil, errdefs.WrapWithStack(err)
	}

	logged := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[API] %s %s", r.Method, r.URL.Path)
		apiHandler.ServeHTTP(w, r)
	})

	return logged, nil
}