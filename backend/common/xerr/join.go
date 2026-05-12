package xerr

import (
	"fmt"
	"strings"
)

type multierr interface {
	Error() string
	Unwrap() []error
}

type joined struct {
	errs []error
}

var _ error = (*joined)(nil)

func Join(errs ...error) error {
	j := joined{
		errs: make([]error, 0, len(errs)),
	}
	for _, e := range errs {
		if e != nil {
			j.errs = append(j.errs, e)
		}
	}
	if len(j.errs) == 0 {
		return nil
	}

	return &j
}

func (j *joined) Error() string {
	var b strings.Builder
	for i, e := range j.errs {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(e.Error())
	}
	return b.String()
}

func (j *joined) Unwrap() []error {
	return j.errs
}

func (j *joined) Format(f fmt.State, verb rune) {
	if j == nil {
		return
	}

	ffallback := "%" + string(verb)
	for _, e := range j.errs {
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
