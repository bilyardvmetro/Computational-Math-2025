package main

import (
	"bufio"
	"fmt"
	"math"
)

type function func(x float64, y float64) float64

type equationExtend struct {
	funcs       []function
	derivativeX []function
	derivativeY []function
	x           float64
	y           float64
	xPlace      []float64
	yPlace      []float64
	accuracy    float64
	vectorError []float64
}

// Получение первой системы уравнений
func getFirstSystem(eq *equationExtend) {
	eq.funcs = make([]function, 2)
	eq.funcs[0] = func(x float64, y float64) float64 {
		return 0.3 - 0.1*x*x - 0.2*y*y
	}
	eq.funcs[1] = func(x float64, y float64) float64 {
		return 0.7 - 0.2*x*x - 0.1*x*y
	}

	eq.derivativeX = make([]function, 2)
	eq.derivativeX[0] = func(x float64, y float64) float64 {
		return -0.2 * x
	}
	eq.derivativeX[1] = func(x float64, y float64) float64 {
		return -0.4*x - 0.1*y
	}

	eq.derivativeY = make([]function, 2)
	eq.derivativeY[0] = func(x float64, y float64) float64 {
		return -0.4 * y
	}
	eq.derivativeY[1] = func(x float64, y float64) float64 {
		return -0.1 * x
	}
}

// Получение второй системы уравнений
func getSecondSystem(eq *equationExtend) {
	eq.funcs = make([]function, 2)
	eq.funcs[0] = func(x float64, y float64) float64 {
		return 1 + math.Sin(y+1)
	}
	eq.funcs[1] = func(x float64, y float64) float64 {
		return 1.5 - math.Sin(x-1)
	}

	eq.derivativeX = make([]function, 2)
	eq.derivativeX[0] = func(x float64, y float64) float64 {
		return 0
	}
	eq.derivativeX[1] = func(x float64, y float64) float64 {
		return math.Cos(x - 1)
	}

	eq.derivativeY = make([]function, 2)
	eq.derivativeY[0] = func(x float64, y float64) float64 {
		return math.Cos(y + 1)
	}
	eq.derivativeY[1] = func(x float64, y float64) float64 {
		return 0
	}
}

// Получение области изоляции корня, начальное приближение и точность
func getJunkInfo(in *bufio.Reader, eq *equationExtend) {
	fmt.Print("Введите область изоляции корня (x_start x_end y_start y_end): ")
	var x1, x2, y1, y2 float64
	ReadFloat(in, &x1, false, "левой границы изоляции по x")
	ReadFloat(in, &x2, false, "правой границы изоляции по x")
	ReadFloat(in, &y1, false, "левой границы изоляции по y")
	ReadFloat(in, &y2, true, "правой границы изоляции по y")
	eq.xPlace = make([]float64, 2)
	eq.xPlace[0], eq.xPlace[1] = x1, x2
	eq.yPlace = make([]float64, 2)
	eq.yPlace[0], eq.yPlace[1] = y1, y2

	fmt.Print("Введите начальное приближение (x y): ")
	var x, y float64
	ReadFloat(in, &x, false, "root x")
	ReadFloat(in, &y, true, "root y")
	eq.x, eq.y = x, y

	fmt.Print("Введите точность: ")
	var epsilon float64
	ReadFloat(in, &epsilon, true, "точности")
	eq.accuracy = epsilon
}

// Проверка правильности использования метода простых итераций для первой системы уравнений
func checkSystem(eq equationExtend) bool {
	var x, y = MaxAbs(eq.xPlace[0], eq.xPlace[1]), MaxAbs(eq.yPlace[0], eq.yPlace[1])
	for i := 0; i < 2; i++ {
		if Abs(eq.derivativeX[i](x, y))+Abs(eq.derivativeY[i](x, y)) >= 1 {
			return false
		}
	}
	return true
}

// Решение системы уравнений методом простых итераций
func solveSystem(eq *equationExtend, M int) int {
	var maximum = -1.0
	var k = 0
	eq.vectorError = make([]float64, 2)
	for Abs(maximum) >= eq.accuracy && k < M {
		var x, y float64
		x = eq.funcs[0](eq.x, eq.y)
		y = eq.funcs[1](eq.x, eq.y)
		eq.vectorError[0] = Abs(x - eq.x)
		maximum = eq.vectorError[0]
		eq.vectorError[1] = Abs(y - eq.y)
		maximum = max(maximum, eq.vectorError[1])
		eq.x = x
		eq.y = y
		k += 1
	}
	return k
}

// LinearSystem Запуск программы по решению системы нелинейных уравнений
func LinearSystem(in *bufio.Reader, out *bufio.Writer) {
	var eq equationExtend
	fmt.Print("Какую систему вы хотите решить? (введите номер)\n 1) 0.1x^2 + 0.2y^2 + x - 0.3 = 0\n    0.2x^2 + 0.1xy + y - 0.7 = 0\n 2) sin(x - 1) + y = 1.5\n    x - sin(y + 1) = 1\n Enter: ")
	var option int
	ReadInt(in, &option, true)

	if option == 1 {
		getFirstSystem(&eq)
	} else {
		getSecondSystem(&eq)
	}

	getJunkInfo(in, &eq)
	if !checkSystem(eq) {
		GetOut(SystemError{})
	}

	err := DrawTwoFunctions(func(x float64, y float64) float64 { return x - eq.funcs[0](x, y) },
		func(x float64, y float64) float64 { return y - eq.funcs[1](x, y) }, eq.xPlace[0], eq.xPlace[1], eq.yPlace[0], eq.yPlace[1], eq.accuracy)

	if err != nil {
		GetOut(err)
	}

	var iterations = solveSystem(&eq, 1000000)

	fmt.Fprintln(out, "Вектор неизвестных: ", eq.x, eq.y)
	fmt.Fprintln(out, "Количество итераций: ", iterations)
	fmt.Fprintln(out, "Вектор погрешностей: ", eq.vectorError)
}
