package textfold

import "strings"

// Fold returns a lowercase form with whitespace collapsed to single spaces.
func Fold(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}
