package satellite

import (
	"testing"
)

func TestDetermineCSP(t *testing.T) {
	t.Logf("CSP: %s", DetermineCSP())
}
