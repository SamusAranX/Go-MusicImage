package musicoders

type TestEncoderOpts struct {
	EncoderOpts

	Turns float64 `short:"T" long:"turns" default:"100" description:"Only used for debugging"`
}

type EncoderOpts struct {
	Mode string `short:"m" long:"mode" default:"spiral" description:"Unused"`
}

type DecoderOpts struct {
	SampleRate uint32 `short:"r" long:"rate" default:"11025" description:"Sample rate of destination audio file, in Hz"`
	ChannelNum uint16 `short:"c" long:"channels" default:"1" description:""`
	BitDepth   uint16 `short:"b" long:"bitdepth" default:"8" description:""`
}

type SharedOptions struct {
	InFile  string `short:"i" long:"infile" required:"yes" description:"Input file"`
	OutFile string `short:"o" long:"outfile" required:"yes" description:"Output file"`

	Diameter   uint32  `short:"d" long:"diameter" default:"64" description:"Diameter of vinyl label, in pixels"`
	Separation float64 `short:"s" long:"separation" default:"1" description:"Distance between spiral turns, in pixels"`
	DeepColor  bool    `short:"D" long:"deep" description:"Use 64-bit pixels to store more data in a single pixel"`

	EncoderOptions     EncoderOpts     `command:"encode" alias:"e" description:"Encode a WAV file into a PNG image"`
	DecoderOptions     DecoderOpts     `command:"decode" alias:"d" description:"Decode an image into a WAV file"`
	TestEncoderOptions TestEncoderOpts `command:"test" description:"Test the drawing code"`
}
