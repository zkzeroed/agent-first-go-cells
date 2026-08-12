// Package api defines the public contract for the text-fold cell.
package api

// Folder canonicalizes a label for display-independent comparison.
type Folder interface {
	Fold(value string) string
}
