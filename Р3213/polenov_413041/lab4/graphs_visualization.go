package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"log"
	"math"
	"os"
)

func DrawLinearGraph(X, Y []float64, a, b float64, filename string) {
	p := plot.New()
	p.Title.Text = ""
	p.X.Label.Text = "Ox"
	p.Y.Label.Text = "Oy"

	pts := make(plotter.XYs, len(X))
	for i := range X {
		pts[i].X = X[i]
		pts[i].Y = Y[i]
	}

	scatter, _ := plotter.NewScatter(pts)
	scatter.GlyphStyle.Color = plotter.DefaultGlyphStyle.Color
	scatter.GlyphStyle.Radius = vg.Points(3)

	p.Add(scatter)
	p.Legend.Add("Исходные точки", scatter)

	line := plotter.NewFunction(func(x float64) float64 {
		return a*x + b
	})
	line.Color = plotter.DefaultLineStyle.Color
	p.Add(line)
	p.Legend.Add(fmt.Sprintf("y = %.3fx + %.3f", a, b), line)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func DrawQuadraticGraph(X, Y []float64, a, b, c float64, filename string) {
	p := plot.New()
	p.Title.Text = ""
	p.X.Label.Text = "Ox"
	p.Y.Label.Text = "Oy"

	pts := make(plotter.XYs, len(X))
	for i := range X {
		pts[i].X = X[i]
		pts[i].Y = Y[i]
	}

	scatter, _ := plotter.NewScatter(pts)
	scatter.GlyphStyle.Radius = vg.Points(3)
	p.Add(scatter)
	p.Legend.Add("Исходные точки", scatter)

	line := plotter.NewFunction(func(x float64) float64 {
		return a*x*x + b*x + c
	})
	line.Color = plotter.DefaultLineStyle.Color
	p.Add(line)
	p.Legend.Add(fmt.Sprintf("y = %.3fx² + %.3fx + %.3f", a, b, c), line)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func DrawCubicGraph(X, Y []float64, a, b, c, d float64, filename string) {
	p := plot.New()
	p.Title.Text = ""
	p.X.Label.Text = "Ox"
	p.Y.Label.Text = "Oy"

	pts := make(plotter.XYs, len(X))
	for i := range X {
		pts[i].X = X[i]
		pts[i].Y = Y[i]
	}

	scatter, _ := plotter.NewScatter(pts)
	scatter.GlyphStyle.Radius = vg.Points(3)
	p.Add(scatter)
	p.Legend.Add("Исходные точки", scatter)

	line := plotter.NewFunction(func(x float64) float64 {
		return a*x*x*x + b*x*x + c*x + d
	})
	line.Color = plotter.DefaultLineStyle.Color
	p.Add(line)
	p.Legend.Add(fmt.Sprintf("y = %.3fx³ + %.3fx² + %.3fx + %.3f", a, b, c, d), line)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func DrawExponentialGraph(X, Y, accX []float64, a, b float64, filename string) {
	p := plot.New()
	p.Title.Text = ""
	p.X.Label.Text = "Ox"
	p.Y.Label.Text = "Oy"

	pts := make(plotter.XYs, len(X))
	for i := range X {
		pts[i].X = X[i]
		pts[i].Y = Y[i]
	}

	scatter, _ := plotter.NewScatter(pts)
	scatter.GlyphStyle.Radius = vg.Points(3)
	p.Add(scatter)
	p.Legend.Add("Исходные точки", scatter)

	line := plotter.NewFunction(func(x float64) float64 {
		return a * math.Exp(b*x)
	})
	line.Color = plotter.DefaultLineStyle.Color
	p.Add(line)
	p.Legend.Add(fmt.Sprintf("y = %.3fe^{%.3fx}", a, b), line)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		panic(err)
	}
}

// DrawLogGraph Логарифмическая аппроксимация: y = a*ln(x) + b
func DrawLogGraph(X, Y, accX []float64, a, b float64, filename string) {
	p := plot.New()
	p.Title.Text = ""
	p.X.Label.Text = "Ox"
	p.Y.Label.Text = "Oy"

	pts := make(plotter.XYs, len(X))
	for i := range X {
		pts[i].X = X[i]
		pts[i].Y = Y[i]
	}

	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatalf("error creating scatter: %v", err)
	}
	scatter.GlyphStyle.Radius = vg.Points(3)
	p.Add(scatter)
	p.Legend.Add("Исходные точки", scatter)

	line := plotter.NewFunction(func(x float64) float64 {
		if x <= 0 {
			return 0
		}
		return a*math.Log(x) + b
	})
	line.Color = plotter.DefaultLineStyle.Color
	p.Add(line)
	p.Legend.Add(fmt.Sprintf("y = %.3fln(x) + %.3f", a, b), line)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		panic(err)
	}
}

