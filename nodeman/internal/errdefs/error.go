package errdefs

import (
	"fmt"
	"strings"

	"github.com/go-faster/errors"
)

// - short form of wrap errors
// 	WrapWith(err, "details")
// 	WrapWithf(err, "details %s", "str")
//  WrapWithFile(err, path),
//  WrapWithStack(err, path)
//
// - wrap errors
// Wrap(err error,
// 	With("details"),
// 	Withf("details %s", "str"),
// 	WithFile(path),
//	WithStack())
//
// - make new errors:
// New(text,
// 	With("details"),
// 	Withf("details %s", "str"),
// 	WithFile(path),
//	WithStack()) (on by default)
//  WithoutStack()
//
// - common error types:
// NewNilCall for nil object call
// NewNilArg(name string) for nil arg passed

type baseError struct {
	err   error
	with  []string
	stack []string
}

var _ error = (*baseError)(nil)
var _ errors.Wrapper = (*baseError)(nil)

func (b *baseError) Error() string {
	if b == nil || b.err == nil {
		return ""
	}
	return b.err.Error()
}

func (b *baseError) Format(f fmt.State, verb rune) {
	if b == nil || b.err == nil {
		return
	}

	format := "%" + string(verb)
	fmt.Fprintf(f, format, b.err.Error())
	if f.Flag('+') {
		fmt.Fprintf(f, format, b.details())
	}
}

func (b *baseError) details() string {
	text := ""
	if len(b.stack) > 0 {
		text = fmt.Sprintf("-> %s:\n\t%s", strings.Join(b.stack, "\n-> "), text)
	}
	if len(b.with) > 0 {
		text = fmt.Sprintf("%s\nwith %s", text, strings.Join(b.with, "\nwith "))
	}
	return text
}

func (b *baseError) Unwrap() error {
	if b == nil {
		return nil
	}
	return b.err
}
