package greeting

import (
	"errors"

	greetingrender "example.com/agent-first-library/internal/cells/greeting-render"
	renderapi "example.com/agent-first-library/internal/cells/greeting-render/api"
	textnormalize "example.com/agent-first-library/internal/cells/text-normalize"
)

// ErrInvalidName reports an empty or whitespace-only consumer input.
var ErrInvalidName = errors.New("greeting: invalid name")

// Greeter creates a greeting for a caller-provided name.
type Greeter interface {
	Greet(string) (string, error)
}

type greeter struct{ renderer renderapi.Renderer }

// New constructs the public API over its private render cell.
func New() Greeter {
	return greeter{renderer: greetingrender.New(greetingrender.Deps{Normalizer: textnormalize.New()})}
}

func (g greeter) Greet(name string) (string, error) {
	message, err := g.renderer.Render(name)
	if errors.Is(err, renderapi.ErrInvalidName) {
		return "", ErrInvalidName
	}
	return message, err
}
