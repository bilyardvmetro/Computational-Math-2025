package main

import "fmt"

// Point представляет точку (x, y)
type Point struct {
	X, Y float64
}

// Points - это срез точек
type Points []Point

// String для красивого вывода точки
func (p Point) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", p.X, p.Y)
}

// String для красивого вывода среза точек
func (pts Points) String() string {
	s := "["
	for i, p := range pts {
		s += p.String()
		if i < len(pts)-1 {
			s += ", "
		}
	}
	s += "]"
	return s
}

// Len, Swap, Less для сортировки Points по X
func (pts Points) Len() int           { return len(pts) }
func (pts Points) Swap(i, j int)      { pts[i], pts[j] = pts[j], pts[i] }
func (pts Points) Less(i, j int) bool { return pts[i].X < pts[j].X }
