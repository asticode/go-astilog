package astilog

import "github.com/rs/xlog"

// NopLogger returns a nop logger
func NopLogger() Logger {
	return xlog.NopLogger
}

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

// NewXlogConfig creates a new xlog.Config
func NewXlogConfig(c Configuration) (o xlog.Config) {
	// Init
	o = xlog.Config{
		Fields: xlog.F{
			"app_name": c.AppName,
		},
		Level:  xlog.LevelInfo,
		Output: DefaultOutput(c),
	}

	// Verbose
	if c.Verbose {
		o.Level = xlog.LevelDebug
	}
	return
}

// New creates a new Logger
func New(c Configuration) Logger {
	return xlog.New(NewXlogConfig(c))
}
