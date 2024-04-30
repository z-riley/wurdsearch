package search

import (
	"math"
	"testing"
)

func TestTheta(t *testing.T) {
	a := vector{
		label: "a",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
	b := vector{
		label: "b",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
	th := theta(a, b)
	expected := 0.0
	if th != expected {
		t.Errorf("Actual (%f) did not equal expected (%f)", th, expected)
	}

	a = vector{
		label: "a",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
	b = vector{
		label: "b",
		val: map[string]float64{
			"three": -3,
			"one":   -1,
			"two":   -2,
		},
	}
	th = theta(a, b)
	expected = math.Pi
	if th != expected {
		t.Errorf("Actual (%f) did not equal expected (%f)", th, expected)
	}
}

func TestDotProduct(t *testing.T) {
	a := vector{
		label: "a",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
			"four":  4,
		},
	}
	b := vector{
		label: "b",
		val: map[string]float64{
			"four":  4,
			"two":   2,
			"three": 3,
			"one":   1,
		},
	}

	dp := dotProduct(a, b)
	expected := 30.0
	if dp != expected {
		t.Error("Actual did not equal expected")
	}
}

func TestMod(t *testing.T) {
	v := vector{
		label: "a",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
			"four":  -4,
		},
	}
	mod := v.mod()
	expected := math.Sqrt(30.0)
	if mod != expected {
		t.Error("Actual did not equal expected")
	}

	v = vector{
		label: "a",
		val: map[string]float64{
			"one": 0,
			"two": 0,
		},
	}
	mod = v.mod()
	expected = 0.0
	if mod != expected {
		t.Error("Actual did not equal expected")
	}
}
