package main

import (
	"fmt"
	"math"
	"os"
)

func printTable(n int, X, Y []float64, phi func(float64) float64, outputFile *os.File) {
	fmt.Fprintf(outputFile, "%10s|", "X")
	for _, x := range X {
		fmt.Fprintf(outputFile, "%12.4f|", x)
	}
	fmt.Fprintln(outputFile)
	fmt.Fprintln(outputFile, stringRepeat("-", (n+1)*13))

	fmt.Fprintf(outputFile, "%10s|", "Y")
	for _, y := range Y {
		fmt.Fprintf(outputFile, "%12.4f|", y)
	}
	fmt.Fprintln(outputFile)
	fmt.Fprintln(outputFile, stringRepeat("-", (n+1)*13))

	fmt.Fprintf(outputFile, "%10s|", "phi(X)")
	for _, x := range X {
		fmt.Fprintf(outputFile, "%12.4f|", phi(x))
	}
	fmt.Fprintln(outputFile)
	fmt.Fprintln(outputFile, stringRepeat("-", (n+1)*13))

	e_i := make([]float64, n)
	for i := 0; i < n; i++ {
		e_i[i] = phi(X[i]) - Y[i]
	}
	fmt.Fprintf(outputFile, "%10s|", "e_i")
	for _, e := range e_i {
		fmt.Fprintf(outputFile, "%12.4f|", e)
	}
	fmt.Fprintln(outputFile)

	e_i2 := make([]float64, n)
	for i, e := range e_i {
		e_i2[i] = e * e
	}
	sum := 0.0
	for _, e2 := range e_i2 {
		sum += e2
	}
	rmse := math.Sqrt(sum / float64(n))
	fmt.Fprintf(outputFile, "СРЕДНЕКВАДРАТИЧНОЕ ОТКЛОНЕНИЕ: %.4f\n\n", rmse)
}

func stringRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
