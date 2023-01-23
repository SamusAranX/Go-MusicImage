package curves

import (
	"fmt"
	"math"
)

type Spiral struct {
	Center IntegralPoint

	lastPoint  IntegralPoint
	lastPoints map[IntegralPoint]struct{}

	Separation, StartRadius, Theta, radius float64
}

func NewSpiral(diameter uint16, separation uint8) Spiral {
	sp := Spiral{}
	sp.Separation = float64(separation)
	sp.StartRadius = float64(diameter) / 2
	sp.Theta = 0
	sp.radius = sp.StartRadius
	if sp.radius < 1 {
		sp.radius = 1
	}

	sp.lastPoints = map[IntegralPoint]struct{}{}

	return sp
}

func (s Spiral) String() string {
	return fmt.Sprintf("{%d %d} %.2f", s.Center.X, s.Center.Y, s.radius)
}

func (s Spiral) Radius() float64 {
	return s.radius
}

// next returns the next point in the Spiral as a Point
func (s *Spiral) next() Point {
	p := Point{}
	p.X = s.radius*math.Cos(s.Theta) + float64(s.Center.X)
	p.Y = s.radius*math.Sin(s.Theta) + float64(s.Center.Y)

	// advance after constructing Point
	s.Theta += s.Separation / s.radius
	s.radius = s.StartRadius + s.Separation*s.Theta/math.Pi

	return p
}

// NextIntegral returns the next point in the Spiral as an IntegralPoint
func (s *Spiral) NextIntegral() IntegralPoint {
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
