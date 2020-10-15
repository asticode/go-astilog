package astilog

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asticode/go-astikit"
)

type formatter interface {
	format(msg string, level int, fs map[string]interface{}) []byte
}

type textFormatter struct {
	c         Configuration
	createdAt time.Time
}

func newTextFormatter(c Configuration, createdAt time.Time) *textFormatter {
	return &textFormatter{
		c:         c,
		createdAt: createdAt,
	}
}

func (f *textFormatter) format(msg string, level int, fs map[string]interface{}) (b []byte) {
	// Add level
	switch level {
	case levelDebug:
		b = append(b, []byte("DEBUG")...)
	case levelWarn:
		b = append(b, []byte(" WARN")...)
	case levelError:
		b = append(b, []byte("ERROR")...)
	case levelFatal:
		b = append(b, []byte("FATAL")...)
	default:
		b = append(b, []byte(" INFO")...)
	}

	// Add timestamp
	b = append(b, []byte("[")...)
	if f.c.TimestampFormat == "" {
		b = append(b, astikit.BytesPad([]byte(strconv.Itoa(int(now().Sub(f.createdAt).Seconds()))), '0', 4)...)
	} else {
		b = append(b, []byte(now().Format(f.c.TimestampFormat))...)
	}
	b = append(b, []byte("]")...)

	// Add msg
	b = append(b, []byte(msg)...)

	// Add fields
	if len(fs) > 0 {
		// Add spaces
		b = append(b, []byte("  ")...)

		// Sort fields
		var vs []string
		for k, v := range fs {
			vs = append(vs, k+"="+fmt.Sprintf("%v", v))
		}
		sort.Strings(vs)
		b = append(b, []byte(strings.Join(vs, " "))...)
	}

	// Add newline
	b = append(b, []byte("\n")...)
	return
}

type jsonFormatter struct {
	c         Configuration
	createdAt time.Time
	msgKey    string
}

func newJSONFormatter(c Configuration, createdAt time.Time) (f *jsonFormatter) {
	f = &jsonFormatter{
		c:         c,
		createdAt: createdAt,
		msgKey:    "msg",
	}
	if c.MessageKey != "" {
		f.msgKey = c.MessageKey
	}
	return
}

func (f *jsonFormatter) format(msg string, level int, fs map[string]interface{}) []byte {
	// Add msg
	fs[f.msgKey] = msg

	// Add level
	switch level {
	case levelDebug:
		fs["level"] = LevelDebug
	case levelError:
		fs["level"] = LevelError
	case levelFatal:
		fs["level"] = LevelFatal
	case levelWarn:
		fs["level"] = LevelWarn
	default:
		fs["level"] = LevelInfo
	}

	// Add timestamp
	if f.c.TimestampFormat == "" {
		fs["time"] = int(now().Sub(f.createdAt).Seconds())
	} else {
		fs["time"] = now().Format(f.c.TimestampFormat)
	}

	// Marshal
	b, err := json.Marshal(fs)
	if err != nil {
		log.Println(fmt.Errorf("astilog: marshaling failed: %w", err))
		return nil
	}

	// Add newline
	b = append(b, []byte("\n")...)
	return b
}

type minimalistFormatter struct{}

func newMinimalistFormatter() *minimalistFormatter {
	return &minimalistFormatter{}
}

func (f *minimalistFormatter) format(msg string, level int, fs map[string]interface{}) []byte {
	return append([]byte(msg), []byte("\n")...)
}
