package main

import (
	"fmt"
	"math"
)

// BuildDividedDifferenceTable строит таблицу разделенных разностей
// Возвращает таблицу, где первый столбец - это f[xi], f[xi, xi+1], f[xi, xi+1, xi+2], etc.
// Точнее, возвращает коэффициенты f[x0], f[x0,x1], f[x0,x1,x2], ...
// Которые являются верхним диагональным рядом полной таблицы разделенных разностей.
//
// Таблица выглядит так (coeffs[i] = table[0][i]):
// y0
// y1   f[x0,x1]
// y2   f[x1,x2]   f[x0,x1,x2]
// y3   f[x2,x3]   f[x1,x2,x3]   f[x0,x1,x2,x3]
//
// Возвращаем массив coeffs = [y0, f[x0,x1], f[x0,x1,x2], f[x0,x1,x2,x3]]
// И полную таблицу для вывода пользователю.
func BuildDividedDifferenceTable(points Points) ([][]float64, []float64, error) {
	n := len(points)
	if n == 0 {
		return nil, nil, NewInterpolationError("нет точек для построения таблицы разделенных разностей")
	}

	// Проверка на уникальность X координат
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if math.Abs(points[i].X-points[j].X) < floatTolerance {
				return nil, nil, NewInterpolationError("дублирующиеся X-координаты узлов не допускаются для разделенных разностей")
			}
		}
	}

	// table[i][j] будет хранить j-ую разделенную разность, начинающуюся с x_i
	// table[i][0] = y_i
	// table[i][1] = f[x_i, x_{i+1}]
	// table[i][2] = f[x_i, x_{i+1}, x_{i+2}]
	// ...
	table := make([][]float64, n)
	for i := range table {
		table[i] = make([]float64, n-i)
		table[i][0] = points[i].Y // Нулевая разделенная разность f[x_i] = y_i
	}

	for j := 1; j < n; j++ { // Порядок разности (столбец)
		for i := 0; i < n-j; i++ { // Начальная точка для разности (строка)
			// f[x_i, ..., x_{i+j}] = (f[x_{i+1}, ..., x_{i+j}] - f[x_i, ..., x_{i+j-1}]) / (x_{i+j} - x_i)
			numerator := table[i+1][j-1] - table[i][j-1]
			denominator := points[i+j].X - points[i].X
			if math.Abs(denominator) < floatTolerance {
				return nil, nil, NewInterpolationError("деление на ноль при вычислении разделенных разностей (возможно, совпадающие X узлов)")
			}
			table[i][j] = numerator / denominator
		}
	}

	// Коэффициенты полинома Ньютона - это верхняя диагональ: table[0][0], table[0][1], ..., table[0][n-1]
	coeffs := make([]float64, n)
	for j := 0; j < n; j++ {
		coeffs[j] = table[0][j]
	}

	// divDiffDisplayTable[i] будет строкой для x_i, y_i, f[x_i,x_{i+1}], f[x_i,x_{i+1},x_{i+2}], ...
	divDiffDisplayTable := make([][]float64, n)
	for i := 0; i < n; i++ {
		divDiffDisplayTable[i] = make([]float64, n-i) // Первая разность - это y_i
		// divDiffDisplayTable[i][0] = points[i].Y // Уже в table[i][0]
		for j := 0; j < n-i; j++ {
			divDiffDisplayTable[i][j] = table[i][j]
		}
	}

	return divDiffDisplayTable, coeffs, nil
}

// NewtonDividedDifferenceInterpolate вычисляет значение интерполяционного многочлена Ньютона
// с использованием разделенных разностей в точке xEval.
func NewtonDividedDifferenceInterpolate(points Points, coeffs []float64, xEval float64) (float64, error) {
	n := len(points)
	if n == 0 || len(coeffs) == 0 {
		return math.NaN(), NewInterpolationError("нет точек или коэффициентов для интерполяции Ньютона")
	}
	if len(coeffs) != n {
		return math.NaN(), NewInterpolationError(fmt.Sprintf("количество коэффициентов (%d) не совпадает с количеством точек (%d)", len(coeffs), n))
	}

	result := coeffs[0]
	termProduct := 1.0 // (x - x0), (x - x0)(x - x1), ...

	for k := 1; k < n; k++ {
		termProduct *= (xEval - points[k-1].X)
		result += coeffs[k] * termProduct
	}
	return result, nil
}

// NewtonDividedDifferencePolynomialFunction возвращает функцию, представляющую многочлен Ньютона (разделенные разности)
func NewtonDividedDifferencePolynomialFunction(points Points, coeffs []float64) (func(x float64) float64, error) {
	if len(points) == 0 || len(coeffs) == 0 {
		return nil, NewInterpolationError("нет точек или коэффициентов для построения многочлена Ньютона")
	}
	if len(coeffs) != len(points) {
		return nil, NewInterpolationError(fmt.Sprintf("количество коэффициентов (%d) не совпадает с количеством точек (%d)", len(coeffs), len(points)))
	}

	return func(x float64) float64 {
		val, _ := NewtonDividedDifferenceInterpolate(points, coeffs, x) // Ошибки должны быть обработаны при получении coeffs
		return val
	}, nil
}
