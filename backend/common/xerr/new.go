package xerr

import "fmt"

func New(text string) error {
	var err error = &xerr{text: text}
	return Wrap(err, withStack(1))
}

func Newf(text string, args ...any) error {
	var err error = &xerr{text: fmt.Sprintf(text, args...)}
	return Wrap(err, withStack(1))
}
