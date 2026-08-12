package modulus

import "example.com/research-library/internal/cells/modulus-reduce"

const Prime = 97

// Reduce returns the canonical representative of value modulo Prime.
func Reduce(value int) int {
	return modulusreduce.Reduce(value, Prime)
}
