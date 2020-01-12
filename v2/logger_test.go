package astilog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/asticode/go-astikit"
)

func TestNew(t *testing.T) {
	l := New(Configuration{AppName: "app"})
	defer l.Close()
	if e, g := map[string]interface{}{"app_name": "app"}, l.fs; !reflect.DeepEqual(e, g) {
		t.Errorf("expected %+v, got %+v", e, g)
	}
}

func TestSetWriter(t *testing.T) {
	// Create temp dir
	d, err := ioutil.TempDir("", "astilog_")
	if err != nil {
		t.Fatal(fmt.Errorf("creating temp dir failed: %w", err))
	}

	// Make sure to delete directory
	defer os.RemoveAll(d)

	// Create logger
	l := NewFromFlags()
	defer l.Close()

	// Default to stdout
	if !reflect.DeepEqual(l.w, astikit.NopCloser(os.Stdout)) {
		t.Error("expected false, got true")
	}

	// File
	f := filepath.Join(d, "f1.log")
	l.setWriter(Configuration{Filename: f})
	switch tp := l.w.(type) {
	case *os.File:
	default:
		t.Errorf("expected *os.File, got %T", tp)
	}
	l.w.Write([]byte("test"))
	b, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatal(fmt.Errorf("reading %s failed: %w", f, err))
	}
	if e := []byte("test"); !bytes.Equal(e, b) {
		t.Errorf("expected %s, got %s", e, b)
	}

	// File not working defaults to stdout
	l.setWriter(Configuration{Filename: filepath.Join("testdata/invalidpath")})
	if !reflect.DeepEqual(l.w, astikit.NopCloser(os.Stdout)) {
		t.Error("expected false, got true")
	}

	// Syslog
	l.setWriter(Configuration{Out: OutSyslog})
	switch tp := l.w.(type) {
	case *syslog.Writer:
	default:
		t.Errorf("expected *os.File, got %T", tp)
	}

	// Bypass newSyslogWriter
	old := newSyslogWriter
	newSyslogWriter = func(c Configuration) (io.WriteCloser, error) { return nil, errors.New("dummy") }
	defer func() { newSyslogWriter = old }()

	// syslog not working defaults to stdout
	l.setWriter(Configuration{Out: OutSyslog})
	if !reflect.DeepEqual(l.w, astikit.NopCloser(os.Stdout)) {
		t.Error("expected false, got true")
	}

	// Stderr
	l.setWriter(Configuration{Out: OutStderr})
	if !reflect.DeepEqual(l.w, astikit.NopCloser(os.Stderr)) {
		t.Error("expected false, got true")
	}
}

func TestSetLevel(t *testing.T) {
	l := NewFromFlags()
	defer l.Close()
	l.setLevel(Configuration{Level: LevelDebug})
	if e, g := levelDebug, l.l; e != g {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	l.setLevel(Configuration{Level: LevelInfo})
	if e, g := levelInfo, l.l; e != g {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	l.setLevel(Configuration{Level: LevelWarn})
	if e, g := levelWarn, l.l; e != g {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	l.setLevel(Configuration{Level: LevelError})
	if e, g := levelError, l.l; e != g {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	l.setLevel(Configuration{Level: LevelFatal})
	if e, g := levelFatal, l.l; e != g {
		t.Errorf("expected %+v, got %+v", e, g)
	}
}

func TestSetFormatter(t *testing.T) {
	l := NewFromFlags()
	defer l.Close()
	switch tp := l.f.(type) {
	case *textFormatter:
	default:
		t.Errorf("expected *textFormatter, got %T", tp)
	}
	l.setFormatter(Configuration{Format: FormatJSON}, time.Unix(0, 0))
	switch tp := l.f.(type) {
	case *jsonFormatter:
	default:
		t.Errorf("expected *jsonFormatter, got %T", tp)
	}
}

func TestSource(t *testing.T) {
	if e, g := "logger_test.go:138", source(); e != g {
		t.Errorf("expected %s, got %s", e, g)
	}
}

func TestWrite(t *testing.T) {
	b := &bytes.Buffer{}
	l := NewFromFlags()
	defer l.Close()
	l.w = astikit.NopCloser(b)

	// Level is not sufficient
	l.l = levelInfo
	l.write(context.Background(), msgFunc("test"), levelDebug)
	if e, g := "", b.String(); e != g {
		t.Errorf("expected %s, got %s", e, g)
	}

	// Context
	l.fs = map[string]interface{}{"k1": "v1"}
	l.write(ContextWithField(context.Background(), "k2", "v2"), msgFunc("test"), levelInfo)
	if e, g := " INFO[0000]test  k1=v1 k2=v2\n", b.String(); e != g {
		t.Errorf("expected %s, got %s", e, g)
	}

	// Source
	b.Reset()
	l.c.Source = true
	l.fs = map[string]interface{}{}
	l.write(context.Background(), msgFunc("test"), levelInfo)
	if e, g := " INFO[0000]test  source=logger_test.go:167\n", b.String(); e != g {
		t.Errorf("expected %s, got %s", e, g)
	}
}

func TestLogger(t *testing.T) {
	// Bypass exit
	old := exit
	count := 0
	exit = func() { count++ }
	defer func() { exit = old }()

	// Setup
	l := New(Configuration{Level: LevelDebug})
	defer l.Close()
	b := &bytes.Buffer{}
	l.w = astikit.NopCloser(b)
	ctx := ContextWithField(context.Background(), "k", "v")

	// Run
	l.Print("print")
	l.Printf("printf %s", "test")
	l.Debug("debug")
	l.Debugf("debug %s", "test")
	l.DebugC(ctx, "debug")
	l.DebugCf(ctx, "debug %s", "test")
	l.Info("info")
	l.Infof("info %s", "test")
	l.InfoC(ctx, "info")
	l.InfoCf(ctx, "info %s", "test")
	l.Warn("debug")
	l.Warnf("warn %s", "test")
	l.WarnC(ctx, "warn")
	l.WarnCf(ctx, "warn %s", "test")
	l.Error("error")
	l.Errorf("error %s", "test")
	l.ErrorC(ctx, "error")
	l.ErrorCf(ctx, "error %s", "test")
	l.Fatal("fatal")
	l.Fatalf("fatal %s", "test")
	l.FatalC(ctx, "fatal")
	l.FatalCf(ctx, "fatal %s", "test")

	// Assert
	if e, g := ` INFO[0000]print
 INFO[0000]printf test
DEBUG[0000]debug
DEBUG[0000]debug test
DEBUG[0000]debug  k=v
DEBUG[0000]debug test  k=v
 INFO[0000]info
 INFO[0000]info test
 INFO[0000]info  k=v
 INFO[0000]info test  k=v
 WARN[0000]debug
 WARN[0000]warn test
 WARN[0000]warn  k=v
 WARN[0000]warn test  k=v
ERROR[0000]error
ERROR[0000]error test
ERROR[0000]error  k=v
ERROR[0000]error test  k=v
FATAL[0000]fatal
FATAL[0000]fatal test
FATAL[0000]fatal  k=v
FATAL[0000]fatal test  k=v
`, b.String(); e != g {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e := 4; e != count {
		t.Errorf("expected %v, got %v", e, count)
	}

	// With field
	fs := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	l.WithFields(fs)
	if g := l.fs; !reflect.DeepEqual(fs, g) {
		t.Errorf("expected %+v, got %+v", fs, g)
	}
}
