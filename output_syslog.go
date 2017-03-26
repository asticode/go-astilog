// +build !windows

package astilog

import (
	"log/syslog"
	"os"

	"github.com/rs/xlog"
)

// DefaultOutput is the default output
func DefaultOutput(c Configuration) xlog.Output {
	return xlog.NewConsoleOutputW(os.Stderr, xlog.NewSyslogOutputFacility("", "", c.AppName, syslog.LOG_LOCAL0))
}
