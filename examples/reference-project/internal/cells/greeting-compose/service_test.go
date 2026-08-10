package greetingcompose

import (
	"errors"
	"testing"
)

func TestCompose(t *testing.T) {
	composer := New(Deps{Store: NewStore("Hello")})
	got, err := composer.Compose(t.Context(), " Ada ")
	if err != nil || got != "Hello, Ada!" {
		t.Fatalf("Compose() = %q, %v", got, err)
	}
	_, err = composer.Compose(t.Context(), " ")
	if !errors.Is(err, ErrEmptyName) {
		t.Fatalf("Compose(empty) error = %v", err)
	}
}
