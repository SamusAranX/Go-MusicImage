package curves

import (
	"fmt"
)

// https://oeis.org/A174344
// https://www.desmos.com/calculator/augmltextm

type Square struct {
	Center IntegralPoint

	lastPoint  IntegralPoint
	lastPoints map[IntegralPoint]struct{}

	Progress uint64
}

func NewSquare() Square {
	sq := Square{}
	sq.lastPoints = map[IntegralPoint]struct{}{}

	return sq
}

func (sq Square) String() string {
	return fmt.Sprintf("{%d %d} %d", sq.Center.X, sq.Center.Y, sq.Progress)
}

// next returns the next point in the Square as a Point
func (s *Square) next() Point {
	p := Point{}
	//p.X = s.radius*math.Cos(s.Theta) + float64(s.Center.X)
	//p.Y = s.radius*math.Sin(s.Theta) + float64(s.Center.Y)
	//
	//// advance after constructing Point
	//s.Theta += s.Separation / s.radius
	//s.radius = s.StartRadius + s.Separation*s.Theta/math.Pi

	return p
}

// NextIntegral returns the next point in the Square as an IntegralPoint
func (s *Square) NextIntegral() IntegralPoint {
	pr := s.next().Round()

	_, exists := s.lastPoints[pr]
	for exists {
		pr = s.next().Round()
		_, exists = s.lastPoints[pr]
	}

	if !exists {
		s.lastPoints[pr] = struct{}{}
	}

	return pr
}
