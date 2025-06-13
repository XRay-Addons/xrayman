package errdefs

import "errors"

var (
	ErrConfig = errors.New("invalid config")
	ErrXRay   = errors.New("xray apicall")
)
