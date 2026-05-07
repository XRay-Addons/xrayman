package errdefs

/*func WithOgen() option {
	return func(e *baseError) {
		// add status code if exists
		e.with = append(e.with, e.err.Error())
		var sc interface{ StatusCode() int }
		if errors.As(e.err, &sc) {
			e.with = append(e.with, fmt.Sprintf("Status: %d", sc.StatusCode()))
		} else {
			e.with = append(e.with, "Status: Transport error")
		}
		// add url path if exists
		var ue *url.Error
		if errors.As(e.err, &ue) {
			e.with = append(e.with, fmt.Sprintf("URL: %s", ue.URL))
		}
		// replace error
		e.err = ErrConnection
	}
}*/
