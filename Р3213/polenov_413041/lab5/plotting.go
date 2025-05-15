package main

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// PlotInterpolations строит и сохраняет график
func PlotInterpolations(
	basePoints Points,
	targetX float64,
	originalFunc OriginalFunctionType,
	originalFuncName string,
	lagrangeFunc func(x float64) float64,
	newtonDivFunc func(x float64) float64,
	newtonFinFunc func(x float64) float64, // Может быть nil, если неприменимо
	filename string,
) error {
	p := plot.New()

	p.Title.Text = "Интерполяция функций"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Legend.Top = true

	// Определяем диапазон для X на графике
	minX, maxX := math.Inf(1), math.Inf(-1)
	if len(basePoints) > 0 {
		// Сортируем basePoints если они еще не отсортированы (хотя должны быть)
		sortedBasePoints := make(Points, len(basePoints))
		copy(sortedBasePoints, basePoints)
		sort.Sort(sortedBasePoints)

		minX = sortedBasePoints[0].X
		maxX = sortedBasePoints[len(sortedBasePoints)-1].X

		// Добавим небольшой отступ
		padding := (maxX - minX) * 0.1
		if padding == 0 { // если все точки X одинаковы или одна точка
			padding = 1.0
		}
		minX -= padding
		maxX += padding
	} else { // если нет узлов, но есть, например, targetX
		minX = targetX - 5
		maxX = targetX + 5
	}
	if targetX < minX {
		minX = targetX - (maxX-targetX)*0.1 - 1
	}
	if targetX > maxX {
		maxX = targetX + (targetX-minX)*0.1 + 1
	}

	// 1. Исходная функция (если есть)
	if originalFunc != nil {
		lineOrig := plotter.NewFunction(originalFunc)
		lineOrig.Color = color.RGBA{B: 255, A: 255} // Синий
		lineOrig.Width = vg.Points(2)
		lineOrig.Samples = 200 // Больше точек для гладкости
		p.Add(lineOrig)
		p.Legend.Add(fmt.Sprintf("Исходная: %s", originalFuncName), lineOrig)
	}

	// 2. Многочлен Лагранжа
	if lagrangeFunc != nil {
		lineLagrange := plotter.NewFunction(lagrangeFunc)
		lineLagrange.Color = color.RGBA{R: 255, A: 255} // Красный
		lineLagrange.Width = vg.Points(1.5)
		lineLagrange.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(3)} // Пунктир
		lineLagrange.Samples = 200
		p.Add(lineLagrange)
		p.Legend.Add("Лагранж", lineLagrange)
	}

	// 3. Многочлен Ньютона (разделенные разности)
	// Обычно он идентичен Лагранжу, но для демонстрации можно добавить
	if newtonDivFunc != nil {
		lineNewtonDiv := plotter.NewFunction(newtonDivFunc)
		lineNewtonDiv.Color = color.RGBA{G: 150, R: 50, A: 255} // Темно-зеленый/Оранжевый
		lineNewtonDiv.Width = vg.Points(1.5)
		lineNewtonDiv.LineStyle.Dashes = []vg.Length{vg.Points(2), vg.Points(2)} // Точки
		lineNewtonDiv.Samples = 200
		p.Add(lineNewtonDiv)
		p.Legend.Add("Ньютон (разд. разн.)", lineNewtonDiv)
	}

	// 4. Многочлен Ньютона (конечные разности)
	if newtonFinFunc != nil {
		lineNewtonFin := plotter.NewFunction(newtonFinFunc)
		lineNewtonFin.Color = color.RGBA{R: 128, G: 0, B: 128, A: 255} // Фиолетовый
		lineNewtonFin.Width = vg.Points(1.5)
		lineNewtonFin.LineStyle.Dashes = []vg.Length{vg.Points(8), vg.Points(4), vg.Points(2), vg.Points(4)} // Штрих-пунктир
		lineNewtonFin.Samples = 200
		p.Add(lineNewtonFin)
		p.Legend.Add("Ньютон (кон. разн.)", lineNewtonFin)
	}

	// 5. Узлы интерполяции
	if len(basePoints) > 0 {
		ptsPlot := make(plotter.XYs, len(basePoints))
		for i, pt := range basePoints {
			ptsPlot[i].X = pt.X
			ptsPlot[i].Y = pt.Y
		}
		scatter, err := plotter.NewScatter(ptsPlot)
		if err != nil {
			return fmt.Errorf("ошибка создания scatter plot для узлов: %v", err)
		}
		scatter.GlyphStyle.Color = color.RGBA{A: 255} // Черный
		scatter.GlyphStyle.Radius = vg.Points(4)
		p.Add(scatter)
		p.Legend.Add("Узлы интерполяции", scatter)
	}

	// Сохраняем график в файл
	if err := p.Save(8*vg.Inch, 6*vg.Inch, filename); err != nil {
		return fmt.Errorf("ошибка сохранения графика: %v", err)
	}
	fmt.Printf("График сохранен в файл: %s\n", filename)
	return nil
}
