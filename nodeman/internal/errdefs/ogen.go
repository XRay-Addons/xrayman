package errdefs

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
)

func OgenErr(err error) error {
	details := ""

	var sc interface{ StatusCode() int }
	if errors.As(err, &sc) {
		details += fmt.Sprintf("Status: %d;", sc.StatusCode())
	} else {
		details += "Status: Transport error"
	}
	// add url path if exists
	var ue *url.Error
	if errors.As(err, &ue) {
		details += fmt.Sprintf("; URL: %s", ue.URL)
	}

	return xerr.WrapWith(ErrConnection, details)
}
