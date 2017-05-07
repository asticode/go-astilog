// +build !windows

package astilog

import (
	"io"
	"log/syslog"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

// DefaultOutput is the default output
func DefaultOut(c Configuration) (w io.Writer) {
	if logrus.IsTerminal() {
		return os.Stdout
	}
	var err error
	if w, err = syslog.New(syslog.LOG_LOCAL0, c.AppName); err != nil {
		panic(errors.Wrap(err, "new syslog failed"))
	}
	return
}
