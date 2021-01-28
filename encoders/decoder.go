package encoders

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"

	wav "github.com/youpy/go-wav"
	"musicimage/encoders/curves"
)

type Decoder struct {
	SharedOptions
}

func (d Decoder) Decode() error {
	imgFile, err := os.Open(d.InFile)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	wavFile, err := os.Create(d.OutFile)
	if err != nil {
		return err
	}
	defer wavFile.Close()

	img, _, err := image.Decode(imgFile)
	imgSize := img.Bounds().Size()

	useDeepColor := d.DeepColor || img.ColorModel() == color.RGBA64Model || img.ColorModel() == color.NRGBA64Model
	var colorDepth uint32 = 24
	if useDeepColor {
		colorDepth *= 2
	}

	bitsPerSample := d.DecoderOptions.BitDepth
	channelNum := d.DecoderOptions.ChannelNum
	sampleRate := d.DecoderOptions.SampleRate
	spiral := curves.NewSpiral(d.Diameter, d.Separation)
	spiral.Center = curves.IntegralPoint{X: imgSize.X / 2, Y: imgSize.Y / 2}

	sampleMask := uint64(math.Pow(2, float64(bitsPerSample)) - 1)

	// allocate enough space for 10 minutes of 44100 Hz stereo audio
	var numSamples uint32 = 2 * 10 * 60 * 44100
	var samples = make([]wav.Sample, numSamples)

	// offset := int(math.Pow(2, float64(d.DecoderOptions.BitDepth)) / 2)

	writer := wav.NewWriter(wavFile, numSamples, channelNum, sampleRate, bitsPerSample)

	var sampleIdx uint32

	var pixelInt uint64
	var shiftedBy uint32
	for {
		p := spiral.Next()

		col := img.At(p.X, p.Y)
		_r, _g, _b, _a := col.RGBA()

		if _a < 32768 || sampleIdx >= 10000000 {
			fmt.Printf("Alpha detected, reached end of track at %d samples\n", sampleIdx)
			break
		}

		if !useDeepColor {
			_r = uint32(math.Round(float64(_r) / 65535.0 * 255.0))
			_g = uint32(math.Round(float64(_g) / 65535.0 * 255.0))
			_b = uint32(math.Round(float64(_b) / 65535.0 * 255.0))
		}

		pixelInt <<= bitsPerSample
		pixelInt |= uint64(_r)
		shiftedBy += uint32(bitsPerSample)

		pixelInt <<= bitsPerSample
		pixelInt |= uint64(_g)
		shiftedBy += uint32(bitsPerSample)

		pixelInt <<= bitsPerSample
		pixelInt |= uint64(_b)
		shiftedBy += uint32(bitsPerSample)

		if shiftedBy == colorDepth {
			for shiftedBy > 0 {
				shiftedBy -= uint32(bitsPerSample)

				s := wav.Sample{}

				// fmt.Printf("%024b\n", pixelInt)
				sampleVal := pixelInt >> uint64(shiftedBy)
				sampleValInt := int(sampleVal & sampleMask)

				// fmt.Printf("%024b\n", sampleVal)
				s.Values[0] = sampleValInt

				if channelNum == 2 {
					shiftedBy -= uint32(bitsPerSample)
					s.Values[1] = int(pixelInt >> uint64(shiftedBy))
				} else {
					s.Values[1] = sampleValInt
				}

				samples[sampleIdx] = s
				sampleIdx += uint32(channelNum)
			}

			pixelInt = 0
			shiftedBy = 0
		}
	}

	fmt.Println("Writing samples to file…")
	err = writer.WriteSamples(samples[:sampleIdx])
	if err != nil {
		return err
	}

	return nil
}
