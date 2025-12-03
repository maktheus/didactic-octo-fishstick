package logger

import (
	"io"
	"log"
	"os"
)

// Logger exposes a minimal logging interface compatible with the standard library log package.
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// Option allows customizing the logger instance created by New.
type Option func(*log.Logger)

// WithWriter overrides the output writer used by the logger.
func WithWriter(w io.Writer) Option {
	return func(l *log.Logger) {
		l.SetOutput(w)
	}
}

// WithPrefix configures a prefix for all log lines.
func WithPrefix(prefix string) Option {
	return func(l *log.Logger) {
		l.SetPrefix(prefix)
	}
}

// New creates a new logger writing to stdout unless overridden through options.
func New(opts ...Option) Logger {
	l := log.New(os.Stdout, "backend ", log.LstdFlags|log.Lmicroseconds)
	for _, opt := range opts {
		opt(l)
	}
	return l
}
