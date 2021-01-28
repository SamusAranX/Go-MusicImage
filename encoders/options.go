package encoders

type TestEncoderOpts struct {
	EncoderOpts

	Turns float64 `short:"T" long:"turns" default:"128" description:"debugging option"`
}

type EncoderOpts struct {
	Mode string `short:"m" long:"mode" default:"spiral" description:"Currently unused"`
}

type DecoderOpts struct {
	SampleRate uint32 `short:"r" long:"rate" default:"11025" description:"Sample rate of destination audio file, in Hz"`
	ChannelNum uint16 `short:"c" long:"channels" default:"1" description:"Number of channels"`
	BitDepth   uint16 `short:"b" long:"bitdepth" default:"8" description:""`
}

type SharedOptions struct {
	InFile  string `short:"i" long:"infile" required:"yes" description:"Input file"`
	OutFile string `short:"o" long:"outfile" required:"yes" description:"Output file"`

	Diameter   uint32  `short:"d" long:"diameter" default:"64" description:"Diameter of vinyl label, in pixels"`
	Separation float64 `short:"s" long:"separation" default:"1" description:"Distance between spiral turns, in pixels"`
	DeepColor  bool    `short:"D" long:"deep" description:"Use 64 bits per pixel"`

	EncoderOptions     EncoderOpts     `command:"encode" alias:"e" description:"Take a WAV file and create a PNG image"`
	DecoderOptions     DecoderOpts     `command:"decode" alias:"d" description:"Take a PNG image and create a WAV file"`
	TestEncoderOptions TestEncoderOpts `command:"test" description:"debug option"`
}
