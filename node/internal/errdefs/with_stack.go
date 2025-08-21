package errdefs

import (
	"fmt"
	"strings"

	"github.com/go-faster/errors"
)

func WithStack(err error) error {
	return &errWithStack{
		err:   err,
		stack: getTrace(3),
	}
}

type errWithStack struct {
	err   error
	stack []string
}

var _ error = (*errWithStack)(nil)
var _ errors.Wrapper = (*errWithStack)(nil)

func (w *errWithStack) Error() string {
	if w.stack == nil {
		return w.err.Error()
	}
	return fmt.Sprintf("-> %s:\n\t%s",
		strings.Join(w.stack, "\n-> "), w.err.Error())
}

func (e *errWithStack) Unwrap() error {
	return e.err
}
