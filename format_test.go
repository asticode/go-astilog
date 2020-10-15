package astilog

import (
	"bytes"
	"testing"
	"time"
)

func TestTextFormatter(t *testing.T) {
	oldNow := now
	defer func() { now = oldNow }()
	now = func() time.Time { return time.Unix(5, 0).UTC() }

	f := newTextFormatter(Configuration{}, time.Unix(0, 0).UTC())
	if e, g := []byte("DEBUG[0005]msg  k1=v1 k2=v2\n"), f.format("msg", levelDebug, map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e, g := []byte(" INFO[0005]msg\n"), f.format("msg", levelInfo, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e, g := []byte(" WARN[0005]msg\n"), f.format("msg", levelWarn, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e, g := []byte("ERROR[0005]msg\n"), f.format("msg", levelError, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e, g := []byte("FATAL[0005]msg\n"), f.format("msg", levelFatal, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}

	f = newTextFormatter(Configuration{TimestampFormat: time.RFC3339}, time.Unix(0, 0))
	if e, g := []byte(" INFO[1970-01-01T00:00:05Z]msg\n"), f.format("msg", levelInfo, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
}

func TestJSONFormatter(t *testing.T) {
	oldNow := now
	defer func() { now = oldNow }()
	now = func() time.Time { return time.Unix(5, 0).UTC() }

	f := newJSONFormatter(Configuration{}, time.Unix(0, 0).UTC())
	if e, g := []byte(`{"k1":"v1","k2":"v2","level":"debug","msg":"msg","time":5}`+"\n"), f.format("msg", levelDebug, map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e, g := []byte(`{"level":"info","msg":"msg","time":5}`+"\n"), f.format("msg", levelInfo, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e, g := []byte(`{"level":"warn","msg":"msg","time":5}`+"\n"), f.format("msg", levelWarn, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e, g := []byte(`{"level":"error","msg":"msg","time":5}`+"\n"), f.format("msg", levelError, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
	if e, g := []byte(`{"level":"fatal","msg":"msg","time":5}`+"\n"), f.format("msg", levelFatal, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}

	f = newJSONFormatter(Configuration{
		MessageKey:      "msg_test",
		TimestampFormat: time.RFC3339,
	}, time.Unix(0, 0))
	if e, g := []byte(`{"level":"info","msg_test":"msg","time":"1970-01-01T00:00:05Z"}`+"\n"), f.format("msg", levelInfo, map[string]interface{}{}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
}

func TestMinimalistFormatter(t *testing.T) {
	f := newMinimalistFormatter()
	if e, g := []byte("msg\n"), f.format("msg", levelDebug, map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}); !bytes.Equal(e, g) {
		t.Errorf("expected %s, got %s", e, g)
	}
}
