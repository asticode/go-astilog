package astilog

import (
	"io"
	"os"

	"context"

	"golang.org/x/crypto/ssh/terminal"
)

// Logger represents a logger
type Logger interface {
	Debug(v ...interface{})
	DebugC(ctx context.Context, v ...interface{})
	DebugCf(ctx context.Context, format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	InfoC(ctx context.Context, v ...interface{})
	InfoCf(ctx context.Context, format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	WarnC(ctx context.Context, v ...interface{})
	WarnCf(ctx context.Context, format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	ErrorC(ctx context.Context, v ...interface{})
	ErrorCf(ctx context.Context, format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	FatalC(ctx context.Context, v ...interface{})
	FatalCf(ctx context.Context, format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	WithField(k string, v interface{})
	WithFields(fs Fields)
}

// LoggerSetter represents a logger setter
type LoggerSetter interface {
	SetLogger(l Logger)
}

// New creates a new Logger
func New(c Configuration) Logger {
	// Create logger
	l := newLogrus(c)

	// Default fields
	l.WithFields(Fields{"app_name": c.AppName})
	return l
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
