package micoders

import (
	// "bufio"
	"fmt"
	// "image"
	"math"
	// "os"
	// flags "github.com/jessevdk/go-flags"
)

func max(is ...float64) float64 {
	max := is[0]
	for _, i := range is[1:] {
		if i > max {
			max = i
		}
	}
	return max
}

type spiral struct {
	Center IntegralPoint

	lastPoint  IntegralPoint
	lastPoints map[IntegralPoint]struct{}

	Separation, StartRadius, Theta, Radius float64
}

func NewSpiral(diameter uint32, separation float64) spiral {
	sp := spiral{}
	sp.Center = IntegralPoint{127, 127}
	sp.Separation = separation
	sp.StartRadius = float64(diameter) / 2
	sp.Theta = 0
	sp.Radius = max(1, sp.StartRadius)

	sp.lastPoints = map[IntegralPoint]struct{}{}

	return sp
}

func (s spiral) String() string {
	return fmt.Sprintf("{%d %d} %.2f", s.Center.X, s.Center.Y, s.Radius)
}

// returns a Point with float64 coordinates
func (s *spiral) next() Point {
	p := Point{}

	s.Theta += s.Separation / s.Radius
	s.Radius = s.StartRadius + s.Separation*s.Theta/math.Pi

	p.X = s.Radius*math.Cos(s.Theta) + float64(s.Center.X)
	p.Y = s.Radius*math.Sin(s.Theta) + float64(s.Center.Y)

	return p
}

// returns an IntegralPoint with int coordinates
func (s *spiral) Next() IntegralPoint {
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
