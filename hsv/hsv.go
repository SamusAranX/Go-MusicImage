package hsv

import (
	"image/color"
	"math"
)

// HSVColor defines a color in the Hue-Saturation-Value scheme.
// Hue is a double [0 - 360.0] specifying the color
// Saturation is a double [0 - 100.0] specifying the strength of the color
// Value is a double [0 - 100.0] specifying the brightness of the color
type Color struct {
	H    float64
	S, V float64
}

func (hsv Color) rgba() (r, g, b float64) {
	h := math.Mod(hsv.H, 360)
	s := hsv.S / 100
	v := hsv.V / 100

	if hsv.S > 100 || hsv.V > 100 {
		panic("Invalid saturation or value values")
	}

	if s == 0 {
		return v, v, v
	}

	h /= 60
	i, f := math.Modf(h)

	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	switch i {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	default:
		return v, p, q
	}
}

func (hsv Color) RGBA() color.RGBA {
	fr, fg, fb := hsv.rgba()
	r := uint8(math.Round(fr * 0xff))
	g := uint8(math.Round(fg * 0xff))
	b := uint8(math.Round(fb * 0xff))
	return color.RGBA{r, g, b, 0xff}
}

func (hsv Color) RGBA64() color.RGBA64 {
	fr, fg, fb := hsv.rgba()
	r := uint16(math.Round(fr * 0xffff))
	g := uint16(math.Round(fg * 0xffff))
	b := uint16(math.Round(fb * 0xffff))
	return color.RGBA64{r, g, b, 0xffff}
}
