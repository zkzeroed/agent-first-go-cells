// Package api defines the public contract for the token-render cell.
package api

import "net/http"

// Handler serves issued tokens over HTTP.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
