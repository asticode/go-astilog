// +build !windows

package astilog

import (
	"io"
	"log/syslog"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// DefaultOutput is the default output
func DefaultOut(c Configuration) (w io.Writer) {
	if logrus.IsTerminal(os.Stdout) {
		return os.Stdout
	}
	var err error
	if w, err = syslog.New(syslog.LOG_LOCAL0, c.AppName); err != nil {
		panic(errors.Wrap(err, "new syslog failed"))
	}
	return
}
