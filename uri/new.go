package uri

func New(opts ...Option) *URI {
	d := &URI{}

	for _, opt := range opts {
		opt(d)
	}

	d.parsing()

	return d
}
