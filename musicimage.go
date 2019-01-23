package main

import (
	"fmt"
	"github.com/SamusAranX/musicimage/musicoders"
	flags "github.com/jessevdk/go-flags"
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
		} else {
			fmt.Println("Done!")
		}
	case "decode":
		dec := musicoders.Decoder{opts}
		if err := dec.Decode(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Done!")
		}
	case "test":
		testEnc := musicoders.TestEncoder{opts}
		if err := testEnc.Encode(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Done!")
		}
	default:
		fmt.Println("Invalid command")
	}
}
