package greetingrender

import (
	"errors"
	"testing"

	"example.com/agent-first-library/internal/cells/greeting-render/api"
)

func TestRender(t *testing.T) {
	renderer := New(Deps{Normalizer: fakeNormalizer{}})
	message, err := renderer.Render("grace hopper")
	if err != nil || message != "Hello, Grace Hopper!" {
		t.Fatalf("Render() = %q, %v", message, err)
	}
	_, err = renderer.Render("")
	if !errors.Is(err, api.ErrInvalidName) {
		t.Fatalf("Render(empty) error = %v", err)
	}
}

type fakeNormalizer struct{}

func (fakeNormalizer) Normalize(name string) string {
	if name == "grace hopper" {
		return "Grace Hopper"
	}
	return ""
}
