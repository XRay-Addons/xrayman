package dbstorage

type ConvertOption[From any, To any] = func(from *From, to *To) error

func With[From any, To any](f func(from *From, to *To)) ConvertOption[From, To] {
	return func(from *From, to *To) error {
		f(from, to)
		return nil
	}
}

func WithE[From any, To any](f func(from *From, to *To) error) ConvertOption[From, To] {
	return f
}

func Convert[From any, To any](from *From, opts ...ConvertOption[From, To]) (*To, error) {
	var to To
	for _, opt := range opts {
		if err := opt(from, &to); err != nil {
			return nil, err
		}
	}

	return &to, nil
}

func ConvertArray[From any, To any](from []From, opts ...ConvertOption[From, To]) ([]To, error) {
	to := make([]To, len(from), len(from))

	for i := range from {
		t, err := Convert[From, To](&from[i], opts...)
		if err != nil {
			return nil, err
		}
		to[i] = *t
	}

	return to, nil
}
