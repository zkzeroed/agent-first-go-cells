// Package api defines the public contract for the greeting-compose cell.
package api

import (
	"context"
	"errors"
)

var ErrEmptyName = errors.New("greeting-compose: empty name")

// Composer creates a greeting for a non-empty name.
type Composer interface {
	Compose(context.Context, string) (string, error)
}
