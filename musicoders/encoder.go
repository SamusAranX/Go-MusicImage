package musicoders

import (
	"errors"
	"fmt"
	"github.com/SamusAranX/musicimage/musicoders/curves"
	wav "github.com/youpy/go-wav"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
)

type Encoder struct {
	SharedOptions
}

func (e Encoder) Encode() error {
	wavFile, err := os.Open(e.InFile)
	if err != nil {
		return err
	}
	defer wavFile.Close()

	reader := wav.NewReader(wavFile)
	format, err := reader.Format()
	if err != nil {
		return err
	}

	var totalSamples int

	// color depth in bits per pixel, ignoring alpha channel
	var colorDepth uint32 = 24
	useDeepColor := e.DeepColor
	if useDeepColor {
		colorDepth *= 2
	}

	// how many samples are going to be used per pixel (alpha is always 100%)
	samplesPerPixel := colorDepth / uint32(format.BitsPerSample) / uint32(format.NumChannels)

	fmt.Println(colorDepth, format.BitsPerSample, format.NumChannels)
	fmt.Println(colorDepth, uint32(format.BitsPerSample), uint32(format.NumChannels))
	fmt.Println(samplesPerPixel)

	if samplesPerPixel < 1 {
		return errors.New("Enable Deep Color Mode (-D) for this WAV format")
	}

	if format.NumChannels > 2 {
		return errors.New("WAV files with more than 2 channels are not supported")
	} else if format.BitsPerSample < 8 || format.BitsPerSample > 24 || format.BitsPerSample%8 != 0 {
		return errors.New("Only WAV files with 8, 16, or 24 bits per sample are supported")
	} else if format.BitsPerSample == 16 && format.NumChannels == 2 {
		return errors.New("16-bit stereo WAV files are not supported in any mode; please change the bit depth to either 8 or 24 bit or downmix the file to mono first")
	} else if format.BitsPerSample == 16 && !useDeepColor {
		return errors.New("Use Deep Color Mode (-D) for 16-bit WAV files")
	} else if format.NumChannels == 2 && !useDeepColor {
		return errors.New("Use Deep Color Mode (-D) for stereo WAV files")
	}

	rect := image.Rect(0, 0, 2048, 2048)

	var img image.RGBA
	var img64 image.RGBA64

	if useDeepColor {
		img64 = *image.NewRGBA64(rect)
	} else {
		img = *image.NewRGBA(rect)
	}

	spiral := curves.NewSpiral(e.Diameter, e.Separation)
	spiral.Center = curves.IntegralPoint{X: 1024, Y: 1024}

	var pixelInt uint64
	var shiftedBy uint16
	var isLastIteration = false

	fmt.Printf("Samples per pixel: %d\n", samplesPerPixel)

	for {
		// 48 is the LCM of 8, 16, 24, and 48
		samples, err := reader.ReadSamples(samplesPerPixel)
		if err == io.EOF {
			isLastIteration = true
			break
		}

		for sIdx := 0; sIdx < len(samples); sIdx++ {
			sample := samples[sIdx]

			pixelInt <<= format.BitsPerSample
			shiftedBy += format.BitsPerSample
			pixelInt |= uint64(sample.Values[0])

			if format.NumChannels == 2 {
				pixelInt <<= format.BitsPerSample
				shiftedBy += format.BitsPerSample
				pixelInt |= uint64(sample.Values[0])
			}

			totalSamples++

			if uint32(shiftedBy) == colorDepth || isLastIteration {
				p := spiral.Next()
				if useDeepColor {
					r := uint16(pixelInt & 0xffff00000000 >> 32)
					g := uint16(pixelInt & 0x0000ffff0000 >> 16)
					b := uint16(pixelInt & 0x00000000ffff)

					col := color.RGBA64{r, g, b, 0xffff}
					img64.Set(p.X, p.Y, col)
				} else {
					r := uint8(pixelInt & 0xff0000 >> 16)
					g := uint8(pixelInt & 0x00ff00 >> 8)
					b := uint8(pixelInt & 0x0000ff)

					col := color.RGBA{r, g, b, 0xff}
					img.Set(p.X, p.Y, col)
				}

				pixelInt = 0
				shiftedBy = 0
			}
		}
	}

	fmt.Printf("%d samples written\n", totalSamples)

	radRounded := int(math.Ceil(spiral.Radius))
	subRect := image.Rect(spiral.Center.X-radRounded, spiral.Center.Y-radRounded, spiral.Center.X+radRounded, spiral.Center.Y+radRounded)

	fmt.Println(radRounded)
	fmt.Println(subRect)

	// return true

	f, err := os.Create(e.OutFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if useDeepColor {
		png.Encode(f, img64.SubImage(subRect))
	} else {
		png.Encode(f, img.SubImage(subRect))
	}

	return nil
}
