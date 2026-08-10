package main

import (
	"log"
	"net/http"
	"time"

	greetingcompose "example.com/agent-first-reference/internal/cells/greeting-compose"
	greetingrender "example.com/agent-first-reference/internal/cells/greeting-render"
	tokenissue "example.com/agent-first-reference/internal/cells/token-issue"
	tokenrender "example.com/agent-first-reference/internal/cells/token-render"
)

func main() {
	composer := greetingcompose.New(greetingcompose.Deps{Store: greetingcompose.NewStore("Hello")})
	handler := greetingrender.NewHandler(greetingrender.Deps{Composer: composer})
	issuer := tokenissue.New(tokenissue.Deps{Store: tokenissue.NewStore()})
	tokenHandler := tokenrender.NewHandler(tokenrender.Deps{Issuer: issuer})
	mux := http.NewServeMux()
	mux.Handle("GET /greeting/{name}", handler)
	mux.Handle("GET /token/{subject}", tokenHandler)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
