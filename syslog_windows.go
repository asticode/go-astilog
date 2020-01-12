package astilog

import (
	"errors"
	"io"
)

var newSyslogWriter = func(c Configuration) (io.WriteCloser, error) {
	return nil, errors.New("astilog: syslog is not implemented")
}
