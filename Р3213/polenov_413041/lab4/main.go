package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	n   int
	X   []float64
	Y   []float64
	res = Results{
		Linear:    make([]float64, 2),
		Quadratic: make([]float64, 3),
		Cubic:     make([]float64, 4),
		Exp:       make([]interface{}, 4),
		Log:       make([]interface{}, 4),
		Pow:       make([]interface{}, 4),
	}
)

func validateData(n int, xLine, yLine string) (int, []float64, []float64, error) {
	if n < 8 || n > 12 {
		return 0, nil, nil, fmt.Errorf("n out of range")
	}

	xStr := strings.Fields(xLine)
	yStr := strings.Fields(yLine)
	if len(xStr) != n || len(yStr) != n {
		return 0, nil, nil, fmt.Errorf("invalid data count")
	}

	X, Y := make([]float64, n), make([]float64, n)
	for i := 0; i < n; i++ {
		var err error
		X[i], err = strconv.ParseFloat(strings.ReplaceAll(xStr[i], ",", "."), 64)
		if err != nil {
			return 0, nil, nil, err
		}
		Y[i], err = strconv.ParseFloat(strings.ReplaceAll(yStr[i], ",", "."), 64)
		if err != nil {
			return 0, nil, nil, err
		}
	}

	return n, X, Y, nil
}

//func readData() (int, []float64, []float64) {
//	scanner := bufio.NewScanner(os.Stdin)
//	for {
//		fmt.Print("введите f если хотите прочитать из файла и любой другой символ в ином случае:")
//		scanner.Scan()
//		key := scanner.Text()
//		var xLine, yLine string
//		var n int
//		if strings.TrimSpace(key) == "f" {
//			fmt.Print("введите название файла:")
//			scanner.Scan()
//			filename := scanner.Text()
//			file, err := os.Open(filename)
//			if err != nil {
//				fmt.Println("ошибка открытия файла:", err)
//				continue
//			}
//			defer file.Close()
//			fileScanner := bufio.NewScanner(file)
//			fileScanner.Scan()
//			n, _ = strconv.Atoi(fileScanner.Text())
//			fileScanner.Scan()
//			xLine = strings.ReplaceAll(fileScanner.Text(), ",", ".")
//			fileScanner.Scan()
//			yLine = strings.ReplaceAll(fileScanner.Text(), ",", ".")
//		} else {
//			fmt.Print("введите количество узлов функции:")
//			scanner.Scan()
//			n, _ = strconv.Atoi(scanner.Text())
//			fmt.Print("введите через пробел значения x_i:")
//			scanner.Scan()
//			xLine = scanner.Text()
//			fmt.Print("введите через пробел значения y_i:")
//			scanner.Scan()
//			yLine = scanner.Text()
//		}
//
//		xStr := strings.Fields(xLine)
//		yStr := strings.Fields(yLine)
//		if len(xStr) != n || len(yStr) != n {
//			fmt.Println("error: некорректный ввод :( попробуйте еще раз!")
//			continue
//		}
//		X := make([]float64, n)
//		Y := make([]float64, n)
//		for i := 0; i < n; i++ {
//			X[i], _ = strconv.ParseFloat(xStr[i], 64)
//			Y[i], _ = strconv.ParseFloat(yStr[i], 64)
//		}
//		return n, X, Y
//	}
//}

func readData() (int, []float64, []float64) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("введите f если хотите прочитать из файла и любой другой символ в ином случае:")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)
		var n int
		var xLine, yLine string
		if key == "f" {
			fmt.Print("введите название файла:")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)
			file, err := os.Open(filename)
			if err != nil {
				fmt.Println("ошибка открытия файла:", err)
				continue
			}
			scanner := bufio.NewScanner(file)
			if scanner.Scan() {
				n, _ = strconv.Atoi(scanner.Text())
			}
			if scanner.Scan() {
				xLine = scanner.Text()
			}
			if scanner.Scan() {
				yLine = scanner.Text()
			}
			file.Close()
		} else {
			fmt.Print("введите количество узлов функции:")
			fmt.Scanln(&n)
			fmt.Print("введите через пробел значения x_i:")
			xLine, _ = reader.ReadString('\n')
			fmt.Print("введите через пробел значения y_i:")
			yLine, _ = reader.ReadString('\n')
		}

		n, X, Y, err := validateData(n, xLine, yLine)
		if err != nil {
			fmt.Println("error: некорректный ввод :( попробуйте еще раз!")
			continue
		}
		return n, X, Y
	}
}

