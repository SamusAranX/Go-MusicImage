package curves

import (
	"fmt"
	"math"
)

type Spiral struct {
	Center IntegralPoint

	lastPoint  IntegralPoint
	lastPoints map[IntegralPoint]struct{}

	Separation, StartRadius, Theta, Radius float64
}

func NewSpiral(diameter uint32, separation float64) Spiral {
	sp := Spiral{}
	sp.Separation = separation
	sp.StartRadius = float64(diameter) / 2
	sp.Theta = 0
	sp.Radius = sp.StartRadius
	if sp.Radius < 1 {
		sp.Radius = 1
	}

	sp.lastPoints = map[IntegralPoint]struct{}{}

	return sp
}

func (s Spiral) String() string {
	return fmt.Sprintf("{%d %d} %.2f", s.Center.X, s.Center.Y, s.Radius)
}

// returns a Point with float64 coordinates
func (s *Spiral) next() Point {
	p := Point{}
	p.X = s.Radius*math.Cos(s.Theta) + float64(s.Center.X)
	p.Y = s.Radius*math.Sin(s.Theta) + float64(s.Center.Y)

	// advance after constructing Point
	s.Theta += s.Separation / s.Radius
	s.Radius = s.StartRadius + s.Separation*s.Theta/math.Pi

	return p
}

// returns an IntegralPoint with int coordinates
func (s *Spiral) Next() IntegralPoint {
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
