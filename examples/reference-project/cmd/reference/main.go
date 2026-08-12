package main

import (
	"log"
	"net/http"
	"time"

	"example.com/agent-first-reference/internal/app"
)

func main() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           app.NewHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
