// Package api defines the public contract for the greeting-render cell.
package api

import "errors"

// ErrInvalidName reports a blank normalized name.
var ErrInvalidName = errors.New("greeting-render: invalid name")

// Renderer formats a greeting message.
type Renderer interface {
	Render(string) (string, error)
}
