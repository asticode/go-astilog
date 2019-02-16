package astilog

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// Logger represents a logger
type Logger interface {
	Clone() Logger
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	WithField(k string, v interface{})
	WithFields(fs Fields)
}

// Fields represents logger fields
type Fields map[string]interface{}

// LoggerSetter represents a logger setter
type LoggerSetter interface {
	SetLogger(l Logger)
}

// New creates a new Logger
func New(c Configuration) Logger {
	// Create logger
	l := newLogrus(c)

	// Default fields
	// Split this from the new method so that the Clone method only copies existing fields
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
