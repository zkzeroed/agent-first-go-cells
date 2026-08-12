// Package api defines the public contract for the text-normalize cell.
package api

// Normalizer canonicalizes a display name for rendering.
type Normalizer interface {
	Normalize(string) string
}
