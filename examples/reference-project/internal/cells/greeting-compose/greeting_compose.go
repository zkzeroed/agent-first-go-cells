package greetingcompose

import "example.com/agent-first-reference/internal/cells/greeting-compose/api"

type Deps struct{ Store Store }

func New(deps Deps) api.Composer { return service{store: deps.Store} }
