package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"musicimage/musicoders"
	"os"
	// . "strings"
)

func main() {
	opts := musicoders.SharedOptions{}
	optsParser := flags.NewParser(&opts, flags.Default)
	optsParser.SubcommandsOptional = false

	_, err := optsParser.Parse()
	if err != nil {
		os.Exit(-1)
	}

	switch optsParser.Active.Name {
	case "encode":
		enc := musicoders.Encoder{opts}
		if err := enc.Encode(); err != nil {
			fmt.Println(err)
		}
	case "decode":
		dec := musicoders.Decoder{opts}
		dec.Decode()
	default:
		fmt.Println("Invalid command")
	}

	// infile := opts.Args.InFile
	// outfile := opts.Args.OutFile

	// fmt.Printf("Infile:  %v\n", opts.Args.InFile)
	// fmt.Printf("Outfile: %v\n", opts.Args.OutFile)
	// fmt.Println("------------")
	// fmt.Printf("Diameter:    %v\n", opts.Diameter)
	// fmt.Printf("Separation:  %v\n", opts.Separation)
	// fmt.Printf("Sample Rate: %v\n", opts.SampleRate)
	// fmt.Println("------------")

	// infileIsImage := HasSuffix(infile, "png") || HasSuffix(infile, "jpg") || HasSuffix(infile, "jpeg") || HasSuffix(infile, "gif")
	// infileIsSound := HasSuffix(infile, "wav")
	// outfileIsImage := HasSuffix(outfile, "png") || HasSuffix(outfile, "jpg") || HasSuffix(outfile, "jpeg") || HasSuffix(outfile, "gif")
	// outfileIsSound := HasSuffix(outfile, "wav")

	// if infileIsSound && outfileIsImage {
	// 	fmt.Println("encode wav to png")

	// 	enc := musicoders.Encoder{infile, opts.Diameter, opts.Separation, opts.DeepColor}

	// 	if enc.Encode(opts.Args.OutFile) {
	// 		fmt.Println("Encoded successfully!")
	// 	} else {
	// 		fmt.Println("Encoding failed")
	// 	}
	// } else if infileIsImage && outfileIsSound {
	// 	fmt.Println("decode png to wav")

	// 	dec := musicoders.Decoder{opts.Args.InFile, opts.Diameter, opts.Separation, 0, 0, 0}
	// 	dec.ChannelNum = 1
	// 	dec.SampleBits = 24
	// 	dec.SampleRate = 11025

	// 	if dec.Decode(opts.Args.OutFile) {
	// 		fmt.Println("Decoded successfully!")
	// 	} else {
	// 		fmt.Println("Decoding failed")
	// 	}
	// } else if infile == "test" {
	// 	test := TestEncoder{opts.Turns, opts.Diameter, opts.Separation}
	// 	test.Encode(outfile)
	// } else {
	// 	fmt.Printf("inImage: %t\ninSound: %t\noutImage: %t\noutSound: %t\n", infileIsImage, infileIsSound, outfileIsImage, outfileIsSound)
	// }
}
