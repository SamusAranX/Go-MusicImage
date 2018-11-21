package musicoders

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"musicimage/hsv"
	"os"
)

type TestEncoder struct {
	Turns float64

	Diameter   uint32
	Separation float64
}

func (e TestEncoder) Encode(OutputImage string) bool {
	rect := image.Rect(0, 0, 2048, 2048)
	img := image.NewRGBA(rect)

	spiral := NewSpiral(e.Diameter, e.Separation)
	spiral.Center = IntegralPoint{1024, 1024}

	hsv := hsv.HSVColor{0, 255, 255}

	var totalPixels uint32
	for spiral.Theta < (math.Pi * e.Turns * 2) {
		p := spiral.Next()
		img.Set(p.X, p.Y, hsv.RGBA())

		hsv.H += 2
		hsv.H %= 360

		totalPixels++

	}

	fmt.Printf("%d pixels drawn\n", totalPixels)

	radRounded := int(math.Ceil(spiral.Radius))

	subRect := image.Rect(spiral.Center.X-radRounded, spiral.Center.Y-radRounded, spiral.Center.X+radRounded, spiral.Center.Y+radRounded)

	f, err := os.Create(OutputImage)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f.Close()

	png.Encode(f, img.SubImage(subRect))

	return true
}
