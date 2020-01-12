package astilog

import (
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

// Levels
const (
	levelDebug = iota
	levelInfo
	levelWarn
	levelError
	levelFatal
)

// Logger represents an object that can log stuff
type Logger struct {
	c         Configuration
	createdAt time.Time
	f         formatter
	fs        map[string]interface{}
	mf        *sync.RWMutex // Locks fs
	l         int           // Level
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
	switch c.Level {
	case LevelDebug:
		l.l = levelDebug
	case LevelWarn:
		l.l = levelWarn
	case LevelError:
		l.l = levelError
	case LevelFatal:
		l.l = levelFatal
	default:
		l.l = levelInfo
	}
}

func (l *Logger) setFormatter(c Configuration, createdAt time.Time) {
	switch c.Format {
	case FormatJSON:
		l.f = newJSONFormatter(c, createdAt)
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

func (l *Logger) write(ctx context.Context, msgFunc func() string, level int) {
	// Check level
	if l.l > level {
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
	for k, v := range fieldsFromContext(ctx) {
		fs[k] = v
	}

	// Write
	if _, err := l.w.Write(l.f.format(msgFunc(), level, fs)); err != nil {
		log.Println(fmt.Errorf("astilog: writing failed: %w", err))
		return
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
	l.write(context.Background(), msgFunc(v...), levelDebug)
}

func (l *Logger) DebugC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), levelDebug)
}

func (l *Logger) DebugCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), levelDebug)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), levelDebug)
}

func (l *Logger) Info(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), levelInfo)
}

func (l *Logger) InfoC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), levelInfo)
}

func (l *Logger) InfoCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), levelInfo)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), levelInfo)
}

func (l *Logger) Warn(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), levelWarn)
}

func (l *Logger) WarnC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), levelWarn)
}

func (l *Logger) WarnCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), levelWarn)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), levelWarn)
}

func (l *Logger) Error(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), levelError)
}

func (l *Logger) ErrorC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), levelError)
}

func (l *Logger) ErrorCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), levelError)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), levelError)
}

var exit = func() { os.Exit(1) }

func (l *Logger) Fatal(v ...interface{}) {
	l.write(context.Background(), msgFunc(v...), levelFatal)
	exit()
}

func (l *Logger) FatalC(ctx context.Context, v ...interface{}) {
	l.write(ctx, msgFunc(v...), levelFatal)
	exit()
}

func (l *Logger) FatalCf(ctx context.Context, format string, v ...interface{}) {
	l.write(ctx, msgFuncf(format, v...), levelFatal)
	exit()
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.write(context.Background(), msgFuncf(format, v...), levelFatal)
	exit()
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
