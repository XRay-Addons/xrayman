package errdefs

import "errors"

var (
	ErrConfig    = errors.New("invalid config")
	ErrXRay      = errors.New("xray apicall error")
	ErrIPE       = errors.New("internal program error")
	ErrService   = errors.New("xray service error")
	ErrAccess    = errors.New("access error")
	ErrCancelled = errors.New("operation cancelled")
)
