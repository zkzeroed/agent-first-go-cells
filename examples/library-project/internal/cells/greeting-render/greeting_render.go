package greetingrender

import (
	"fmt"

	"example.com/agent-first-library/internal/cells/greeting-render/api"
	normalizeapi "example.com/agent-first-library/internal/cells/text-normalize/api"
)

type renderer struct{ normalizer normalizeapi.Normalizer }

// Deps supplies the declared private-cell contracts.
type Deps struct{ Normalizer normalizeapi.Normalizer }

// New constructs the private renderer over its explicit cell dependency.
func New(deps Deps) api.Renderer { return renderer{normalizer: deps.Normalizer} }

func (r renderer) Render(name string) (string, error) {
	name = r.normalizer.Normalize(name)
	if name == "" {
		return "", api.ErrInvalidName
	}
	return fmt.Sprintf("Hello, %s!", name), nil
}
