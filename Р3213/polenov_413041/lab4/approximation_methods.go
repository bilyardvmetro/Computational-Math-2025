package main

import (
	"fmt"
	"log"
	"math"
)

func countSums1(n int, X, Y []float64) (float64, float64, float64, float64) {
	var S_X, S_Y, S_XX, S_XY float64
	for i := 0; i < n; i++ {
		S_X += X[i]
		S_Y += Y[i]
		S_XX += X[i] * X[i]
		S_XY += X[i] * Y[i]
	}
	return S_X, S_Y, S_XX, S_XY
}

func countSums2(n int, X, Y []float64) (float64, float64, float64) {
	var S_X3, S_X2Y, S_X4 float64
	for i := 0; i < n; i++ {
		S_X3 += math.Pow(X[i], 3)
		S_X2Y += math.Pow(X[i], 2) * Y[i]
		S_X4 += math.Pow(X[i], 4)
	}
	return S_X3, S_X2Y, S_X4
}

func countSums3(n int, X, Y []float64) (float64, float64, float64) {
	var S_X5, S_X6, S_X3Y float64
	for i := 0; i < n; i++ {
		S_X5 += math.Pow(X[i], 5)
		S_X6 += math.Pow(X[i], 6)
		S_X3Y += math.Pow(X[i], 3) * Y[i]
	}
	return S_X5, S_X6, S_X3Y
}

func LinearApproximation(n int, X, Y []float64) (float64, float64) {
	S_X, S_Y, S_XX, S_XY := countSums1(n, X, Y)

	delta := S_XX*float64(n) - S_X*S_X
	delta1 := S_XY*float64(n) - S_X*S_Y
	delta2 := S_XX*S_Y - S_X*S_XY

	return delta1 / delta, delta2 / delta
}

func solveLinearSystem(A [][]float64, B []float64) ([]float64, error) {
	n := len(B)
	X := make([]float64, n)
	for i := range A {
		for j := range A[i] {
			if math.IsNaN(A[i][j]) || math.IsInf(A[i][j], 0) {
				return nil, fmt.Errorf("invalid matrix value")
			}
		}
	}
	// Gaussian elimination
	for i := 0; i < n; i++ {
		pivot := A[i][i]
		if pivot == 0 {
			return nil, fmt.Errorf("singular matrix")
		}
		for j := i; j < n; j++ {
			A[i][j] /= pivot
		}
		B[i] /= pivot

		for k := 0; k < n; k++ {
			if k != i {
				factor := A[k][i]
				for j := i; j < n; j++ {
					A[k][j] -= factor * A[i][j]
				}
				B[k] -= factor * B[i]
			}
		}
	}
	copy(X, B)
	return X, nil
}

func QuadraticApproximation(n int, X, Y []float64) (float64, float64, float64) {
	S_X, S_Y, S_X2, S_XY := countSums1(n, X, Y)
	S_X3, S_X2Y, S_X4 := countSums2(n, X, Y)

	A := [][]float64{
		{float64(n), S_X, S_X2},
		{S_X, S_X2, S_X3},
		{S_X2, S_X3, S_X4},
	}
	B := []float64{S_Y, S_XY, S_X2Y}

	result, err := solveLinearSystem(A, B)
	if err != nil {
		log.Fatal(err)
	}
	return result[2], result[1], result[0]
}

func CubicApproximation(n int, X, Y []float64) ([]float64, error) {
	S_X, S_Y, S_X2, S_XY := countSums1(n, X, Y)
	S_X3, S_X2Y, S_X4 := countSums2(n, X, Y)
	S_X5, S_X6, S_X3Y := countSums3(n, X, Y)

	A := [][]float64{
		{S_X6, S_X5, S_X4, S_X3},
		{S_X5, S_X4, S_X3, S_X2},
		{S_X4, S_X3, S_X2, S_X},
		{S_X3, S_X2, S_X, float64(n)},
	}
	B := []float64{S_X3Y, S_X2Y, S_XY, S_Y}

	return solveLinearSystem(A, B)
}

func ExponentialApproximation(X, Y []float64) (float64, float64, []float64, []float64) {
	var cleanX []float64
	var cleanY []float64

	for i := range X {
		if Y[i] > 0 {
			cleanX = append(cleanX, X[i])
			cleanY = append(cleanY, Y[i])
		}
	}

	if len(cleanY) < 2 {
		return 0, 0, nil, nil
	}

	lnY := make([]float64, len(cleanY))
	for i, y := range cleanY {
		lnY[i] = math.Log(y)
	}
	B, A := LinearApproximation(len(lnY), cleanX, lnY)
	return math.Exp(A), B, cleanX, cleanY
}

func LogarithmApproximation(X, Y []float64) (float64, float64, []float64, []float64) {
	var cleanX []float64
	var cleanY []float64

	for i := range X {
		if X[i] > 0 {
			cleanX = append(cleanX, X[i])
			cleanY = append(cleanY, Y[i])
		}
	}

	if len(cleanX) < 2 {
		return 0, 0, nil, nil
	}

	lnX := make([]float64, len(cleanX))
	for i, x := range cleanX {
		lnX[i] = math.Log(x)
	}

	a, b := LinearApproximation(len(lnX), lnX, cleanY)
	return a, b, cleanX, cleanY
}

func PowerApproximation(X, Y []float64) (float64, float64, []float64, []float64) {
	var cleanX []float64
	var cleanY []float64

	for i := range X {
		if X[i] > 0 && Y[i] > 0 {
			cleanX = append(cleanX, X[i])
			cleanY = append(cleanY, Y[i])
		}
	}

	if len(cleanX) < 2 {
		return 0, 0, nil, nil
	}

	lnX := make([]float64, len(cleanX))
	lnY := make([]float64, len(cleanY))
	for i := range cleanX {
		lnX[i] = math.Log(cleanX[i])
		lnY[i] = math.Log(cleanY[i])
	}

	b, lnA := LinearApproximation(len(lnX), lnX, lnY)
	a := math.Exp(lnA)

	return a, b, cleanX, cleanY
}
