# MusicImage (but written in Go this time)
Creates fancy vinyl-like PNGs from music files

## Features
* Can "encode" mono/stereo 8-24 bit WAV files to PNG files
	* (with the exception of 16-bit stereo files)
* Can "decode" PNG files to re-create the original WAV
* Support for "deep color" PNG files to fit more samples into a single pixel

## Download
There are no releases, you'll have to build MusicImage yourself. If your Go development environment is set up properly, this should be as easy as running:

```
go get
go build
```

## Usage

To create PNGs from WAV files:

```
$ musicimage -i infile-11025hz.wav -o encoded.png
```

To recreate the WAV file from a PNG:

```
$ musicimage -i encoded.png -o outfile-decoded.wav
```

### General Options

* `-d`/`--diameter`: The diameter of the cutout in the middle of the image. Can be used to insert album covers in a vinyl label-like fashion. *(Optional. Default: **64**)*
* `-s`/`--separation`: The distance between spiral turns, in pixels. *(Optional. Default: **1**)*
* `-D`/`--deep`: Treats PNGs as 64-bit PNGs to stuff more data into or read more data from a single pixel. *(Optional. Default: **1**)*

### PNG → WAV Options

* `-r`/`--bitrate`: Sample rate override. *(Optional. Default: **11025**)*
* `-b`/`--bitdepth`: Bit depth override. *(Optional. Allowed values: 8, 16, 24. Default: **8**)*
* `-c`/`--channels`: Channel number override. *(Optional. Allowed values: 1, 2. Default: **1**)*

## Example
![MusicImage-encoded excerpt from Ultra Sheriff's "Leviathan"](https://i.peterwunder.de/leviathan.png)

This image decodes to [this WAV file](https://i.peterwunder.de/leviathan.wav) (with the standard settings: `-d 64 -s 1 -r 11025 -b 8 -c 1`)

(The full song is available [on iTunes](https://itunes.apple.com/us/album/deception-oil-and-laser-beams-ep/1105412287) and [on Spotify](https://open.spotify.com/track/4NRyBYL1pyMX696XcRgeWw), by the way. It's great.)

## Limitations
MusicImage only supports WAV files with:

* one or two channels
* 8, 16, or 24 bits per sample
	* 16-bit stereo WAV files are not supported at all

## Feedback
Just tweet at me [@SamusAranX](https://twitter.com/SamusAranX).
Feel free to file an issue if you encounter any crashes, bugs, etc.: https://github.com/SamusAranX/Go-MusicImage/issues