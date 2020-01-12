// +build !windows

package astilog

import (
	"io"
	"log/syslog"
)

var newSyslogWriter = func(c Configuration) (io.WriteCloser, error) {
	return syslog.New(syslog.LOG_INFO|syslog.LOG_USER, c.AppName)
}
