package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const ACCURACY = 0.0001

func DrawSingleFunction(f func(x float64) float64, start float64, end float64, accuracy float64) error {
	p := plot.New()

	p.Title.Text = "График функции"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	points := make(plotter.XYs, int((end-start)/accuracy))
	for i := range points {
		x := start + accuracy*float64(i)
		points[i].X = x
		points[i].Y = f(x)
	}

	line, err := plotter.NewLine(points)
	if err != nil {
		return err
	}

	p.Add(line)

	if err := p.Save(6*vg.Inch, 6*vg.Inch, "plot.png"); err != nil {
		return err
	}
	return nil
}

func DrawTwoFunctions(f1 func(x float64, y float64) float64, f2 func(x float64, y float64) float64, xStart float64, xEnd float64, yStart float64, yEnd float64, accuracy float64) error {
	p := plot.New()

	p.Title.Text = "График системы функций"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	points1 := make(plotter.XYs, 0)
	points2 := make(plotter.XYs, 0)

	for x := xStart; x <= xEnd; x += accuracy {
		for y := yStart; y <= yEnd; y += accuracy {
			if Abs(f1(x, y)) <= accuracy {
				points1 = append(points1, plotter.XY{X: x, Y: y})
			}
			if Abs(f2(x, y)) <= accuracy {
				points2 = append(points2, plotter.XY{X: x, Y: y})
			}
		}
	}

	pointsList := make([]plotter.XYs, 2)
	pointsList[0], pointsList[1] = points1, points2
	for index, points := range pointsList {
		line, err := plotter.NewLine(points)
		if err != nil {
			return err
		}
		line.Color = plotutil.Color(index)
		p.Add(line)
	}

	if err := p.Save(6*vg.Inch, 6*vg.Inch, "plot.png"); err != nil {
		return err
	}
	return nil

}
