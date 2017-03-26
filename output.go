// +build windows

package astilog

import (
	"os"

	"github.com/rs/xlog"
)

// DefaultOutput is the default output
func DefaultOutput(c Configuration) xlog.Output {
	return xlog.NewConsoleOutputW(os.Stderr, nil)
}
