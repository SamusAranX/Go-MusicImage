package curves

type Curve interface {
	next() Point
	NextIntegral() IntegralPoint

	Radius() float64
}
