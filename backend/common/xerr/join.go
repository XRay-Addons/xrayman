package xerr

import (
	"fmt"
	"strings"
)

type xjoined struct {
	errs []error
}

var _ error = (*xjoined)(nil)

func Join(errs ...error) error {
	j := xjoined{
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
	if len(j.errs) == 1 {
		return j.errs[0]
	}

	return &j
}

func (j *xjoined) Error() string {
	var b strings.Builder
	for i, e := range j.errs {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(e.Error())
	}
	return b.String()
}

func (j *xjoined) Unwrap() []error {
	return j.errs
}

func (j *xjoined) Format(f fmt.State, verb rune) {
	if j == nil {
		return
	}

	ffallback := "%" + string(verb)
	for i, e := range j.errs {
		if f.Flag('+') {
			fmt.Fprint(f, "\n\t* ")
		} else if i > 0 {
			fmt.Fprint(f, "; ")
		}
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
