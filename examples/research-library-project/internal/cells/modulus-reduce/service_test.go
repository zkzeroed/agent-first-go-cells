package modulusreduce

import "testing"

func TestReduceCanonicalizesNegativeValues(t *testing.T) {
	if got := Reduce(-1, 97); got != 96 {
		t.Fatalf("Reduce(-1, 97) = %d, want 96", got)
	}
}
