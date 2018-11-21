package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	. "musicimage/musicoders"
	. "strings"
)

func main() {
	var opts struct {
		Args struct {
			InFile  string
			OutFile string
		} `positional-args:"yes" required:"yes"`

		Diameter   uint32  `short:"d" long:"diameter" default:"64" description:"Diameter of vinyl label. in pixels"`
		Separation float64 `short:"s" long:"separation" default:"1" description:"Distance between spiral turns, in pixels"`
		SampleRate uint32  `short:"r" long:"rate" default:"11025" description:"Sample rate of destination audio file, in Hz"`
		Turns      float64 `short:"t" long:"turns" default:"100" description:"Only used for debugging."`
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

	infileIsImage := HasSuffix(infile, "png") || HasSuffix(infile, "jpg") || HasSuffix(infile, "jpeg") || HasSuffix(infile, "gif")
	infileIsSound := HasSuffix(infile, "wav")
	outfileIsImage := HasSuffix(outfile, "png") || HasSuffix(outfile, "jpg") || HasSuffix(outfile, "jpeg") || HasSuffix(outfile, "gif")
	outfileIsSound := HasSuffix(outfile, "wav")

	if infileIsSound && outfileIsImage {
		fmt.Println("encode wav to png")

		enc := Encoder{infile, opts.Diameter, opts.Separation}

		if enc.Encode(opts.Args.OutFile) {
			fmt.Println("Encoded successfully!")
		} else {
			fmt.Println("Encoding failed")
		}
	} else if infileIsImage && outfileIsSound {
		fmt.Println("decode png to wav")

		dec := Decoder{opts.Args.InFile, opts.Diameter, opts.Separation, 0, 0, 0}
		dec.ChannelNum = 1
		dec.SampleBits = 24
		dec.SampleRate = 11025

		if dec.Decode(opts.Args.OutFile) {
			fmt.Println("Decoded successfully!")
		} else {
			fmt.Println("Decoding failed")
		}
	} else if infile == "test" {
		test := TestEncoder{opts.Turns, opts.Diameter, opts.Separation}
		test.Encode(outfile)
	} else {
		fmt.Printf("inImage: %t\ninSound: %t\noutImage: %t\noutSound: %t\n", infileIsImage, infileIsSound, outfileIsImage, outfileIsSound)
	}
}
