package greetingrender

import (
	composeapi "example.com/agent-first-reference/internal/cells/greeting-compose/api"
	"example.com/agent-first-reference/internal/cells/greeting-render/api"
)

type Deps struct{ Composer composeapi.Composer }

func NewHandler(deps Deps) api.Handler { return handler{composer: deps.Composer} }
