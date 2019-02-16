package astilog

import (
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Logrus represents a logrus logger
type Logrus struct {
	c      Configuration
	fields Fields
	l      *logrus.Logger
}

func newLogrus(c Configuration) (l *Logrus) {
	// Init
	l = &Logrus{
		c:      c,
		fields: make(Fields),
		l:      logrus.New(),
	}

	// Out
	var out string
	l.l.Out, out = logrusOut(c)

	// Formatter
	l.l.Formatter = logrusFormatter(c, out)

	// Level
	l.l.Level = logrusLevel(c)
	return
}

func logrusOut(c Configuration) (w io.Writer, out string) {
	switch c.Out {
	case OutStdOut:
		return stdOut(), c.Out
	case OutSyslog:
		return syslogOut(c), c.Out
	default:
		if isTerminal(os.Stdout) {
			w = stdOut()
			out = OutStdOut
		} else {
			w = syslogOut(c)
			out = OutSyslog
		}
		if len(c.Filename) > 0 {
			f, err := os.OpenFile(c.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Println(errors.Wrapf(err, "astilog: creating %s failed", c.Filename))
			} else {
				w = f
				out = OutFile
			}
		}
		return
	}
}

func logrusFormatter(c Configuration, out string) logrus.Formatter {
	switch c.Format {
	case FormatJSON:
		return logrusJSONFormatter(c)
	case FormatText:
		return logrusTextFormatter(c, out)
	default:
		switch out {
		case OutFile, OutStdOut:
			return logrusTextFormatter(c, out)
		default:
			return logrusJSONFormatter(c)
		}
	}
}

func logrusJSONFormatter(c Configuration) logrus.Formatter {
	f := &logrus.JSONFormatter{
		FieldMap:        make(logrus.FieldMap),
		TimestampFormat: c.TimestampFormat,
	}
	if len(c.MessageKey) > 0 {
		f.FieldMap[logrus.FieldKeyMsg] = c.MessageKey
	}
	return f
}

func logrusTextFormatter(c Configuration, out string) logrus.Formatter {
	return &logrus.TextFormatter{
		DisableColors:   c.DisableColors || out == OutFile,
		ForceColors:     !c.DisableColors && out != OutFile,
		FullTimestamp:   c.FullTimestamp,
		TimestampFormat: c.TimestampFormat,
	}
}

func logrusLevel(c Configuration) logrus.Level {
	if c.Verbose {
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}

// Clone implements the Logger interface
func (l *Logrus) Clone() Logger {
	// Create logger
	n := newLogrus(l.c)

	// Copy fields
	n.WithFields(l.fields)
	return n
}

// Debug implements the Logger interface
func (l *Logrus) Debug(v ...interface{}) { l.l.Debug(v...) }

// Debugf implements the Logger interface
func (l *Logrus) Debugf(format string, v ...interface{}) { l.l.Debugf(format, v...) }

// WithField implements the Logger interface
func (l *Logrus) Info(v ...interface{}) { l.l.Info(v...) }

// WithField implements the Logger interface
func (l *Logrus) Infof(format string, v ...interface{}) { l.l.Infof(format, v...) }

// WithField implements the Logger interface
func (l *Logrus) Warn(v ...interface{}) { l.l.Warn(v...) }

// WithField implements the Logger interface
func (l *Logrus) Warnf(format string, v ...interface{}) { l.l.Warnf(format, v...) }

// WithField implements the Logger interface
func (l *Logrus) Error(v ...interface{}) { l.l.Error(v...) }

// WithField implements the Logger interface
func (l *Logrus) Errorf(format string, v ...interface{}) { l.l.Errorf(format, v...) }

// WithField implements the Logger interface
func (l *Logrus) Fatal(v ...interface{}) { l.l.Fatal(v...) }

// WithField implements the Logger interface
func (l *Logrus) Fatalf(format string, v ...interface{}) { l.l.Fatalf(format, v...) }

// WithField implements the Logger interface
func (l *Logrus) WithField(k string, v interface{}) {
	l.fields[k] = v
	l.l.AddHook(newWithFieldHook(k, v))
}

// WithFields implements the Logger interface
func (l *Logrus) WithFields(fs Fields) {
	for k, v := range fs {
		l.WithField(k, v)
	}
}
