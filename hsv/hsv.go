// Copyright Brian Starkey <stark3y@gmail.com> 2017

// image/color.Color implementation for HSV representation
package hsv

import (
	"image/color"
)

// HSVColor defines a color in the Hue-Saturation-Value scheme.
// Hue is a value [0 - 360) specifying the color
// Saturation is a value [0 - 255] specifying the strength of the color
// Value is a value [0 - 255] specifying the brightness of the color
type HSVColor struct {
	H    uint16
	S, V uint8
}

func (h HSVColor) rgba() (r, g, b, a uint32) {
	// Direct implementation of the graph in this image:
	// https://en.wikipedia.org/wiki/HSL_and_HSV#/media/File:HSV-RGB-comparison.svg
	max := uint32(h.V) * 255
	min := uint32(h.V) * uint32(255-h.S)

	h.H %= 360
	segment := h.H / 60
	offset := uint32(h.H % 60)
	mid := ((max - min) * offset) / 60

	switch segment {
	case 0:
		return max, min + mid, min, 0xffff
	case 1:
		return max - mid, max, min, 0xffff
	case 2:
		return min, max, min + mid, 0xffff
	case 3:
		return min, max - mid, max, 0xffff
	case 4:
		return min + mid, min, max, 0xffff
	case 5:
		return max, min, max - mid, 0xffff
	}

	return 0, 0, 0, 0xffff
}

func (h HSVColor) RGBA() color.RGBA {
	r, g, b, _ := h.rgba()

	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xff}
}
