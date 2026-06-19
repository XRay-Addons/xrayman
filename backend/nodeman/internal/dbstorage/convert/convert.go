package convert

type withFn[From any, To any] = func(f *From, t *To)

type withEFn[From any, To any] = func(f *From, t *To) error

func cnv[From any, To any](from *From, w withFn[From, To], we ...withEFn[From, To]) (*To, error) {
	var to To
	if w != nil {
		w(from, &to)
	}
	for _, c := range we {
		if err := c(from, &to); err != nil {
			return nil, err
		}
	}
	return &to, nil
}

func cnvNoErr[From any, To any](from *From, w withFn[From, To]) *To {
	var to To
	if w != nil {
		w(from, &to)
	}

	return &to
}

func cnvArr[From any, To any](from []From, w withFn[From, To], we ...withEFn[From, To]) ([]To, error) {
	to := make([]To, len(from), len(from))

	for i := range from {
		t, err := cnv[From, To](&from[i], w, we...)
		if err != nil {
			return nil, err
		}
		to[i] = *t
	}

	return to, nil
}

func cnvArrNoErr[From any, To any](from []From, w withFn[From, To]) []To {
	to := make([]To, len(from), len(from))

	for i := range from {
		t := cnvNoErr[From, To](&from[i], w)
		to[i] = *t
	}

	return to
}
