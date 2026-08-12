package canonical

import "example.com/canonical/internal/cells/text-fold"

// Fold returns a lowercase, single-space form of a user-facing label.
func Fold(value string) string {
	return textfold.Fold(value)
}
