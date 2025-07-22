package errdefs

import "errors"

var (
	ErrExec   error = errors.New("command exec error")
	ErrConfig error = errors.New("config error")
)
