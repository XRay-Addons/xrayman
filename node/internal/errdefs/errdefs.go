package errdefs

import (
	"errors"
	"fmt"
)

var (
	ErrConfig          = errors.New("invalid config")
	ErrXRay            = errors.New("xray apicall error")
	ErrIPE             = errors.New("internal program error")
	ErrService         = errors.New("xray service error")
	ErrServiceNotReady = fmt.Errorf("%w: not ready", ErrService)

	ErrAccess        = errors.New("access error")
	ErrGRPC          = errors.New("grpc call error")
	ErrGRPCNotReady  = fmt.Errorf("%w: not ready", ErrGRPC)
	ErrCancelled     = errors.New("operation cancelled")
	ErrNilObjectCall = errors.New("nil object call")
	ErrNilArgPassed  = errors.New("nil argument passed")
	ErrCmd           = errors.New("command exec error")

	ErrUnsupportedPlatform = errors.New("unsupported platform")
)
