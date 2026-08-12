// Package app composes the reference application's private cells.
package app

import (
	"net/http"

	greetingcompose "example.com/agent-first-reference/internal/cells/greeting-compose"
	greetingrender "example.com/agent-first-reference/internal/cells/greeting-render"
	tokenissue "example.com/agent-first-reference/internal/cells/token-issue"
	tokenrender "example.com/agent-first-reference/internal/cells/token-render"
)

// NewHandler assembles the reference application's HTTP routes.
func NewHandler() http.Handler {
	composer := greetingcompose.New(greetingcompose.Deps{Store: greetingcompose.NewStore("Hello")})
	greetingHandler := greetingrender.NewHandler(greetingrender.Deps{Composer: composer})
	issuer := tokenissue.New(tokenissue.Deps{Store: tokenissue.NewStore()})
	tokenHandler := tokenrender.NewHandler(tokenrender.Deps{Issuer: issuer})
	mux := http.NewServeMux()
	mux.Handle("GET /greeting/{name}", greetingHandler)
	mux.Handle("GET /token/{subject}", tokenHandler)
	return mux
}
