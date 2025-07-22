package errdefs

import "errors"

var (
	ErrExec error = errors.New("command exec error")
)
