package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// OriginalFunctionType представляет тип исходной функции для интерполяции
type OriginalFunctionType func(x float64) float64

var predefinedFunctions = map[int]struct {
	Name string
	Func OriginalFunctionType
}{
	1: {"sin(x)", math.Sin},
	2: {"4x+5", func(x float64) float64 { return 4*x + 5 }},
	3: {"sqrt(x)", func(x float64) float64 {
		if x < 0 {
			return math.NaN() // Обработка отрицательных значений
		}
		return math.Sqrt(x)
	}},
	4: {"x^2", func(x float64) float64 { return x * x }},
	5: {"x^3", func(x float64) float64 { return x * x * x }},
}

// ReadPointsFromKeyboard читает узлы интерполяции с клавиатуры
func ReadPointsFromKeyboard() (Points, float64, error) {
	reader := bufio.NewReader(os.Stdin)
	var points Points
	var targetX float64

	fmt.Print("Введите точку x для интерполяции: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("некорректное значение для точки x: %v", err)
	}
	targetX = val

	fmt.Println("Введите узлы интерполяции (пары x y, разделенные пробелом). Напишите 'quit' для завершения:")
	for {
		fmt.Print("Узел (x y) или 'quit': ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if strings.ToLower(line) == "quit" {
			break
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Println("Некорректный формат ввода. Пожалуйста, введите два числа, разделенных пробелом, или 'quit'.")
			continue
		}

		x, errX := strconv.ParseFloat(parts[0], 64)
		y, errY := strconv.ParseFloat(parts[1], 64)

		if errX != nil || errY != nil {
			fmt.Println("Ошибка парсинга чисел. Пожалуйста, введите корректные числа.")
			continue
		}
		points = append(points, Point{X: x, Y: y})
	}

	if len(points) < 2 {
		return points, targetX, fmt.Errorf("необходимо как минимум 2 узла для интерполяции")
	}
	sort.Sort(points) // Сортируем точки по X для удобства
	return points, targetX, nil
}

// ReadPointsFromFile читает узлы интерполяции из файла
func ReadPointsFromFile() (Points, float64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите имя файла с данными: ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка открытия файла %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var points Points
	var targetX float64
	isFirstLine := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" { // Пропускаем пустые строки
			continue
		}

		if isFirstLine {
			val, err := strconv.ParseFloat(line, 64)
			if err != nil {
				return nil, 0, fmt.Errorf("некорректное значение для точки x в первой строке файла: %v", err)
			}
			targetX = val
			isFirstLine = false
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, targetX, fmt.Errorf("некорректный формат узла в файле: '%s'. Ожидается 'x y'", line)
		}

		x, errX := strconv.ParseFloat(parts[0], 64)
		y, errY := strconv.ParseFloat(parts[1], 64)

		if errX != nil || errY != nil {
			return nil, targetX, fmt.Errorf("ошибка парсинга чисел в узле '%s': %v, %v", line, errX, errY)
		}
		points = append(points, Point{X: x, Y: y})
	}

	if err := scanner.Err(); err != nil {
		return nil, targetX, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	if isFirstLine { // Если файл был пуст или содержал только пустые строки
		return nil, 0, fmt.Errorf("файл не содержит данных для точки x")
	}
	if len(points) < 2 {
		return points, targetX, fmt.Errorf("необходимо как минимум 2 узла для интерполяции из файла")
	}
	sort.Sort(points) // Сортируем точки по X
	return points, targetX, nil
}

// GeneratePointsFromFunction генерирует узлы на основе выбранной функции и интервала
func GeneratePointsFromFunction() (Points, float64, OriginalFunctionType, string, error) {
	reader := bufio.NewReader(os.Stdin)
	var points Points
	var targetX float64
	var chosenFunc OriginalFunctionType
	var funcName string

	fmt.Println("Выберите функцию для интерполяции:")
	for i := 1; i <= len(predefinedFunctions); i++ {
		fmt.Printf("%d: %s\n", i, predefinedFunctions[i].Name)
	}

	var choice int
	for {
		fmt.Print("Ваш выбор (номер): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		c, err := strconv.Atoi(input)
		if err == nil && c >= 1 && c <= len(predefinedFunctions) {
			choice = c
			chosenFunc = predefinedFunctions[choice].Func
			funcName = predefinedFunctions[choice].Name
			break
		}
		fmt.Println("Некорректный выбор. Пожалуйста, введите номер из списка.")
	}

	var a, b float64
	var numPoints int

	fmt.Print("Введите начало интервала a: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return nil, 0, nil, "", fmt.Errorf("некорректное значение для a: %v", err)
	}
	a = val

	fmt.Print("Введите конец интервала b: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return nil, 0, nil, "", fmt.Errorf("некорректное значение для b: %v", err)
	}
	b = val

	if a >= b {
		return nil, 0, nil, "", fmt.Errorf("начало интервала 'a' должно быть меньше конца 'b'")
	}

	for {
		fmt.Print("Введите количество точек на интервале (>=2): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		np, err := strconv.Atoi(input)
		if err == nil && np >= 2 {
			numPoints = np
			break
		}
		fmt.Println("Некорректное количество точек. Должно быть целое число >= 2.")
	}

	// Генерируем точки
	step := (b - a) / float64(numPoints-1)
	for i := 0; i < numPoints; i++ {
		x := a + float64(i)*step
		y := chosenFunc(x)
		if math.IsNaN(y) {
			// Например, для sqrt(x) при x < 0 на интервале
			return nil, 0, nil, "", fmt.Errorf("ошибка вычисления функции для x=%.2f (например, sqrt от отрицательного числа). Пожалуйста, проверьте интервал", x)
		}
		points = append(points, Point{X: x, Y: y})
	}

	fmt.Print("Введите точку x для интерполяции: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	val, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return nil, 0, nil, "", fmt.Errorf("некорректное значение для точки x: %v", err)
	}
	targetX = val

	// Точки уже отсортированы, так как генерировались по порядку
	return points, targetX, chosenFunc, funcName, nil
}
