package modulusreduce

// Reduce computes the Euclidean remainder of value modulo a positive modulus.
func Reduce(value, modulus int) int {
	remainder := value % modulus
	if remainder < 0 {
		remainder += modulus
	}
	return remainder
}
