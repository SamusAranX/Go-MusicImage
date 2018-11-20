package main

import (
	// "bufio"
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"image"
	"image/png"
	"math"
	"musicimage/hsv"
	"musicimage/micoders"
	"os"
	"strings"
)

func main() {
	var opts struct {
		Args struct {
			InFile  string
			OutFile string
		} `positional-args:"yes" required:"yes"`

		Diameter   uint32  `short:"d" long:"diameter" default:"64" description:"Diameter of vinyl label. in pixels"`
		Separation float64 `short:"s" long:"separation" default:"1" description:"Distance between spiral turns, in pixels"`
		SampleRate uint32  `short:"r" long:"rate" default:"8000" description:"Sample rate of destination audio file, in Hz"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	infile := opts.Args.InFile
	outfile := opts.Args.OutFile

	fmt.Printf("Infile:  %v\n", opts.Args.InFile)
	fmt.Printf("Outfile: %v\n", opts.Args.OutFile)
	fmt.Println("------------")
	fmt.Printf("Diameter:    %v\n", opts.Diameter)
	fmt.Printf("Separation:  %v\n", opts.Separation)
	fmt.Printf("Sample Rate: %v\n", opts.SampleRate)
	fmt.Println("------------")

	if strings.HasSuffix(infile, "wav") && strings.HasSuffix(outfile, "png") {
		fmt.Println("encode wav to png")
	} else if strings.HasSuffix(infile, "png") && strings.HasSuffix(outfile, "wav") {
		fmt.Println("decode png to wav")
	}

	s := micoders.NewSpiral(opts.Diameter, opts.Separation)

	p1 := image.Point{0, 0}
	p2 := image.Point{256, 256}

	r1 := image.Rectangle{p1, p2}

	img := image.NewRGBA(r1)

	var turns float64 = 40

	// var fmtStr = "%.0f turns so far"
	// var progStr string

	var hue uint16

	for s.Theta < (math.Pi * turns * 2) {
		p := s.Next()

		hsvCol := hsv.HSVColor{hue, 255, 255}

		img.Set(p.X, p.Y, hsvCol.RGBA())

		hue = (hue + 2) % 360

		// tmpStr := fmt.Sprintf(fmtStr, s.Theta/math.Pi/2)
		// if tmpStr != progStr {
		// 	fmt.Println(tmpStr)
		// 	progStr = tmpStr
		// }
	}

	f, err := os.Create("/Users/peterwunder/go/src/musicimage/spiraltest.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
	fmt.Println("done")

}
