package methods

func SolveByRectangleLeft(f func(x float64) float64, a float64, b float64, n int) float64 {
	h := (b - a) / float64(n)
	var sum = 0.0

	for i := 0; i < n; i++ {
		sum += f(a + float64(i)*h)
	}

	return sum * h
}

func SolveByRectangleRight(f func(x float64) float64, a float64, b float64, n int) float64 {
	h := (b - a) / float64(n)
	var sum = 0.0

	for i := 1; i <= n; i++ {
		sum += f(a + float64(i)*h)
	}

	return sum * h
}

func SolveByRectangleMid(f func(x float64) float64, a float64, b float64, n int) float64 {
	h := (b - a) / float64(n)
	var sum = 0.0

	for i := 0; i < n; i++ {
		sum += f(a + (float64(i)+0.5)*h)
	}

	return sum * h
}

func SolveByTrapezia(f func(x float64) float64, a float64, b float64, n int) float64 {
	h := (b - a) / float64(n)
	sum := (f(a) + f(b)) / 2

	for i := 1; i < n-1; i++ {
		sum += f(a + float64(i)*h)
	}

	return sum * h
}

func SolveBySimpson(f func(x float64) float64, a float64, b float64, n int) float64 {
	h := (b - a) / float64(n)
	sum := f(a) + f(b)

	sumOdd := 0.0
	sumEven := 0.0

	for i := 1; i < n; i++ {
		if i%2 == 0 {
			sumEven += f(a + float64(i)*h)
		} else {
			sumOdd += f(a + float64(i)*h)
		}
	}

	return (sum + 2*sumEven + 4*sumOdd) * (h / 3)
}
