package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"musicimage/encoders"
)

func main() {
	opts := encoders.SharedOptions{}
	optsParser := flags.NewParser(&opts, flags.Default)
	optsParser.SubcommandsOptional = false

	_, err := optsParser.Parse()
	if err != nil {
		os.Exit(-1)
	}

	switch optsParser.Active.Name {
	case "encode":
		enc := encoders.Encoder{opts}
		if err := enc.Encode(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Done!")
		}
	case "decode":
		dec := encoders.Decoder{opts}
		if err := dec.Decode(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Done!")
		}
	case "test":
		testEnc := encoders.TestEncoder{opts}
		if err := testEnc.Encode(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Done!")
		}
	default:
		fmt.Println("Invalid command")
	}
}
