package tokenrender

import (
	issueapi "example.com/agent-first-reference/internal/cells/token-issue/api"
	"example.com/agent-first-reference/internal/cells/token-render/api"
)

type Deps struct{ Issuer issueapi.Issuer }

func NewHandler(deps Deps) api.Handler { return handler{issuer: deps.Issuer} }
