package musicoders

import (
	"fmt"
	wav "github.com/youpy/go-wav"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
)

type Encoder struct {
	InputSound string

	Diameter   uint32
	Separation float64

	DeepColor bool
}

func (e Encoder) Encode(OutputImage string) bool {
	wavFile, err := os.Open(e.InputSound)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer wavFile.Close()

	reader := wav.NewReader(wavFile)
	format, err := reader.Format()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if format.BitsPerSample > 24 {
		fmt.Println("WAV files with more than 24 bits per sample are not supported.")
		return false
	}

	offset := int(math.Pow(2, float64(format.BitsPerSample)) / 2)

	rect := image.Rect(0, 0, 2048, 2048)
	img := image.NewRGBA(rect)

	spiral := NewSpiral(e.Diameter, e.Separation)
	spiral.Center = IntegralPoint{1024, 1024}

	var totalSamples uint32

	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}

		for _, sample := range samples {
			p := spiral.Next()

			val := uint32(sample.Values[0])

			b1 := uint8(val & 0xff0000 >> 16)
			b2 := uint8(val & 0x00ff00 >> 8)
			b3 := uint8(val & 0x0000ff)

			col := color.RGBA{b1, b2, b3, 0xff}

			if totalSamples > 200000 && totalSamples < 200064 {
				fmt.Printf("0x%X: %d %d %d\n", val, b1, b2, b3)
				fmt.Printf("%d %d %d\n", sample.Values[0], offset, int(sample.Values[0])+offset)
				fmt.Println("----------")
			}

			img.Set(p.X, p.Y, col)

			totalSamples++

			// fmt.Printf("(%d) L/R: %d/%d\n", i, reader.IntValue(sample, 0), reader.IntValue(sample, 1))
		}
	}

	fmt.Printf("%d samples written\n", totalSamples)

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
