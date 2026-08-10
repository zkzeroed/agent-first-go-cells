package tokenissue

import "example.com/agent-first-reference/internal/cells/token-issue/api"

type Deps struct{ Store Store }

func New(deps Deps) api.Issuer { return service{store: deps.Store} }