func printData() {
	fmt.Print("введите f, если хотите записать результаты в файл и любой другой символ, если в консоль):")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	key := scanner.Text()
	var out = os.Stdout
	if strings.TrimSpace(key) == "f" {
		f, err := os.Create("output.txt")
		if err != nil {
			fmt.Println("не удалось создать файл вывода")
			return
		}
		defer f.Close()
		out = f
	}

	writeTitleAndTable := func(title string, x, y []float64, f func(float64) float64) {
		fmt.Fprintln(out, strings.Repeat(" ", 20)+title)
		printTable(len(x), x, y, f, out)
	}

	writeTitleAndTable("МЕТОД ЛИНЕЙНОЙ АППРОКСИМАЦИИ:", X, Y, func(x float64) float64 {
		return res.Linear[0]*x + res.Linear[1]
	})

	writeTitleAndTable("МЕТОД КВАДРАТИЧНОЙ АППРОКСИМАЦИИ:", X, Y, func(x float64) float64 {
		return res.Quadratic[0]*x*x + res.Quadratic[1]*x + res.Quadratic[2]
	})

	writeTitleAndTable("МЕТОД КУБИЧЕСКОЙ АППРОКСИМАЦИИ:", X, Y, func(x float64) float64 {
		return res.Cubic[0]*x*x*x + res.Cubic[1]*x*x + res.Cubic[2]*x + res.Cubic[3]
	})

	if res.Exp[2].([]float64) != nil && res.Exp[3].([]float64) != nil {
		xVals := res.Exp[2].([]float64)
		yVals := res.Exp[3].([]float64)
		writeTitleAndTable("МЕТОД ЭКСПОНЕНЦИАЛЬНОЙ АППРОКСИМАЦИИ:", xVals, yVals, func(x float64) float64 {
			return res.Exp[0].(float64) * math.Exp(res.Exp[1].(float64)*x)
		})
	}

	if res.Log[2].([]float64) != nil && res.Log[3].([]float64) != nil {
		xVals := res.Log[2].([]float64)
		yVals := res.Log[3].([]float64)
		writeTitleAndTable("МЕТОД ЛОГАРИФМИЧЕСКОЙ АППРОКСИМАЦИИ:", xVals, yVals, func(x float64) float64 {
			if x <= 0 {
				return math.NaN()
			}
			return res.Log[0].(float64)*math.Log(x) + res.Log[1].(float64)
		})
	}

	if res.Pow[2].([]float64) != nil && res.Pow[3].([]float64) != nil {
		xVals := res.Pow[2].([]float64)
		yVals := res.Pow[3].([]float64)
		writeTitleAndTable("МЕТОД СТЕПЕННОЙ АППРОКСИМАЦИИ:", xVals, yVals, func(x float64) float64 {
			if x <= 0 {
				return math.NaN()
			}
			return res.Pow[0].(float64) * math.Pow(x, res.Pow[1].(float64))
		})
	}
}

func main() {
	n, X, Y = readData()

	a, b := LinearApproximation(n, X, Y)
	res.Linear[0], res.Linear[1] = a, b
	DrawLinearGraph(X, Y, a, b, "graphics/linear.png")

	a2, b2, c2 := QuadraticApproximation(n, X, Y)
	res.Quadratic[0], res.Quadratic[1], res.Quadratic[2] = a2, b2, c2
	DrawQuadraticGraph(X, Y, a2, b2, c2, "graphics/quadratic.png")

	res.Cubic, _ = CubicApproximation(n, X, Y)
	DrawCubicGraph(X, Y, res.Cubic[0], res.Cubic[1], res.Cubic[2], res.Cubic[3], "graphics/cubic.png")

	if a, b, xAccepted, yAccepted := ExponentialApproximation(X, Y); xAccepted != nil && yAccepted != nil {
		res.Exp = []interface{}{a, b, xAccepted, yAccepted}
		DrawExponentialGraph(X, Y, xAccepted, a, b, "graphics/exp.png")
	} else {
		fmt.Println("error: недостаточно точек с положительными значениями функции для экспоненциальной аппроксимации")
	}

	if a, b, xAccepted, yAccepted := LogarithmApproximation(X, Y); xAccepted != nil && yAccepted != nil {
		res.Log = []interface{}{a, b, xAccepted, yAccepted}
		DrawLogGraph(X, Y, xAccepted, a, b, "graphics/log.png")
	} else {
		fmt.Println("error: недостаточно точек с положительными значениями функции для логарифмической аппроксимации")
	}

	if a, b, xAccepted, yAccepted := PowerApproximation(X, Y); xAccepted != nil && yAccepted != nil {
		res.Pow = []interface{}{a, b, xAccepted, yAccepted}
		DrawPowerGraph(X, Y, xAccepted, a, b, "graphics/pow.png")
	} else {
		fmt.Println("error: недостаточно точек с положительными значениями функции для степенной аппроксимации")
	}

	printData()
	drawAllApproximations(X, Y, res)
}
