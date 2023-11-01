package log

func New(opts ...Option) *Logging {
	logging := &Logging{
		Level: DebugLevel,
	}

	for _, opt := range opts {
		opt(logging)
	}

	return logging
}
