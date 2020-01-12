package astilog

import "testing"

func TestConfiguration(t *testing.T) {
	*Level = LevelInfo
	*Verbose = true
	if e, g := LevelDebug, FlagConfig().Level; e != g {
		t.Errorf("expected %+v, got %+v", e, g)
	}
}
