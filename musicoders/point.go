package musicoders

import (
	"math"
)

type Point struct {
	X, Y float64
}

func (p Point) Round() IntegralPoint {
	ip := IntegralPoint{}
	ip.X = int(math.Round(p.X))
	ip.Y = int(math.Round(p.Y))
	return ip
}

type IntegralPoint struct {
	X, Y int
}

func (p IntegralPoint) Add(ap IntegralPoint) IntegralPoint {
	return IntegralPoint{p.X + ap.X, p.Y + ap.Y}
}

func (p IntegralPoint) Subtract(sp IntegralPoint) IntegralPoint {
	ap := IntegralPoint{-sp.X, -sp.Y}
	return p.Add(ap)
}

func (p IntegralPoint) DistanceFrom(dp IntegralPoint) float64 {
	x1, y1 := p.X, p.Y
	x2, y2 := dp.X, dp.Y

	left := math.Pow(float64(x2-x1), 2)
	right := math.Pow(float64(y2-y1), 2)

	return math.Sqrt(left + right)
}
