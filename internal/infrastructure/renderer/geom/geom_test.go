package geom

import "testing"

func TestAbs(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-100, 100},
	}

	for _, tt := range tests {
		got := Abs(tt.input)
		if got != tt.want {
			t.Errorf("Abs(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		input float64
		want  float64
		delta float64
	}{
		{4.0, 2.0, 0.001},
		{9.0, 3.0, 0.001},
		{0.0, 0.0, 0.001},
		{-1.0, 0.0, 0.001},
		{2.0, 1.414, 0.01},
	}

	for _, tt := range tests {
		got := Sqrt(tt.input)
		diff := got - tt.want
		if diff < 0 {
			diff = -diff
		}
		if diff > tt.delta {
			t.Errorf("Sqrt(%f) = %f, want %f (delta %f)", tt.input, got, tt.want, tt.delta)
		}
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{1, 2, 2},
		{5, 3, 5},
		{-1, -2, -1},
		{0, 0, 0},
	}

	for _, tt := range tests {
		got := MaxInt(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("MaxInt(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMinInt(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{-1, -2, -2},
		{0, 0, 0},
	}

	for _, tt := range tests {
		got := MinInt(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("MinInt(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
