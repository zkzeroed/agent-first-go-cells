package modulus

import "testing"

func TestReduce(t *testing.T) {
	tests := map[int]int{0: 0, 96: 96, 97: 0, 98: 1, -1: 96, -98: 96}
	for input, want := range tests {
		if got := Reduce(input); got != want {
			t.Errorf("Reduce(%d) = %d, want %d", input, got, want)
		}
	}
}
