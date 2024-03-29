package log

var (
	defaultLogger *Logging
	consoleLogger *Logging
)

func init() {
	defaultLogger = DefaultLogger()
	consoleLogger = ConsoleLogger()
}

func New(opts ...Option) *Logging {
	logging := &Logging{
		Level: DebugLevel,
	}

	for _, opt := range opts {
		opt(logging)
	}

	if logging.FileName != "" {
		if logging.MaxSize == 0 {
			logging.MaxSize = 100
		}
	}

	return logging
}

func DefaultLogger() *Logging {
	if defaultLogger == nil {
		defaultLogger = New(WithLevel(DebugLevel), WithStreamHandler())
	}

	return defaultLogger
}

func ConsoleLogger() *Logging {
	if consoleLogger == nil {
		consoleLogger = New(WithLevel(DebugLevel), WithConsoleHandler())
	}

	return consoleLogger
}
