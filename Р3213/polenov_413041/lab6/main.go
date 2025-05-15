package main

import (
	"bufio"
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"math"
	"os"
	"strconv"
	"strings"
)

const MAX_ITERS = 20

type ODEFunc func(x, y float64) float64
type ExactFunc func(x, x0, y0 float64) float64

func input(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	if text == "q" {
		fmt.Println("! Выход из программы.")
		os.Exit(0)
	}
	return text
}

func parseFloat(prompt string) float64 {
	for {
		if f, err := strconv.ParseFloat(prompt, 64); err == nil {
			return f
		}
		fmt.Println("! Некорректный ввод. Попробуйте еще раз")
	}
}

func parseInt(prompt string) int {
	for {
		if i, err := strconv.Atoi(prompt); err == nil {
			return i
		}
		fmt.Println("! Некорректный ввод. Попробуйте еще раз")
	}
}

func selectODE() (ODEFunc, ExactFunc) {
	fmt.Println("ОДУ:")
	fmt.Println("1. y + (1 + x)*y^2")
	fmt.Println("2. x + y")
	fmt.Println("3. sin(x) - y")
	fmt.Println("4. y / x")
	fmt.Println("5. e^x\n")

	for {
		val := parseInt(input("> Выберите ОДУ [1/2/3/4/5]: "))
		switch val {
		case 1:
			f := func(x, y float64) float64 {
				return y + (1+x)*y*y
			}
			exact := func(x, x0, y0 float64) float64 {
				expX := math.Exp(x)
				expX0 := math.Exp(x0)
				num := -expX
				denom := x*expX - (x0*expX0*y0+expX0)/y0
				return num / denom
			}
			return f, exact
		case 2:
			return func(x, y float64) float64 {
					return x + y
				}, func(x, x0, y0 float64) float64 {
					return math.Exp(x-x0)*(y0+x0+1) - x - 1
				}
		case 3:
			return func(x, y float64) float64 {
					return math.Sin(x) - y
				}, func(x, x0, y0 float64) float64 {
					num := 2*math.Exp(x0)*y0 - math.Exp(x0)*math.Sin(x0) + math.Exp(x0)*math.Cos(x0)
					return num/(2*math.Exp(x)) + math.Sin(x)/2 - math.Cos(x)/2
				}
		case 4:
			return func(x, y float64) float64 {
					return y / x
				}, func(x, x0, y0 float64) float64 {
					return x * y0 / x0
				}
		case 5:
			return func(x, y float64) float64 {
					return math.Exp(x)
				}, func(x, x0, y0 float64) float64 {
					return y0 - math.Exp(x0) + math.Exp(x)
				}
		default:
			fmt.Println("! Некорректный ввод. Попробуйте еще раз")
		}
	}
}

func improvedEulerMethod(f ODEFunc, xs []float64, y0 float64) []float64 {
	ys := []float64{y0}
	h := xs[1] - xs[0]

	for i := 0; i < len(xs)-1; i++ {
		y := ys[i]
		x := xs[i]

		yPred := f(x, y)
		yCorr := f(x+h, y+h*yPred)

		yNext := y + 0.5*h*(yPred+yCorr)
		ys = append(ys, yNext)
	}

	return ys
}

func rungeKutta4Method(f ODEFunc, xs []float64, y0 float64) []float64 {
	ys := []float64{y0}
	h := xs[1] - xs[0]

	for i := 0; i < len(xs)-1; i++ {
		x := xs[i]
		y := ys[i]

		k1 := h * f(x, y)
		k2 := h * f(x+h/2, y+k1/2)
		k3 := h * f(x+h/2, y+k2/2)
		k4 := h * f(x+h, y+k3)

		yNext := y + (k1+2*k2+2*k3+k4)/6
		ys = append(ys, yNext)
	}

	return ys
}

func adamsMethod(f ODEFunc, xs []float64, y0 float64, eps float64) []float64 {
	n := len(xs)
	h := xs[1] - xs[0]
	ys := rungeKutta4Method(f, xs[:4], y0)

	for i := 3; i < n-1; i++ {
		// Предиктор (Адамс-Бэшфорт 4-го порядка)
		yp := ys[i] + h/24*(55*f(xs[i], ys[i])-59*f(xs[i-1], ys[i-1])+37*f(xs[i-2], ys[i-2])-9*f(xs[i-3], ys[i-3]))

		// Корректор (Адамс-Мултон 4-го порядка)
		yNext := yp
		for {
			yCorr := ys[i] + h/24*(9*f(xs[i+1], yNext)+19*f(xs[i], ys[i])-5*f(xs[i-1], ys[i-1])+f(xs[i-2], ys[i-2]))
			if math.Abs(yCorr-yNext) < eps {
				yNext = yCorr
				break
			}
			yNext = yCorr
		}

		ys = append(ys, yNext)
	}

	return ys
}

func makeUniformGrid(x0, xn float64, n int) []float64 {
	xs := make([]float64, n)
	h := (xn - x0) / float64(n)
	for i := 0; i < n; i++ {
		xs[i] = x0 + float64(i)*h
	}
	return xs
}

