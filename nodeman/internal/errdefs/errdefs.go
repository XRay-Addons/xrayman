package errdefs

import (
	"errors"
	"fmt"
)

var (
	// incorrect config, access, platform
	ErrConfig error = errors.New("config error")
	ErrAccess       = errors.New("access error")

	// incorrect code (nil dereference etc)
	ErrIPE           = errors.New("internal program error")
	ErrNilObjectCall = fmt.Errorf("%w: nil object call", ErrIPE)
	ErrNilArgPassed  = fmt.Errorf("%w: nil argument passed", ErrIPE)

	// errors from generated code
	ErrGenerated = errors.New("generated code error")
)
