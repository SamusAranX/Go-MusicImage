package encoders

import (
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/youpy/go-wav"
	"musicimage/encoders/curves"
	"musicimage/encoders/header"
)

type Encoder struct {
	SharedOptions
}

func makeRGBChunksFromSamples(samples []wav.Sample, format *wav.WavFormat, bitsPerPixel int) [][3]uint16 {
	// pre-flight check
	if !slices.Contains([]int{24, 48}, bitsPerPixel) {
		log.Fatalf("Unsupported number of bits per pixel %d!", bitsPerPixel)
	}

	// prepare format vars
	numChannels := int(format.NumChannels)
	bitsPerSample := int(format.BitsPerSample)
	bytesPerChannel := bitsPerSample / 8
	fmt.Println("bytes per channel:", bytesPerChannel)

	// create array with just the right capacity to hold all remapped RGB triplets
	totalAmountOfBits := float64(len(samples) * numChannels * bitsPerSample)
	retCapacity := int(math.Ceil(totalAmountOfBits / float64(bitsPerPixel)))
	retValues := make([][3]uint16, retCapacity)

	// create a counter that is incremented every time a triplet is added to retValues
	// this counter marks what retValues index will be populated next
	retCounter := 0

	// do the same for the rgbTriplet array
	// rgbTripletCounter is incremented every time after a value is inserted into rgbTriplet
	// once it reaches 3, rgbTriplet is added to retValues and the counter is reset to 0
	rgbTriplet := [3]uint16{}
	rgbTripletCounter := 0

	// stage 1: allocate a byte array big enough to hold all sample bytes
	// also create a counter to keep track of the array field to assign to
	allSampleBytesNum := len(samples) * numChannels * bytesPerChannel
	allSampleBytesIdx := 0
	allSampleBytes := make([]byte, allSampleBytesNum)

	// then loop through all the samples and write their little-endian
	// representations into allSampleBytes
	for _, sample := range samples {
		for chanIdx := 0; chanIdx < numChannels; chanIdx++ {
			tempBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(tempBytes, uint32(sample.Values[chanIdx]))

			for tempBytesIdx := 4 - bytesPerChannel; tempBytesIdx < 4; tempBytesIdx++ {
				allSampleBytes[allSampleBytesIdx] = tempBytes[tempBytesIdx]
				allSampleBytesIdx++
			}
		}
	}

	// stage 2: allocate another array to hold RGB triplets as [3]byte items
	tripletsNum := int(math.Ceil(float64(allSampleBytesNum) / 3.0))
	tripletsIdx := 0
	triplets := make([][3]byte, tripletsNum)

	const chunkSize int = 3
	for i := 0; i < allSampleBytesNum; i += chunkSize {
		end := i + chunkSize

		chunk := [chunkSize]byte{}
		if end > allSampleBytesNum {

		}

		triplets[tripletsIdx] = [chunkSize]byte{
			allSampleBytes[i],
			allSampleBytes[i+1],
			allSampleBytes[i+2],
		}
		tripletsIdx++
	}

	return retValues
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

	insertHeaderPixel := true

	if !header.IsGoodSampleRate(format.SampleRate) {
		fmt.Printf("Note: The input sample rate of %d Hz means MusicImage can't insert a header pixel into the output image.\n", format.SampleRate)
		insertHeaderPixel = false
	}

	if !header.IsGoodGrooveSeparation(uint32(e.GrooveSeparation)) {
		fmt.Printf("Note: The requested groove separation of %d px means MusicImage can't insert a header pixel into the output image.\n", e.GrooveSeparation)
		insertHeaderPixel = false
	}

	if !header.IsGoodLabelDiameter(uint32(e.LabelDiameter)) {
		fmt.Printf("Note: The requested label diameter of %d px means MusicImage can't insert a header pixel into the output image.\n", e.GrooveSeparation)
		insertHeaderPixel = false
	}

	if !insertHeaderPixel {
		fmt.Println("You will have to specify all necessary options yourself when decoding the image back into an audio file.")
	}

	// color depth in bits per pixel, ignoring alpha channel
	var colorDepth uint32 = 24
	useDeepColor := e.DeepColor
	if useDeepColor {
		colorDepth *= 2
	}

	// how many samples are going to be used per pixel (alpha is always 100%)
	samplesPerPixel := colorDepth / uint32(format.BitsPerSample) / uint32(format.NumChannels)

	fmt.Println(colorDepth, format.BitsPerSample, format.NumChannels)
	//fmt.Println(colorDepth, uint32(format.BitsPerSample), uint32(format.NumChannels))
	//fmt.Println(samplesPerPixel)

	if samplesPerPixel < 1 {
		return errors.New("enable Deep Color Mode (-D) for this WAV format")
	}

	if format.NumChannels == 0 || format.NumChannels > 2 {
		return errors.New("WAV files must have one or two channels")
	} else if !header.IsGoodBitDepth(uint32(format.BitsPerSample)) || format.AudioFormat == wav.AudioFormatIEEEFloat {
		return errors.New("floating-point WAV files are not supported")
	}
	// else if format.BitsPerSample == 16 && format.NumChannels == 2 {
	// 	return errors.New("16-bit stereo WAV files are not supported in any mode; please change the bit depth to either 8 or 24 bit or downmix the file to mono first")
	// } else if format.BitsPerSample == 16 && !useDeepColor {
	// 	return errors.New("use Deep Color Mode (-D) for 16-bit WAV files")
	// } else if format.NumChannels == 2 && !useDeepColor {
	// 	return errors.New("use Deep Color Mode (-D) for stereo WAV files")
	// }

	rectSize := 4096
	rect := image.Rect(0, 0, rectSize, rectSize)

	var img image.RGBA
	var img64 image.RGBA64

	if useDeepColor {
		img64 = *image.NewRGBA64(rect)
	} else {
		img = *image.NewRGBA(rect)
	}

	spiral := curves.NewSpiral(e.LabelDiameter, e.GrooveSeparation)
	spiral.Center = curves.IntegralPoint{X: rectSize / 2, Y: rectSize / 2}

	//var pixelInt uint64
	//var shiftedBy uint16
	//var isLastIteration = false
	var totalSamples uint64

	// 192 is the LCM of 8, 16, 24, 32, 48, 64
	const samplesPerLoop = 192

	//fmt.Printf("Samples per pixel: %d\n", samplesPerPixel)

	for {
		// 48 is the LCM of 8, 16, 24, and 48

		samples, _ := reader.ReadSamples(samplesPerLoop)
		//if err == io.EOF {
		//isLastIteration = true
		//	break
		//}

		pixels := makeRGB24ChunksFromSamples(samples, format)
		totalSamples += uint64(len(samples)) * uint64(format.NumChannels)

		for _, pix := range pixels {
			p := spiral.NextIntegral()

			r, g, b := pix[0], pix[1], pix[2]
			col := color.RGBA{R: r, G: g, B: b, A: 0xff}
			img.Set(p.X, p.Y, col)
		}
	}

	fmt.Printf("%d samples written\n", totalSamples)

	radRounded := int(math.Ceil(spiral.Radius()))
	subRect := image.Rect(spiral.Center.X-radRounded, spiral.Center.Y-radRounded, spiral.Center.X+radRounded, spiral.Center.Y+radRounded)

	fmt.Println(radRounded)
	fmt.Println(subRect)

	// return true

	f, err := os.Create(e.Positional.OutFile)
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
