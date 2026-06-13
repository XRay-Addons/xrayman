package xerr

func Define(name string) error {
	return &xerr{text: name}
}
