package main

import (
	"github.com/jessevdk/go-flags"
	"log"
	"musicimage/encoders"
	"os"
	"path/filepath"
)

func main() {
	opts := encoders.SharedOptions{}
	optsParser := flags.NewParser(&opts, flags.Default)

	_, err := optsParser.Parse()
	if err != nil {
		os.Exit(1)
	}

	isEncode := opts.IsEncode()
	isDecode := opts.IsDecode()

	if !isEncode && !isDecode {
		log.Fatalln("input and output must be audio → image or image → audio")
	}

	if isEncode {
		// encode audio to image
		enc := encoders.Encoder{SharedOptions: opts}
		if err := enc.Encode(); err != nil {
			log.Fatalln(err)
		}
		done(opts)
		return
	}

	if isDecode {
		// decode image to audio
		dec := encoders.Decoder{SharedOptions: opts}
		if err := dec.Decode(); err != nil {
			log.Fatalln(err)
		}
		done(opts)
		return
	}
}

func done(opts encoders.SharedOptions) {
	absPath, err := filepath.Abs(opts.Positional.OutFile)
	if err != nil {
		panic(err)
	}
	log.Printf("Done: %s\n", absPath)
}
