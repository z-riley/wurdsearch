package search

import "math"

// vector holds an n-dimensional vector with with value for each dimension
type vector struct {
	label string
	val   map[string]float64
}

// mod calculates the modulus (magnitude) of a vector
func (v *vector) mod() float64 {
	sum := 0.0
	for _, val := range v.val {
		sum += (val * val)
	}
	return math.Sqrt(sum)
}

// theta returns the angle (in radians) between two vectors
func theta(a vector, b vector) float64 {
	return math.Acos(dotProduct(a, b) / (a.mod() * b.mod()))
}

// dotProduct calculates the dot product of two vectors. It is assumed that vector.vals contain the same keys
func dotProduct(a vector, b vector) float64 {
	sum := 0.0
	for word := range a.val {
		sum += a.val[word] * b.val[word]
	}
	return sum
}
