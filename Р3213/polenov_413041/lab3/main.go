package main

import (
	funcs "CompMathLab3/functions"
	meths "CompMathLab3/methods"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
)

var functions = []func(x float64) float64{funcs.F1, funcs.F2, funcs.F3, funcs.F4, funcs.F5}
var methods = map[string]func(f func(x float64) float64, a float64, b float64, n int) float64{
	"rec_left":  meths.SolveByRectangleLeft,
	"rec_right": meths.SolveByRectangleRight,
	"rec_mid":   meths.SolveByRectangleMid,
	"trapezia":  meths.SolveByTrapezia,
	"simpson":   meths.SolveBySimpson,
}
var rungeRatios = map[string]float64{
	"rec_left":  1,
	"rec_right": 1,
	"rec_mid":   3,
	"trapezia":  3,
	"simpson":   15,
}

func findBreakpoints(f func(x float64) float64, a float64, b float64, n int) []float64 {
	breakpoints := make([]float64, 0, n)
	h := (b - a) / float64(n)

	for i := 0; i <= n; i++ {
		x := a + float64(i)*h
		y := f(x)
		//fmt.Println(y)
		if y == math.Inf(1) || y == math.Inf(-1) {
			if !slices.Contains(breakpoints, x) {
				breakpoints = append(breakpoints, x)
			}
		}
	}

	return breakpoints
}

func tryCalc(f func(x float64) float64, x float64) (float64, error) {
	y := f(x)
	//fmt.Println(math.IsNaN(y))
	if y == math.Inf(1) || y == math.Inf(-1) || math.IsNaN(y) {
		return 0, errors.New("cannot calculate x")
	}
	return y, nil
}

func computeIntegral(f func(x float64) float64, a float64, b float64, eps float64, method string) (float64, int) {
	n := 4
	ratio := rungeRatios[method]
	result := methods[method](f, a, b, n)
	mismatch := math.Inf(1)

	for mismatch > eps {
		n *= 2
		secondResult := methods[method](f, a, b, n)
		mismatch = math.Abs(secondResult-result) / ratio

		result = secondResult
	}

	return result, n
}

func main() {
	fmt.Println("Выберите функцию для вычисления интеграла:")
	fmt.Println("1.\tx^2")
	fmt.Println("2.\t1 / √x")
	fmt.Println("3.\t4x^2 - 2x + 5")
	fmt.Println("4.\t1 / x")
	fmt.Println("5.\tsin(x)")
	fmt.Print("> ")

	var ans = 0
	_, err := fmt.Scanf("%d", &ans)
	if err != nil || (ans > 5 && ans < 1) {
		fmt.Println(fmt.Errorf("вы напечатали бредик"))
		fmt.Println("")
		os.Exit(1)
	}

	fmt.Println("Введите пределы интегрирования в формате: <верхний предел> <нижний предел>")
	fmt.Print("> ")

	var a, b float64
	_, err = fmt.Scanf("%f %f", &a, &b)
	if err != nil || a >= b {
		fmt.Println(fmt.Errorf("вы напечатали бредик"))
		fmt.Println("")
		os.Exit(1)
	}

	breakpoints := findBreakpoints(functions[ans-1], a, b, int(math.Ceil(b-a)*10))

	if len(breakpoints) != 0 {
		fmt.Printf("Найдены разрывы в точках: %v\n", breakpoints)
		epsilon := 0.0001
		converges := true

		for _, b := range breakpoints {
			y1, err1 := tryCalc(functions[ans-1], b-epsilon)
			y2, err2 := tryCalc(functions[ans-1], b+epsilon)

			if err1 == nil && err2 == nil && math.Abs(y1-y2) > epsilon || (y1 == y2 && err1 == nil) {
				converges = false
			}
		}

		if !converges {
			fmt.Println("Интеграл не сходится => решения не существует")
		} else {
			fmt.Println("Введите точность вычисления в формате: <точность>")
			fmt.Print("> ")

			var userEps float64
			_, err = fmt.Scanf("%f", &userEps)
			if err != nil {
				fmt.Println(fmt.Errorf("вы напечатали бредик"))
				fmt.Println("")
				os.Exit(1)
			}

			for method := range methods {
				fmt.Println("=============================================================================================")
				fmt.Printf("Вычисление методом %s...\n", method)

				if len(breakpoints) == 1 {
					if slices.Contains(breakpoints, a) {
						a += epsilon
					} else if slices.Contains(breakpoints, b) {
						b -= epsilon
					}

					result, iterations := computeIntegral(functions[ans-1], a, b, userEps, method)
					if result != 0.0 && iterations != 0 {
						fmt.Printf("Значение интеграла: %.4f | Число разбиений: %d\n", result, iterations)
					}
				}
				//		} else {
				//			var res = 0.0
				//			var n = 0
				//
				//			_, err3 := tryCalc(functions[ans-1], a)
				//			_, err4 := tryCalc(functions[ans-1], breakpoints[0]-eps)
				//			_, err5 := tryCalc(functions[ans-1], breakpoints[0]+eps)
				//			_, err6 := tryCalc(functions[ans-1], b)
				//
				//			if err3 == nil && err4 == nil {
				//				result, iterations := computeIntegral(functions[ans-1], a, breakpoints[0]-eps, eps, method)
				//				res += result
				//				n += iterations
				//			}
				//
				//			if err6 == nil && err5 == nil {
				//				result, iterations := computeIntegral(functions[ans-1], breakpoints[0]+eps, b, eps, method)
				//				res += result
				//				n += iterations
				//			}
				//
				//			for i := range breakpoints {
				//				bCur := breakpoints[i]
				//				bNext := breakpoints[i+1]
				//
				//				_, errCur := tryCalc(functions[ans-1], bCur+eps)
				//				_, errNext := tryCalc(functions[ans-1], bNext-eps)
				//
				//				if errCur == nil && errNext == nil {
				//					result, iterations := computeIntegral(functions[ans-1], bCur+eps, bNext-eps, eps, method)
				//					res += result
				//					n += iterations
				//				}
				//			}
				//
				//			fmt.Printf("Значение интеграла: %.4f | Число разбиений: %d\n", res, n)
				//		}
			}
		}
	} else {
		fmt.Println("Введите точность вычисления в формате: <точность>")
		fmt.Print("> ")

		var eps float64
		_, err = fmt.Scanf("%f", &eps)
		if err != nil {
			fmt.Println(fmt.Errorf("вы напечатали бредик"))
			fmt.Println("")
			os.Exit(1)
		}

		for method := range methods {
			fmt.Println("=============================================================================================")
			fmt.Printf("Вычисление методом %s...\n", method)
			res, iterations := computeIntegral(functions[ans-1], a, b, eps, method)
			fmt.Printf("Значение интеграла: %.4f | Число разбиений: %d\n", res, iterations)
		}
	}
}
