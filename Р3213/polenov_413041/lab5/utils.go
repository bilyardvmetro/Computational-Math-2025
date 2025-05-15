package main

import (
	"fmt"
	"math"
	"sort"
)

const floatTolerance = 1e-9 // Допуск для сравнения float64

// ArePointsEquallySpaced проверяет, являются ли X-координаты точек равноотстоящими
func ArePointsEquallySpaced(points Points) (bool, float64) {
	if len(points) < 2 {
		return true, 0 // Недостаточно точек для определения шага
	}
	if len(points) == 2 {
		return true, points[1].X - points[0].X
	}

	// Сначала отсортируем точки по X, чтобы проверка была корректной
	sortedPoints := make(Points, len(points))
	copy(sortedPoints, points)
	sort.Sort(sortedPoints)

	h := sortedPoints[1].X - sortedPoints[0].X
	if h <= floatTolerance {
		// Если первые две точки очень близки или совпадают по X, то шаг почти нулевой.
		// В контексте интерполяции, совпадающие X-координаты узлов - это ошибка.
		// Здесь мы проверяем только на равноотстояние.
	}

	for i := 2; i < len(sortedPoints); i++ {
		currentH := sortedPoints[i].X - sortedPoints[i-1].X
		if math.Abs(currentH-h) > floatTolerance {
			return false, 0
		}
	}
	return true, h
}

// Factorial вычисляет факториал числа
func Factorial(n int) float64 {
	if n < 0 {
		return math.NaN() // Факториал отрицательного числа не определен
	}
	if n == 0 {
		return 1.0
	}
	res := 1.0
	for i := 1; i <= n; i++ {
		res *= float64(i)
	}
	return res
}

// PrintFiniteDifferenceTable выводит таблицу конечных разностей
func PrintFiniteDifferenceTable(points Points, diffTable [][]float64) {
	fmt.Println("\nТаблица конечных разностей:")
	fmt.Printf("%-10s %-10s", "X", "Y")
	for k := 0; k < len(diffTable[0])-1; k++ {
		fmt.Printf("%-10s", fmt.Sprintf("Δ^%d Y", k+1))
	}
	fmt.Println()

	for i := 0; i < len(points); i++ {
		fmt.Printf("%-10.2f %-10.2f", points[i].X, points[i].Y)
		for k := 0; k < len(diffTable[i])-1; k++ {
			if !math.IsNaN(diffTable[i][k+1]) {
				fmt.Printf("%-10.2f", diffTable[i][k+1])
			} else {
				fmt.Printf("%-10s", "") // Пусто, если значение не существует
			}
		}
		fmt.Println()
	}
}

// PrintDividedDifferenceTable выводит таблицу разделенных разностей
func PrintDividedDifferenceTable(points Points, divDiffTable [][]float64) {
	fmt.Println("\nТаблица разделенных разностей (коэффициенты для многочлена Ньютона):")
	// Коэффициенты - это f[x0], f[x0,x1], f[x0,x1,x2], ...
	fmt.Print("f[x_i...]: ")
	for i := 0; i < len(divDiffTable); i++ {
		if len(divDiffTable[i]) > 0 && !math.IsNaN(divDiffTable[i][0]) {
			fmt.Printf("%.4f ", divDiffTable[i][0])
		}
	}
	fmt.Println()

	// Более подробная таблица
	/*
		fmt.Println("\nПодробная таблица разделенных разностей:")
		for i := 0; i < len(points); i++ {
			fmt.Printf("x%-2d: %-8.2f y%-2d: %-8.2f ", i, points[i].X, i, points[i].Y)
			for j := 0; j < len(divDiffTable[i]); j++ {
				if !math.IsNaN(divDiffTable[i][j]) {
					fmt.Printf("%-10.4f", divDiffTable[i][j])
				}
			}
			fmt.Println()
		}
	*/
}
