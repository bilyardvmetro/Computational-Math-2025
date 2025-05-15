package main

import (
	"math"
	"sort"
)

// BuildFiniteDifferenceTable строит таблицу конечных разностей Δ^k y_i
// table[i][k] будет Δ^k y_i (k=0 для y_i, k=1 для Δy_i, и т.д.)
// Возвращает саму таблицу и шаг h.
func BuildFiniteDifferenceTable(points Points) ([][]float64, float64, error) {
	n := len(points)
	if n == 0 {
		return nil, 0, NewInterpolationError("нет точек для построения таблицы конечных разностей")
	}

	// Сортируем точки по X на всякий случай, хотя для равноотстоящих это должно быть уже так
	sortedPoints := make(Points, len(points))
	copy(sortedPoints, points)
	sort.Sort(sortedPoints) // Убедимся, что точки отсортированы по X

	isEquallySpaced, h := ArePointsEquallySpaced(sortedPoints)
	if !isEquallySpaced {
		return nil, 0, NewInterpolationError("узлы не являются равноотстоящими, метод Ньютона с конечными разностями неприменим")
	}
	if h <= floatTolerance && n > 1 { // Если шаг нулевой или очень маленький при наличии нескольких точек
		return nil, 0, NewInterpolationError("шаг h между узлами слишком мал или равен нулю, что делает метод некорректным")
	}

	// diffTable[i][k] будет хранить k-ую разность для y_i (Δ^k y_i)
	// diffTable[i][0] = y_i
	// diffTable[i][1] = Δy_i
	// diffTable[i][2] = Δ^2y_i
	diffTable := make([][]float64, n)
	for i := 0; i < n; i++ {
		diffTable[i] = make([]float64, n-i) // Максимальный порядок разности для y_i это n-1-i
		for j := range diffTable[i] {
			diffTable[i][j] = math.NaN() // Инициализируем NaN
		}
		diffTable[i][0] = sortedPoints[i].Y // Нулевая разность
	}

	for k := 1; k < n; k++ { // Порядок разности
		for i := 0; i < n-k; i++ { // Индекс y для которого вычисляется разность
			// Δ^k y_i = Δ^{k-1} y_{i+1} - Δ^{k-1} y_i
			diffTable[i][k] = diffTable[i+1][k-1] - diffTable[i][k-1]
		}
	}

	return diffTable, h, nil
}

// NewtonFiniteDifferenceInterpolate вычисляет значение интерполяционного многочлена Ньютона
// (первая формула, для равноотстоящих узлов) в точке xEval.
// diffTable[0][k] содержит Δ^k y_0
func NewtonFiniteDifferenceInterpolate(points Points, diffTable [][]float64, h float64, xEval float64) (float64, error) {
	n := len(points)
	if n == 0 {
		return math.NaN(), NewInterpolationError("нет точек для интерполяции")
	}
	if h == 0 && n > 1 { // h может быть 0 если всего одна точка, что нормально.
		return math.NaN(), NewInterpolationError("шаг h не может быть равен нулю для нескольких точек")
	}

	// Используем первую интерполяционную формулу Ньютона (для прямого хода)
	// P(x) = y0 + s*Δy0 + s(s-1)/2! * Δ^2y0 + s(s-1)(s-2)/3! * Δ^3y0 + ...
	// где s = (x - x0) / h

	x0 := points[0].X // Предполагается, что points отсортированы
	y0 := points[0].Y

	if n == 1 { // Если всего одна точка, то значение функции в ней и есть интерполяция
		return y0, nil
	}

	s := (xEval - x0) / h
	result := y0

	termMultiplier := 1.0    // s, s(s-1), s(s-1)(s-2), ...
	for k := 1; k < n; k++ { // k - порядок разности
		if k > len(diffTable[0])-1 || math.IsNaN(diffTable[0][k]) {
			// Если разность высокого порядка не существует (например, для линейной функции Δ^2y = 0)
			// или если в таблице закончились вычисленные значения для y0
			break
		}
		termMultiplier *= (s - float64(k-1)) // s для k=1; (s-1) для k=2 (итого s(s-1)); (s-2) для k=3 (итого s(s-1)(s-2))

		term := termMultiplier / Factorial(k) * diffTable[0][k] // diffTable[0][k] это Δ^k y_0
		result += term
	}

	return result, nil
}

// NewtonFiniteDifferencePolynomialFunction возвращает функцию, представляющую многочлен Ньютона (конечные разности)
func NewtonFiniteDifferencePolynomialFunction(points Points, diffTable [][]float64, h float64) (func(x float64) float64, error) {
	if len(points) == 0 {
		return nil, NewInterpolationError("нет точек для построения многочлена Ньютона (конечные разности)")
	}
	if h == 0 && len(points) > 1 {
		return nil, NewInterpolationError("шаг h не может быть равен нулю для нескольких точек")
	}
	// Проверка на то, что diffTable не пуста и содержит хотя бы y0
	if len(diffTable) == 0 || len(diffTable[0]) == 0 {
		return nil, NewInterpolationError("таблица конечных разностей пуста или некорректна")
	}

	return func(x float64) float64 {
		val, _ := NewtonFiniteDifferenceInterpolate(points, diffTable, h, x)
		return val
	}, nil
}

// InterpolationError пользовательский тип ошибки для интерполяции
type InterpolationError struct {
	message string
}

func NewInterpolationError(message string) *InterpolationError {
	return &InterpolationError{message: message}
}

func (e *InterpolationError) Error() string {
	return e.message
}
