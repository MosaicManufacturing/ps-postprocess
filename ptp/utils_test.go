package ptp

import (
	"math"
	"testing"
)

func Test_normalizeAngle(t *testing.T) {
	testCases := []struct {
		name     string
		input    float64
		expected float64
	}{
		// Positive angles
		{"Zero", 0, 0},
		{"Pi", math.Pi, math.Pi},
		{"TwoPi", 2 * math.Pi, 0},
		{"JustOver2Pi", 2*math.Pi + 0.1, 0.1},
		{"3Pi", 3 * math.Pi, math.Pi},
		{"4Pi", 4 * math.Pi, 0},
		{"LargeAngle", 10 * math.Pi, 0},
		{"LargeAngle", 101 * math.Pi, math.Pi},

		// Negative angles
		{"NegativeSmall", -0.5, 2*math.Pi - 0.5},
		{"NegativePi", -math.Pi, math.Pi},
		{"Negative2Pi", -2 * math.Pi, 0},
		{"Negative3Pi", -3 * math.Pi, math.Pi},
		{"LargeNegative", -10 * math.Pi, 0},
		{"LargeNegative", -101 * math.Pi, math.Pi},

		// Very small angles
		{"VerySmallPositive", 0.000001, 0.000001},
		{"VerySmallNegative", -0.000001, 2*math.Pi - 0.000001},
	}

	const epsilon = 0.0001

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := normalizeAngle(tc.input)

			if math.Abs((result - tc.expected)) > epsilon {
				t.Errorf("normalizeAngle(%f) = %f, expected %f", tc.input, result, tc.expected)
			}
		})
	}
}