func plotGraph(xs, ysApprox, ysExact []float64, methodName string) {
	seriesApprox := chart.ContinuousSeries{
		Name:    "Приближенное решение (точки)",
		XValues: xs,
		YValues: ysApprox,
		Style: chart.Style{
			DotWidth:    4,
			DotColor:    chart.ColorRed,
			StrokeWidth: 0, // Без линии
		},
	}

	seriesExact := chart.ContinuousSeries{
		Name:    "Точное решение",
		XValues: xs,
		YValues: ysExact,
		Style: chart.Style{
			StrokeColor: chart.ColorBlue,
			StrokeWidth: 2,
		},
	}

	graph := chart.Chart{
		Title:  methodName,
		XAxis:  chart.XAxis{Name: "X"},
		YAxis:  chart.YAxis{Name: "Y"},
		Series: []chart.Series{seriesExact, seriesApprox},
	}

	file, err := os.Create(methodName + ".png")
	if err != nil {
		fmt.Printf("Не удалось создать файл графика: %v\n", err)
		return
	}
	defer file.Close()

	err = graph.Render(chart.PNG, file)
	if err != nil {
		fmt.Printf("Ошибка при сохранении графика: %v\n", err)
		return
	}
	fmt.Printf("График сохранен в файл %s.png\n", methodName)
}

func solve(f ODEFunc, x0, xn float64, n int, y0 float64, exactY func(float64, float64, float64) float64, eps float64) {
	methods := []struct {
		name string
		fn   func(ODEFunc, []float64, float64) []float64
	}{
		{"Усовершенствованный Эйлер", improvedEulerMethod},
		{"Рунге-Кутта 4", rungeKutta4Method},
		{"Адамс", func(f ODEFunc, xs []float64, y0 float64) []float64 {
			return adamsMethod(f, xs, y0, eps)
		}},
	}

	for _, method := range methods {
		fmt.Println(method.name + ":\n")
		ni := n
		iters := 0
		xs := makeUniformGrid(x0, xn, ni)
		ys := method.fn(f, xs, y0)
		inaccuracy := math.Inf(1)

		for inaccuracy > eps {
			if iters >= MAX_ITERS {
				fmt.Printf("! Не удалось достичь точности за %d итераций.\n\n", iters)
				break
			}
			iters++
			ni *= 2
			xs = makeUniformGrid(x0, xn, ni)
			newYs := method.fn(f, xs, y0)

			var p int
			switch method.name {
			case "Усовершенствованный Эйлер":
				p = 2
			case "Рунге-Кутта 4":
				p = 4
			default:
				p = -1
			}

			if p > 0 {
				coef := math.Pow(2, float64(p)) - 1
				inaccuracy = math.Abs(newYs[len(newYs)-1]-ys[len(ys)-1]) / coef
			} else {
				inaccuracy = 0
				for i := 0; i < len(xs); i++ {
					err := math.Abs(exactY(xs[i], x0, y0) - newYs[i])
					if err > inaccuracy {
						inaccuracy = err
					}
				}
			}
			ys = newYs
		}

		h := (xn - x0) / float64(ni)
		fmt.Printf("Для eps = %g использовано n = %d, шаг h = %.6f, итераций = %d\n\n", eps, ni, h, iters)

		if len(xs) <= 100 {
			fmt.Print("y:\t\t[")
			for _, v := range ys {
				fmt.Printf("%.5f ", v)
			}
			fmt.Println("]")

			fmt.Print("y_точн:\t[")
			for _, x := range xs {
				fmt.Printf("%.5f ", exactY(x, x0, y0))
			}
			fmt.Println("]")
		} else {
			fmt.Print("y_точн:\t[точек слишком много, отображение пропущено]\n")
		}

		if method.name == "Адамс" {
			fmt.Printf("\nПогрешность (max|y_iточн - y_i|): %g\n", inaccuracy)
		} else {
			fmt.Printf("\nПогрешность (по правилу Рунге): %g\n", inaccuracy)
		}

		ysExact := make([]float64, len(xs))
		for i, x := range xs {
			ysExact[i] = exactY(x, x0, y0)
		}
		plotGraph(xs, ys, ysExact, method.name)

		fmt.Println(strings.Repeat("-", 30))
	}
}

func main() {
	f, exactY := selectODE()
	var x0, xn, y0, eps float64
	var n int

	for {
		x0 = parseFloat(input("> Введите первый элемент интервала x0: "))
		xn = parseFloat(input("> Введите последний элемент интервала xn: "))
		n = parseInt(input("> Введите количество элементов в интервале n: "))

		y0 = parseFloat(input("> Введите y0: "))
		eps = parseFloat(input("> Введите точность eps: "))

		if xn <= x0 {
			fmt.Println("! xn должен быть больше x0. Введите еще раз.")
		} else if n <= 1 {
			fmt.Println("! Количество элементов n должно быть > 1. Введите еще раз.")
		} else {
			break
		}
	}

	solve(f, x0, xn, n, y0, exactY, eps)
}
