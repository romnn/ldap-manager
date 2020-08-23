package hash

import (
	"testing"
)

// TestDESHash ...
func TestDESHash(t *testing.T) {
	pw := "password"
	hashed := encodeCRYPT(pw)
	expected := "{CRYPT}JQMuyS6H.AGMo"
	if hashed != expected {
		t.Errorf("expected encodeCRYPT(%q) == %q but got %s", pw, expected, hashed)
	}
}
