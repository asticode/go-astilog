package astilog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/asticode/go-astikit"
)

var newLine = []byte("\n")

// Logger represents an object that can log stuff
type Logger struct {
	c         Configuration
	createdAt time.Time
	f         formatter
	fs        map[string]interface{}
	mf        *sync.RWMutex       // Locks fs
	l         astikit.LoggerLevel // Level
	w         io.WriteCloser
}

// NewFromFlags creates a new Logger based on flags
func NewFromFlags() (l *Logger) {
	return New(FlagConfig())
}

var now = func() time.Time { return time.Now() }

// New creates a new Logger
func New(c Configuration) (l *Logger) {
	// Create
	l = &Logger{
		c:         c,
		createdAt: now(),
		fs:        make(map[string]interface{}),
		mf:        &sync.RWMutex{},
	}

	// Add app name field
	if c.AppName != "" {
		l.WithField("app_name", c.AppName)
	}

	// Set writer
	l.setWriter(c)

	// Set level
	l.setLevel(c)

	// Set formatter
	l.setFormatter(c, l.createdAt)
	return
}

// Close closes the logger properly
func (l *Logger) Close() error {
	return l.w.Close()
}

func (l *Logger) setWriter(c Configuration) {
	// File
	if c.Filename != "" {
		// Open file
		f, err := os.OpenFile(c.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err == nil {
			l.w = f
			return
		}

		// Revert to default
		c.Out = ""
		log.Println(fmt.Errorf("astilog: creating %s failed: %w", c.Filename, err))
	}

	// Syslog
	if c.Out == OutSyslog {
		// Create
		s, err := newSyslogWriter(c)
		if err == nil {
			l.w = s
			return
		}

		// Revert to default
		c.Out = ""
		log.Println(fmt.Errorf("astilog: creating syslog failed: %w", err))
	}

	// Stderr
	if c.Out == OutStderr {
		l.w = astikit.NopCloser(os.Stderr)
		return
	}

	// Default is stdout
	l.w = astikit.NopCloser(os.Stdout)
}

func (l *Logger) setLevel(c Configuration) {
	l.l = c.Level
}

func (l *Logger) setFormatter(c Configuration, createdAt time.Time) {
	switch c.Format {
	case FormatJSON:
		l.f = newJSONFormatter(c, createdAt)
	case FormatMinimalist:
		l.f = newMinimalistFormatter()
	default:
		l.f = newTextFormatter(c, createdAt)
	}
}

func source() string {
	// Skip self callers
	i := 0
	_, file, line, ok := runtime.Caller(i)
	for ok && strings.Contains(file, "asticode/go-astilog/logger.go") {
		i++
		_, file, line, ok = runtime.Caller(i)
	}

	// Process file
	if !ok {
		file = "<???>"
		line = 1
	} else {
		file = filepath.Base(file)
	}
	return file + ":" + strconv.Itoa(line)
}

func (l *Logger) write(ctx context.Context, msgFunc func() string, lvl astikit.LoggerLevel) {
	// Check level
	if l.l > lvl {
		return
	}

	// Create fields
	fs := make(map[string]interface{})
	l.mf.RLock()
	for k, v := range l.fs {
		fs[k] = v
	}
	l.mf.RUnlock()

	// Add source
	if l.c.Source {
		fs["source"] = source()
	}

	// Add context fields
	if cfs := fieldsFromContext(ctx); cfs != nil {
		cfs.m.Lock()
		for k, v := range cfs.fs {
			fs[k] = v
		}
		cfs.m.Unlock()
	}

	// Format message
	m := l.f.format(msgFunc(), lvl, fs)

	// Write
	if l.c.MaxWriteLength > 0 && len(m) > l.c.MaxWriteLength {
		// Loop
		var c int
		for {
			// Get boundaries
			from := c * l.c.MaxWriteLength
			to := (c + 1) * l.c.MaxWriteLength

			// We're done
			if from > len(m)-1 {
				break
			}

			// We've reached the end of the message
			if to > len(m) {
				to = len(m)
			}

			// Since append modifies the input slice, we need to create a new one when
			// appending a new line
			wm := m[from:to]
			if to != len(m) {
				wm = make([]byte, to-from)
				copy(wm, m[from:to])
				if !bytes.HasSuffix(wm, newLine) {
					wm = append(wm, newLine...)
				}
			}

			// Write
			if _, err := l.w.Write(wm); err != nil {
				log.Println(fmt.Errorf("astilog: writing failed: %w", err))
				return
			}

			// Increment
			c++
		}
	} else {
		// Write
		if _, err := l.w.Write(m); err != nil {
			log.Println(fmt.Errorf("astilog: writing failed: %w", err))
			return
		}
	}
}

func msgFunc(v ...interface{}) func() string {
	return func() string { return fmt.Sprint(v...) }
}

func msgFuncf(format string, v ...interface{}) func() string {
	return func() string { return fmt.Sprintf(format, v...) }
}

func (l *Logger) Print(v ...interface{}) {
	l.Info(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), astikit.LoggerLevelDebug)
}

func (l *Logger) DebugC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), astikit.LoggerLevelDebug)
}

func (l *Logger) DebugCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), astikit.LoggerLevelDebug)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), astikit.LoggerLevelDebug)
}

func (l *Logger) Info(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), astikit.LoggerLevelInfo)
}

func (l *Logger) InfoC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), astikit.LoggerLevelInfo)
}

func (l *Logger) InfoCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), astikit.LoggerLevelInfo)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), astikit.LoggerLevelInfo)
}

func (l *Logger) Warn(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), astikit.LoggerLevelWarn)
}

func (l *Logger) WarnC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), astikit.LoggerLevelWarn)
}

func (l *Logger) WarnCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), astikit.LoggerLevelWarn)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), astikit.LoggerLevelWarn)
}

func (l *Logger) Error(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), astikit.LoggerLevelError)
}

func (l *Logger) ErrorC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), astikit.LoggerLevelError)
}

func (l *Logger) ErrorCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), astikit.LoggerLevelError)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), astikit.LoggerLevelError)
}

var exit = func() { os.Exit(1) }

func (l *Logger) Fatal(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), astikit.LoggerLevelFatal)
	exit()
}

func (l *Logger) FatalC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), astikit.LoggerLevelFatal)
	exit()
}

func (l *Logger) FatalCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), astikit.LoggerLevelFatal)
	exit()
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), astikit.LoggerLevelFatal)
	exit()
}

func (l *Logger) Write(lv astikit.LoggerLevel, v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), lv)
}

func (l *Logger) WriteC(ctx context.Context, lv astikit.LoggerLevel, v ...interface{}) {
	l.write(ctx, msgFunc(v...), lv)
}

func (l *Logger) WriteCf(ctx context.Context, lv astikit.LoggerLevel, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), lv)
}

func (l *Logger) Writef(lv astikit.LoggerLevel, format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), lv)
}

// WithField adds a field to the logger
func (l *Logger) WithField(k string, v interface{}) {
	l.mf.Lock()
	l.fs[k] = v
	l.mf.Unlock()
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fs map[string]interface{}) {
	for k, v := range fs {
		l.WithField(k, v)
	}
}
