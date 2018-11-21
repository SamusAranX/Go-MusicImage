package musicoders

import (
	"fmt"
	wav "github.com/youpy/go-wav"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// Decoder what
type Decoder struct {
	InputImage string

	Diameter   uint32
	Separation float64

	SampleRate uint32
	ChannelNum uint16
	SampleBits uint16
}

// Decode wat
func (d Decoder) Decode(OutputSound string) bool {
	imgFile, err := os.Open(d.InputImage)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer imgFile.Close()

	wavFile, err := os.Create(OutputSound)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer wavFile.Close()

	img, _, err := image.Decode(imgFile)
	imgSize := img.Bounds().Size()

	spiral := NewSpiral(d.Diameter, d.Separation)
	spiral.Center = IntegralPoint{imgSize.X / 2, imgSize.Y / 2}

	// allocate enough space for 7 minutes of 44100 Hz audio
	var numSamples uint32 = 7 * 60 * 44100
	var samples = make([]wav.Sample, numSamples)

	offset := int(math.Pow(2, float64(d.SampleBits)) / 2)

	writer := wav.NewWriter(wavFile, numSamples, d.ChannelNum, d.SampleRate, d.SampleBits)

	var sampleIdx uint32

	for {
		p := spiral.Next()

		col := img.At(p.X, p.Y)
		_r, _g, _b, _a := col.RGBA()

		r := uint32(math.Round(float64(_r) / 65535.0 * 255.0))
		g := uint32(math.Round(float64(_g) / 65535.0 * 255.0))
		b := uint32(math.Round(float64(_b) / 65535.0 * 255.0))
		a := uint32(math.Round(float64(_a) / 65535.0 * 255.0))

		if a < 255 || (r == 128 && g == 128 && b == 128) || sampleIdx >= 10000000 {
			fmt.Printf("Alpha detected, reached end of track at %d samples\n", sampleIdx)
			break
		}

		val := uint32(0x000000)
		val |= r << 16
		val |= g << 8
		val |= b

		s := wav.Sample{}
		s.Values[0] = int(val)
		s.Values[1] = int(val)

		if sampleIdx > 200000 && sampleIdx < 200100 {
			fmt.Printf("0x%X: %d %d %d\n", val, r, g, b)
			fmt.Printf("%d %d %d\n", val, offset, int(val)-offset)
			fmt.Println("----------")
		}

		samples[sampleIdx] = s
		sampleIdx++
	}

	fmt.Println("Writing samples to file…")
	err = writer.WriteSamples(samples[:sampleIdx])
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
