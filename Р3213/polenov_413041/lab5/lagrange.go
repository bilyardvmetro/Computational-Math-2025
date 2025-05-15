package main

import "math"

// LagrangeInterpolate вычисляет значение интерполяционного многочлена Лагранжа в точке xEval
func LagrangeInterpolate(points Points, xEval float64) (float64, error) {
	if len(points) == 0 {
		return math.NaN(), NewInterpolationError("нет точек для интерполяции")
	}
	n := len(points)
	result := 0.0

	for j := 0; j < n; j++ {
		term := points[j].Y
		for i := 0; i < n; i++ {
			if i == j {
				continue
			}
			if math.Abs(points[j].X-points[i].X) < floatTolerance {
				if math.Abs(xEval-points[j].X) < floatTolerance {
					// если все X уникальны,
					// то это условие (points[j].X - points[i].X == 0) не должно срабатывать.
				}
				// Проверка на случай дублирования X-координат узлов
				return math.NaN(), NewInterpolationError(
					"дублирующиеся X-координаты узлов не допускаются для многочлена Лагранжа, или деление на ноль",
				)
			}
			term = term * (xEval - points[i].X) / (points[j].X - points[i].X)
		}
		result += term
	}
	return result, nil
}

// LagrangePolynomialFunction возвращает функцию, представляющую многочлен Лагранжа
func LagrangePolynomialFunction(points Points) (func(x float64) float64, error) {
	if len(points) == 0 {
		return nil, NewInterpolationError("нет точек для построения многочлена Лагранжа")
	}
	// Проверка на уникальность X координат
	uniqueXs := make(map[float64]bool)
	for _, p := range points {
		if uniqueXs[p.X] {
			return nil, NewInterpolationError("дублирующиеся X-координаты узлов не допускаются")
		}
		uniqueXs[p.X] = true
	}

	return func(x float64) float64 {
		val, _ := LagrangeInterpolate(points, x) // Ошибку уже проверили выше
		return val
	}, nil
}
