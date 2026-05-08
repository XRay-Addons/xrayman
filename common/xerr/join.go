package xerr

import (
	"errors"
	"fmt"
)

type joined struct {
	err error
}

var _ error = (*joined)(nil)

func (j *joined) Error() string {
	return j.err.Error()
}

func Join(errs ...error) error {
	if je := errors.Join(errs...); je != nil {
		return &joined{err: je}
	}
	return nil
}

func (j *joined) Format(f fmt.State, verb rune) {
	if j == nil {
		return
	}

	ffallback := "%" + string(verb)

	u, ok := j.err.(interface{ Unwrap() []error })
	if !ok {
		// nested error is not unwrappable, so ok
		formatErrQuant(j.err, f, verb, ffallback)
		return
	}

	for _, e := range u.Unwrap() {
		// format nested components
		formatErrQuant(e, f, verb, ffallback)
	}
}

// format unsplittable basic error.
// user Format if supported, else just printf
func formatErrQuant(e error, f fmt.State, verb rune, ffallback string) {
	if fe, ok := e.(interface{ Format(fmt.State, rune) }); ok {
		fe.Format(f, verb)
	} else {
		fmt.Fprintf(f, ffallback, e)
	}
}
