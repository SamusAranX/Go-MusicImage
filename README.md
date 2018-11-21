# MusicImage (but written in Go this time)
Creates fancy vinyl-like PNGs from music files

## Features
* Can "encode" 24 bit mono WAV files to PNG files
* Can "decode" those PNG files (and JPEG and GIF files) to re-create the original WAV

## Download
There are no releases, you'll have to build MusicImage yourself. If your Go development environment is set up properly, this should be as easy as running:

```
go get
go build
```

## Usage

To encode sound files in images:

```
$ musicimage infile.wav outfile.png -d 64 -s 0.5
```

To decode those images back to get sound files:

```
$ musicimage infile.png outfile.wav -d 64 -s 0.5 -r 11025
```

* `-d`/`--diameter`: The diameter of the cutout in the middle of the image. Can be used to insert album covers in a vinyl label-like fashion. *(Optional. Default: **64**)*
* `-s`/`--separation`: The distance between spiral turns, in pixels. *(Optional. Default: **1**)*
* `-r`/`--rate`: Sample rate override. Only used when decoding PNGs to WAVs. *(Optional. Default: **11025**)*

## Example
![MusicImage-encoded excerpt from Ultra Sheriff's "Leviathan"](https://i.peterwunder.de/leviathan.png)

This image decodes to [this WAV file](https://i.peterwunder.de/leviathan.wav) (with the standard settings: `-d 64 -s 1 -r 11025`)

(The full song is available [on iTunes](https://itunes.apple.com/us/album/deception-oil-and-laser-beams-ep/1105412287) and [on Spotify](https://open.spotify.com/track/4NRyBYL1pyMX696XcRgeWw), by the way. It's great.)

## Limitations
MusicImage only supports WAV files with:

* one channel
* 24 bits per sample

Other bit rates (and perhaps even support for stereo files) will be implemented at some point.

Attempting to open files with more than one channel will cause MusicImage to only consider the left channel.

## Feedback
Just tweet at me [@SamusAranX](https://twitter.com/SamusAranX).
Feel free to file an issue if you encounter any crashes, bugs, etc.: https://github.com/SamusAranX/Go-MusicImage/issues