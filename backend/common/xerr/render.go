package xerr

import (
	"strings"
)

func Render(err error) string {
	var b strings.Builder
	render(&b, err, 0, false, "")
	return b.String()
}

func render(b *strings.Builder,
	err error,
	indentSize int,
	detailsOnly bool,
	prefix string,
) {
	if err == nil {
		return
	}

	// print error
	if !detailsOnly {
		print(b, indentSize, prefix, err.Error())
		detailsOnly = true
	}

	// print details if exists
	if d, ok := err.(interface{ Details() string }); ok {
		print(b, indentSize, "with ", d.Details())
	}

	// unwrap
	switch u := err.(type) {
	case interface{ Unwrap() error }:
		render(b, u.Unwrap(), indentSize, detailsOnly, "")

	case interface{ Unwrap() []error }:
		for _, e := range u.Unwrap() {
			render(b, e, indentSize+1, false, "* ")
		}
	}
}

func print(b *strings.Builder, indentSize int, prefix string, text string) {
	indent := strings.Repeat("\t", indentSize)

	for i, line := range strings.Split(text, "\n") {
		b.WriteString(indent)
		if i == 0 {
			b.WriteString(prefix)
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
}
