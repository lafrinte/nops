package log

type Option func(l *Logging)

func WithLevel(val Level) Option {
	return func(l *Logging) {
		l.Level = val
	}
}

func WithMaxByte(val int) Option {
	return func(l *Logging) {
		l.MaxByte = val
	}
}

func WithColor() Option {
	return func(l *Logging) {
		l.NoColor = false
	}
}

func WithCompress() Option {
	return func(l *Logging) {
		l.Compress = true
	}
}
func WithUnixTimestamp() Option {
	return func(l *Logging) {
		l.UnixTimestamp = true
	}
}

func WithFileName(val string) Option {
	return func(l *Logging) {
		l.FileName = val
	}
}

func WithMaxSize(val int) Option {
	return func(l *Logging) {
		l.MaxSize = val
	}
}

func WithMaxBackups(val int) Option {
	return func(l *Logging) {
		l.MaxBackups = val
	}
}

func WithMaxAgeInDay(val int) Option {
	return func(l *Logging) {
		l.MaxAgeInDay = val
	}
}

func WithConsoleHandler() Option {
	return func(l *Logging) {
		l.EnableConsoleHandler = true
	}
}

func WithFileHandler() Option {
	return func(l *Logging) {
		l.EnableFileHandler = true
	}
}

func WithStreamHandler() Option {
	return func(l *Logging) {
		l.EnableStreamHandler = true
	}
}
