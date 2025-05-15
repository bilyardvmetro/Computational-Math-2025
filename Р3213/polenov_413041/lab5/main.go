package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var points Points
	var targetX float64
	var originalFunc OriginalFunctionType
	var originalFuncName string
	var err error

	fmt.Println("Выберите способ задания исходных данных:")
	fmt.Println("1: Ввод с клавиатуры")
	fmt.Println("2: Ввод из файла")
	fmt.Println("3: На основе выбранной функции")

	var choiceStr string
	for {
		fmt.Print("Ваш выбор (1-3): ")
		choiceStr, _ = reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		if choiceStr == "1" || choiceStr == "2" || choiceStr == "3" {
			break
		}
		fmt.Println("Некорректный выбор. Пожалуйста, введите 1, 2 или 3.")
	}

	switch choiceStr {
	case "1":
		points, targetX, err = ReadPointsFromKeyboard()
	case "2":
		points, targetX, err = ReadPointsFromFile()
	case "3":
		points, targetX, originalFunc, originalFuncName, err = GeneratePointsFromFunction()
	}

	if err != nil {
		fmt.Printf("Ошибка ввода данных: %v\n", err)
		return
	}

	if len(points) < 2 {
		fmt.Println("Недостаточно узлов для интерполяции (минимум 2).")
		return
	}
	// Убедимся, что точки отсортированы по X для корректной работы всех методов и вывода
	sort.Sort(points)
	fmt.Printf("\nУзлы интерполяции: %s\n", points)
	fmt.Printf("Точка для интерполяции: %.2f\n", targetX)
	if originalFunc != nil {
		fmt.Printf("Исходная функция: %s\n", originalFuncName)
	}

	// --- Интерполяция ---
	fmt.Println("\n--- Результаты интерполяции ---")

	// 1. Многочлен Лагранжа
	yLagrange, errLagrange := LagrangeInterpolate(points, targetX)
	if errLagrange == nil {
		fmt.Printf("Многочлен Лагранжа: P(%.2f) = %.6f\n", targetX, yLagrange)
	} else {
		fmt.Printf("Многочлен Лагранжа: Ошибка - %v\n", errLagrange)
		yLagrange = math.NaN() // Для корректной обработки при построении графика
	}
	lagrangePolyFunc, errLagrangeFunc := LagrangePolynomialFunction(points)
	if errLagrangeFunc != nil {
		fmt.Printf("Не удалось создать функцию многочлена Лагранжа: %v\n", errLagrangeFunc)
		lagrangePolyFunc = nil
	}

	// 2. Многочлен Ньютона с разделенными разностями
	divDiffTable, coeffsNewtonDiv, errNewtonDivTable := BuildDividedDifferenceTable(points)
	var yNewtonDiv = math.NaN()
	var newtonDivPolyFunc func(x float64) float64

	if errNewtonDivTable == nil {
		PrintDividedDifferenceTable(points, divDiffTable) // Выводим таблицу разделенных разностей
		yNewtonDiv, err = NewtonDividedDifferenceInterpolate(points, coeffsNewtonDiv, targetX)
		if err == nil {
			fmt.Printf("Многочлен Ньютона (разд. разн.): P(%.2f) = %.6f\n", targetX, yNewtonDiv)
		} else {
			fmt.Printf("Многочлен Ньютона (разд. разн.): Ошибка вычисления - %v\n", err)
		}
		newtonDivPolyFunc, err = NewtonDividedDifferencePolynomialFunction(points, coeffsNewtonDiv)
		if err != nil {
			fmt.Printf("Не удалось создать функцию многочлена Ньютона (разд. разн.): %v\n", err)
			newtonDivPolyFunc = nil
		}

	} else {
		fmt.Printf("Многочлен Ньютона (разд. разн.): Ошибка построения таблицы - %v\n", errNewtonDivTable)
	}

	// 3. Многочлен Ньютона с конечными разностями
	var yNewtonFin = math.NaN()
	var newtonFinPolyFunc func(x float64) float64
	finiteDiffTable, hStep, errNewtonFinTable := BuildFiniteDifferenceTable(points)

	if errNewtonFinTable == nil {
		fmt.Printf("\nУзлы равноотстоящие с шагом h = %.4f\n", hStep)
		PrintFiniteDifferenceTable(points, finiteDiffTable) // Выводим таблицу конечных разностей

		yNewtonFin, err = NewtonFiniteDifferenceInterpolate(points, finiteDiffTable, hStep, targetX)
		if err == nil {
			fmt.Printf("Многочлен Ньютона (кон. разн.): P(%.2f) = %.6f\n", targetX, yNewtonFin)
		} else {
			fmt.Printf("Многочлен Ньютона (кон. разн.): Ошибка вычисления - %v\n", err)
		}
		newtonFinPolyFunc, err = NewtonFiniteDifferencePolynomialFunction(points, finiteDiffTable, hStep)
		if err != nil {
			fmt.Printf("Не удалось создать функцию многочлена Ньютона (кон. разн.): %v\n", err)
			newtonFinPolyFunc = nil
		}

	} else {
		// Проверяем, является ли ошибка "не равноотстоящие узлы"
		if _, ok := errNewtonFinTable.(*InterpolationError); ok && strings.Contains(errNewtonFinTable.Error(), "не являются равноотстоящими") {
			fmt.Printf("\nМногочлен Ньютона (кон. разн.): %v\n", errNewtonFinTable)
		} else {
			fmt.Printf("\nМногочлен Ньютона (кон. разн.): Ошибка подготовки - %v\n", errNewtonFinTable)
		}
	}

	// Сравнение со значением исходной функции, если она есть
	if originalFunc != nil {
		trueValue := originalFunc(targetX)
		fmt.Printf("\n--- Сравнение с исходной функцией (%s) ---\n", originalFuncName)
		fmt.Printf("Истинное значение f(%.2f) = %.6f\n", targetX, trueValue)
		if !math.IsNaN(yLagrange) {
			fmt.Printf("  Лагранж:         абс. ошибка = %.6e\n", math.Abs(yLagrange-trueValue))
		}
		if !math.IsNaN(yNewtonDiv) {
			fmt.Printf("  Ньютон (разд.):  абс. ошибка = %.6e\n", math.Abs(yNewtonDiv-trueValue))
		}
		if !math.IsNaN(yNewtonFin) {
			fmt.Printf("  Ньютон (конеч.): абс. ошибка = %.6e\n", math.Abs(yNewtonFin-trueValue))
		}
	}

	// --- Построение графика ---
	plotFilename := "interpolation_plot.png"
	fmt.Printf("\nПопытка построить график...\n")

	// Передаем функции для построения многочленов, а не только значения в одной точке
	errPlot := PlotInterpolations(points, targetX, originalFunc, originalFuncName,
		lagrangePolyFunc, newtonDivPolyFunc, newtonFinPolyFunc,
		plotFilename)

	if errPlot != nil {
		fmt.Printf("Ошибка при построении графика: %v\n", errPlot)
	}
}
