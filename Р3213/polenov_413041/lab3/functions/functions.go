package functions

import "math"

func F1(x float64) float64 {
	return x * x
}

func F2(x float64) float64 {
	return 1 / (math.Sqrt(x))
}

func F3(x float64) float64 {
	return 4*x*x - 2*x + 5
}

func F4(x float64) float64 {
	return 1 / x
}

func F5(x float64) float64 {
	return math.Sin(x)
}
