package errdefs

import (
	"errors"
	"fmt"
)

func New(text string) error {
	return WithStack(errors.New(text))
}

func Newf(text string, args ...any) error {
	return WithStack(fmt.Errorf(text, args...))
}
