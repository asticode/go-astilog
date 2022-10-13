package astilog

import (
	"testing"

	"github.com/asticode/go-astikit"
)

func TestConfiguration(t *testing.T) {
	*Level = "info"
	*Verbose = true
	if e, g := astikit.LoggerLevelDebug, FlagConfig().Level; e != g {
		t.Errorf("expected %+v, got %+v", e, g)
	}
}
