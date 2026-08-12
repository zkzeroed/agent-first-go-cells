package textnormalize

import "testing"

func TestNormalize(t *testing.T) {
	if got := New().Normalize("  aDA   loVELACE "); got != "Ada Lovelace" {
		t.Fatalf("Normalize() = %q", got)
	}
}
