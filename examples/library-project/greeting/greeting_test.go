package greeting

import (
	"errors"
	"testing"
)

func TestGreeter(t *testing.T) {
	greeter := New()
	message, err := greeter.Greet(" ada lovelace ")
	if err != nil || message != "Hello, Ada Lovelace!" {
		t.Fatalf("Greet() = %q, %v", message, err)
	}
	_, err = greeter.Greet(" \t")
	if !errors.Is(err, ErrInvalidName) {
		t.Fatalf("Greet(empty) error = %v", err)
	}
}
