package canonical

import "testing"

func TestFold(t *testing.T) {
	if got := Fold("  Agent   First "); got != "agent first" {
		t.Fatalf("Fold() = %q, want %q", got, "agent first")
	}
}
