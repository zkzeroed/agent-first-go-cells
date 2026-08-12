package textnormalize

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"example.com/agent-first-library/internal/cells/text-normalize/api"
)

type normalizer struct{}

// New constructs the private normalizer cell.
func New() api.Normalizer { return normalizer{} }

func (normalizer) Normalize(value string) string {
	words := strings.Fields(strings.ToLower(value))
	for i, word := range words {
		first, size := utf8.DecodeRuneInString(word)
		words[i] = string(unicode.ToUpper(first)) + word[size:]
	}
	return strings.Join(words, " ")
}
