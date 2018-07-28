package astilog

import (
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

// Logger represents a logger
type Logger interface {
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
}

// LoggerSetter represents a logger setter
type LoggerSetter interface {
	SetLogger(l Logger)
}

// LoggerWithField represents a logger that can have fields
type LoggerWithFields interface {
	WithField(k, v string)
	WithFields(fs Fields)
}

type logger struct {
	*logrus.Logger
}

func newLogger() *logger {
	return &logger{Logger: logrus.New()}
}

// WithField implements the LoggerWithFields interface
func (l *logger) WithField(k, v string) {
	l.AddHook(newWithFieldHook(k, v))
}

// Fields represents logger fields
type Fields map[string]string

// WithFields implements the LoggerWithFields interface
func (l *logger) WithFields(fs Fields) {
	for k, v := range fs {
		l.WithField(k, v)
	}
}

// New creates a new Logger
func New(c Configuration) Logger {
	// Init
	var l = newLogger()

	// Hooks
	l.AddHook(newWithFieldHook("app_name", c.AppName))

	// Out
	l.Out = DefaultOut(c)
	if len(c.Filename) > 0 {
		f, err := os.OpenFile(c.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(errors.Wrapf(err, "creating %s failed", c.Filename))
		} else {
			l.Out = f
		}
	}

	// Formatter
	l.Formatter = &logrus.TextFormatter{ForceColors: true}
	if !isTerminal(l.Out) {
		if len(c.Filename) > 0 {
			l.Formatter = &logrus.TextFormatter{DisableColors: true}
		} else {
			f := &logrus.JSONFormatter{FieldMap: make(logrus.FieldMap)}
			if len(c.MessageKey) > 0 {
				f.FieldMap[logrus.FieldKeyMsg] = c.MessageKey
			}
			l.Formatter = f
		}
	}

	// Level
	l.Level = logrus.InfoLevel
	if c.Verbose {
		l.Level = logrus.DebugLevel
	}
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
