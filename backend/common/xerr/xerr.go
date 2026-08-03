package xerr

type xerr struct {
	text string
}

func (d *xerr) Error() string {
	return d.text
}
