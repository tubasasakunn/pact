// Package geom provides shared geometry and math utilities for renderers.
package geom

// MaxInt returns the larger of two integers.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt returns the smaller of two integers.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Abs returns the absolute value of an integer.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Sqrt computes the square root using Newton's method.
func Sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}
