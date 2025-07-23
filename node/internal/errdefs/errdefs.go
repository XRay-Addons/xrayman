package errdefs

import (
	"errors"
	"fmt"
)

var (
	// incorrect config, access, platform
	ErrConfig              error = errors.New("config error")
	ErrAccess                    = errors.New("access error")
	ErrExec                error = errors.New("command exec error")
	ErrUnsupportedPlatform error = errors.New("unsupported platform")

	// incorrect code (nil dereference etc)
	ErrIPE           = errors.New("internal program error")
	ErrNilObjectCall = fmt.Errorf("%w: nil object call", ErrIPE)

	// errors about service commands
	ErrService         = errors.New("service error")
	ErrServiceNotReady = fmt.Errorf("%w: not ready", ErrService)
)
