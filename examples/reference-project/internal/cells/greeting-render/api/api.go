// Package api defines the public contract for the greeting-render cell.
package api

import "net/http"

// Handler serves rendered greetings.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
