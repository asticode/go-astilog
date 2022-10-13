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
	format(msg string, l astikit.LoggerLevel, fs map[string]interface{}) []byte
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

func (f *textFormatter) format(msg string, l astikit.LoggerLevel, fs map[string]interface{}) (b []byte) {
	// Add level
	switch l {
	case astikit.LoggerLevelDebug:
		b = append(b, []byte("DEBUG")...)
	case astikit.LoggerLevelWarn:
		b = append(b, []byte(" WARN")...)
	case astikit.LoggerLevelError:
		b = append(b, []byte("ERROR")...)
	case astikit.LoggerLevelFatal:
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
	b = append(b, newLine...)
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

func (f *jsonFormatter) format(msg string, l astikit.LoggerLevel, fs map[string]interface{}) []byte {
	// Add msg
	fs[f.msgKey] = msg

	// Add level
	fs["level"] = l

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
	b = append(b, newLine...)
	return b
}

type minimalistFormatter struct{}

func newMinimalistFormatter() *minimalistFormatter {
	return &minimalistFormatter{}
}

func (f *minimalistFormatter) format(msg string, l astikit.LoggerLevel, fs map[string]interface{}) []byte {
	return append([]byte(msg), newLine...)
}
