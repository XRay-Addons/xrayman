package xerr

type ErrType struct {
	string
}

func DefineErrType(t string) ErrType {
	return ErrType{string: t}
}