// DrawPowerGraph Степенная аппроксимация: y = a * x^b
func DrawPowerGraph(X, Y, accX []float64, a, b float64, filename string) {
	p := plot.New()
	p.Title.Text = ""
	p.X.Label.Text = "Ox"
	p.Y.Label.Text = "Oy"

	pts := make(plotter.XYs, len(X))
	for i := range X {
		pts[i].X = X[i]
		pts[i].Y = Y[i]
	}

	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatalf("error creating scatter: %v", err)
	}
	scatter.GlyphStyle.Radius = vg.Points(3)
	p.Add(scatter)
	p.Legend.Add("Исходные точки", scatter)

	line := plotter.NewFunction(func(x float64) float64 {
		if x <= 0 {
			return 0
		}
		return a * math.Pow(x, b)
	})
	line.Color = plotter.DefaultLineStyle.Color
	p.Add(line)
	p.Legend.Add(fmt.Sprintf("y = %.3fx^%.3f", a, b), line)

	if err := p.Save(10*vg.Inch, 6*vg.Inch, filename); err != nil {
		panic(err)
	}
}

type Results struct {
	Linear    []float64
	Quadratic []float64
	Cubic     []float64
	Exp       []interface{}
	Log       []interface{}
	Pow       []interface{}
}

func drawAllApproximations(X, Y []float64, results Results) {
	p := plot.New()
	p.Title.Text = "сравнение методов аппроксимации"
	p.X.Label.Text = "Ox"
	p.Y.Label.Text = "Oy"
	p.Add(plotter.NewGrid())

	pts := make(plotter.XYs, len(X))
	for i := range X {
		pts[i].X = X[i]
		pts[i].Y = Y[i]
	}
	scatter, _ := plotter.NewScatter(pts)
	scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(2)
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(scatter)
	p.Legend.Add("исходные точки", scatter)

	xMin, xMax := X[0], X[0]
	for _, x := range X {
		if x < xMin {
			xMin = x
		}
		if x > xMax {
			xMax = x
		}
	}
	xMin -= 0.5
	xMax += 0.5

	addLine := func(f func(x float64) float64, color color.Color, label string) {
		line := plotter.NewFunction(f)
		line.Color = color
		p.Add(line)
		p.Legend.Add(label, line)
	}

	if len(results.Linear) == 2 {
		a, b := results.Linear[0], results.Linear[1]
		addLine(func(x float64) float64 { return a*x + b }, color.RGBA{R: 255, G: 105, B: 180, A: 255},
			fmt.Sprintf("линейн: y = %.3fx + %.3f", a, b))
	}

	if len(results.Quadratic) == 3 {
		a, b, c := results.Quadratic[0], results.Quadratic[1], results.Quadratic[2]
		addLine(func(x float64) float64 { return a*x*x + b*x + c }, color.RGBA{R: 255, G: 20, B: 147, A: 255},
			fmt.Sprintf("квадр: y = %.3fx² + %.3fx + %.3f", a, b, c))
	}

	if len(results.Cubic) == 4 {
		a, b, c, d := results.Cubic[0], results.Cubic[1], results.Cubic[2], results.Cubic[3]
		addLine(func(x float64) float64 { return a*x*x*x + b*x*x + c*x + d }, color.RGBA{R: 238, G: 130, B: 238, A: 255},
			fmt.Sprintf("кубич: y = %.3fx³ + %.3fx² + %.3fx + %.3f", a, b, c, d))
	}

	if len(results.Exp) >= 2 {
		a, b := results.Exp[0].(float64), results.Exp[1].(float64)
		addLine(func(x float64) float64 { return a * math.Exp(b*x) }, color.Gray{Y: 100},
			fmt.Sprintf("экспоненц: y = %.3fe^%.3fx", a, b))
	}

	if len(results.Log) >= 2 {
		a, b := results.Log[0].(float64), results.Log[1].(float64)
		addLine(func(x float64) float64 {
			if x <= 0 {
				return math.NaN()
			}
			return a*math.Log(x) + b
		}, color.RGBA{R: 200, G: 200, B: 200, A: 255},
			fmt.Sprintf("лог: y = %.3fln(x) + %.3f", a, b))
	}

	if len(results.Pow) >= 2 {
		a, b := results.Pow[0].(float64), results.Pow[1].(float64)
		addLine(func(x float64) float64 {
			if x <= 0 {
				return math.NaN()
			}
			return a * math.Pow(x, b)
		}, color.RGBA{R: 47, G: 79, B: 79, A: 255},
			fmt.Sprintf("степен: y = %.3fx^%.3f", a, b))
	}

	if err := p.Save(12*vg.Inch, 8*vg.Inch, "graphics/all_approximations.png"); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка сохранения графика: %v\n", err)
	}
}
