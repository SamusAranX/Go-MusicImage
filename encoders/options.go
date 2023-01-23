package encoders

import "musicimage/utils"

type SharedOptions struct {
	// Shared Options
	InFile     string `short:"i" long:"infile" required:"yes" description:"Input file"`
	Positional struct {
		OutFile string
	} `positional-args:"yes" required:"yes" description:"Output file"`

	// Encoder (WAV -> PNG) Options
	LabelDiameter    uint16 `short:"l" long:"label-diameter" default:"64" description:"Diameter of vinyl label, in pixels"`
	GrooveSeparation uint8  `short:"g" long:"groove-separation" default:"1" description:"Distance between vinyl 'grooves', in pixels"`
	DeepColor        bool   `short:"D" long:"deep-color" description:"Use 64 bits per pixel when creating PNG"`

	// Decoder (PNG -> WAV) Options
	SampleRate  uint32 `short:"r" long:"rate" default:"8000" description:"Sample rate in Hz"`
	NumChannels uint8  `short:"c" long:"channels" default:"2" description:"Number of channels (1 or 2)"`
	BitDepth    uint8  `short:"b" long:"bits" default:"8" description:"Bits per sample"`
}

// IsEncode returns true if the arguments given would result in a wave file being turned into an image
func (o SharedOptions) IsEncode() bool {
	return utils.FileIsWave(o.InFile) && utils.FileIsImage(o.Positional.OutFile)
}

// IsDecode returns true if the arguments given would result in an image being turned into a wave file
func (o SharedOptions) IsDecode() bool {
	return utils.FileIsImage(o.InFile) && utils.FileIsWave(o.Positional.OutFile)
}

// CheckArgValidity returns nil if the arguments are valid, otherwise it returns an error with more information
func (o SharedOptions) CheckArgValidity() error {
	return nil
}
