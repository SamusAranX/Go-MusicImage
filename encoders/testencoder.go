package encoders

import (
	"fmt"
	"image"
	// "image/color"
	"image/png"
	"math"
	"os"

	"musicimage/encoders/curves"
	"musicimage/hsv"
)

type TestEncoder struct {
	SharedOptions
}

func (e TestEncoder) Encode() error {
	rect := image.Rect(0, 0, 2048, 2048)

	spiral := curves.NewSpiral(e.Diameter, e.Separation)
	spiral.Center = curves.IntegralPoint{X: 1024, Y: 1024}

	hsv := hsv.Color{H: 0, S: 100, V: 100}

	var img image.RGBA
	var img64 image.RGBA64

	if e.DeepColor {
		img64 = *image.NewRGBA64(rect)
	} else {
		img = *image.NewRGBA(rect)
	}

	var totalPixels uint32
	for spiral.Theta < (math.Pi * e.TestEncoderOptions.Turns * 2) {
		p := spiral.Next()

		if e.DeepColor {
			col := hsv.RGBA64()
			img64.Set(p.X, p.Y, col)
		} else {
			col := hsv.RGBA()
			img.Set(p.X, p.Y, col)
		}

		hsv.H += 60

		totalPixels++
	}

	fmt.Printf("%d pixels drawn\n", totalPixels)

	radRounded := int(math.Ceil(spiral.Radius))
	subRect := image.Rect(spiral.Center.X-radRounded, spiral.Center.Y-radRounded, spiral.Center.X+radRounded, spiral.Center.Y+radRounded)

	f, err := os.Create(e.OutFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if e.DeepColor {
		png.Encode(f, img64.SubImage(subRect))
	} else {
		png.Encode(f, img.SubImage(subRect))
	}

	return nil
}
