package astilog

import (
	"bytes"
	"regexp"
)

// Vars
var (
	bytesEmpty  = []byte("")
	bytesEOL    = []byte("\n")
	regexpColor = regexp.MustCompile("\\[[\\d]+m")
)

// Writer represents an object capable of writing to the logger
type Writer struct {
	buffer *bytes.Buffer
	fn     func(v ...interface{})
}

// NewWriter creates a new writer
func NewWriter(fn func(v ...interface{})) *Writer {
	return &Writer{
		buffer: &bytes.Buffer{},
		fn:     fn,
	}
}

// Close closes the writer
func (w *Writer) Close() error {
	if w.buffer.Len() > 0 {
		w.write(w.buffer.Bytes())
	}
	return nil
}

// Write implements the io.Writer interface
func (w *Writer) Write(i []byte) (n int, err error) {
	// Update n to avoid broken pipe error
	defer func() {
		n = len(i)
	}()

	// No EOL in the log, write in buffer
	if bytes.Index(i, bytesEOL) == -1 {
		w.buffer.Write(i)
		return
	}

	// Loop in items split by EOL
	var items = bytes.Split(i, bytesEOL)
	for i := 0; i < len(items)-1; i++ {
		// If first item, add the buffer
		if i == 0 {
			items[i] = append(w.buffer.Bytes(), items[i]...)
			w.buffer.Reset()
		}

		// Log
		w.write(items[i])
	}

	// Add remaining to buffer
	w.buffer.Write(items[len(items)-1])
	return
}

func (w *Writer) write(i []byte) {
	// Sanitize text
	text := string(bytes.TrimSpace(regexpColor.ReplaceAll(i, bytesEmpty)))

	// Log
	w.fn(text)
}
