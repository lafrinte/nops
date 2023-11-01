package log

import (
	"fmt"
	"github.com/go-lumberjack/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"io"
	"os"
)

const (
	// KB ...
	KB = 1024 << (10 * iota)
	// MB ...
	MB
	// GB ...
	GB
)

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled

	// TraceLevel defines trace log level.
	TraceLevel Level = -1
)

type Level int8

type Logging struct {
	// global config
	Level         Level
	NoColor       bool
	UnixTimestamp bool

	// file handler config
	FileName    string
	MaxByte     int
	MaxSize     int
	MaxBackups  int
	MaxAgeInDay int
	Compress    bool

	// handler config
	EnableFileHandler    bool
	EnableStreamHandler  bool
	EnableConsoleHandler bool
}

func (l *Logging) getHandler() []io.Writer {
	zerolog.SetGlobalLevel(zerolog.Level(l.Level))
	zerolog.MessageFieldName = "msg"
	zerolog.TimestampFieldName = "tz"

	if l.UnixTimestamp {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	var handlers []io.Writer

	// ConsoleHandler conflict with StreamHandler and FileHandler.
	if l.EnableConsoleHandler {
		handlers = append(handlers, l.ConsoleHandler())
		return handlers
	}

	if l.EnableFileHandler {
		handlers = append(handlers, l.FileHandler())
	}

	if l.EnableStreamHandler {
		handlers = append(handlers, l.StreamHandler())
	}

	return handlers
}

func (l *Logging) FileHandler() io.Writer {
	return &lumberjack.Logger{
		Filename:   l.FileName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}
}

func (l *Logging) StreamHandler() io.Writer {
	return diode.NewWriter(os.Stderr, 10000, 0, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	})
}

func (l *Logging) ConsoleHandler() io.Writer {
	writer := zerolog.NewConsoleWriter()
	writer.NoColor = l.NoColor

	return writer
}

func (l *Logging) GetLogger() zerolog.Logger {
	writers := zerolog.MultiLevelWriter(l.getHandler()...)
	return zerolog.New(writers).With().Timestamp().Logger()
}
