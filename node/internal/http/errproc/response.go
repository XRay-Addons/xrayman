package errproc

import (
	"encoding/json"
	"net/http"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
)

func ResponseErr(w http.ResponseWriter, httpStatus int, details string) {
	w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
	w.WriteHeader(httpStatus)

	errDetails := struct {
		Err     string `json:"error"`
		Details string `json:"details"`
	}{
		Err:     http.StatusText(httpStatus),
		Details: details,
	}

	_ = json.NewEncoder(w).Encode(errDetails)
}
