package search

import "math"

// vector holds an n-dimensional vector with with value for each dimension.
type vector struct {
	label string
	val   map[string]float64
}

// mod calculates the modulus (magnitude) of a vector.
func (v *vector) mod() float64 {
	sum := 0.0
	for _, val := range v.val {
		sum += (val * val)
	}
	return math.Sqrt(sum)
}

// theta returns the angle (in radians) between two vectors.
func theta(a vector, b vector) float64 {
	cosTheta := dotProduct(a, b) / (a.mod() * b.mod())
	cosTheta = clamp(cosTheta, 0, 1) // to stop floating point errors
	x := math.Acos(cosTheta)
	return x
}

// dotProduct calculates the dot product of two vectors. It is assumed that
// vector.vals contain the same keys.
func dotProduct(a vector, b vector) float64 {
	sum := 0.0
	for word := range a.val {
		sum += a.val[word] * b.val[word]
	}
	return sum
}

// clamp constaints a value between the minimum and maximum.
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
