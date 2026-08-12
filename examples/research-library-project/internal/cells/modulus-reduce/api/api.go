// Package api defines the public contract for the modulus-reduce cell.
package api

// Reducer computes a canonical residue for a positive modulus.
type Reducer interface {
	Reduce(value, modulus int) int
}
