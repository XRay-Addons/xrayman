package errdefs

import (
	"errors"
	"fmt"
)

var (
	// incorrect config, access, platform
	ErrConfig              = errors.New("invalid config")
	ErrAccess              = errors.New("access error")
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrExec                = errors.New("exec command error")

	// incorrect code (nil dereference etc)
	ErrIPE           = errors.New("internal program error")
	ErrNilObjectCall = fmt.Errorf("%w: nil object call", ErrIPE)
	ErrNilArgPassed  = fmt.Errorf("%w: nil argument passed", ErrIPE)

	// errors about service commands
	ErrService         = errors.New("service error")
	ErrServiceNotReady = fmt.Errorf("%w: not ready", ErrService)

	// errors about grpc commands
	ErrGRPC         = errors.New("grpc connection error")
	ErrGRPCNotReady = fmt.Errorf("%w: not ready", ErrGRPC)

	// errors about http requests and responses
	ErrWriteResponse   = errors.New("can't write actual response")
	ErrJSONContentType = errors.New("content type is not JSON")
)
