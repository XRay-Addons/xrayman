package xerr

var (
	ErrNilArg     = Define("nil arg")
	ErrNilCall    = Define("nil object call")
	ErrInvalidArg = Define("invalid argument")
	ErrPanic      = Define("panic")
)

func NilArg(name string) error {
	return Wrap(ErrNilArg,
		withStack(1),
		WithInfof("arg name: %s", name),
	)
}

func NilCall() error {
	return Wrap(ErrNilCall,
		withStack(1),
	)
}

// invalid argument error. contains argument, its value
// and optionally details from value check (details may be nil)
func InvalidArg(name string, value any, details error) error {
	options := []option{}
	if details != nil {
		options = append(options, WithInfof("details: %v", details))
	}
	options = append(options,
		WithInfof("argument '%s' value '%v'", name, value),
	)
	return Wrap(ErrInvalidArg, options...)
}

func Panic(p any) error {
	return Wrap(ErrPanic,
		withStack(1),
		WithInfof("panic info: %v", p),
	)
}
