package utils

import (
	"bufio"
	"encoding/binary"
	"golang.org/x/exp/slices"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var supportedFormats = []string{"jpeg", "png"}

func FileIsWave(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	fileSize := fileInfo.Size()

	f, err := os.Open(filePath)
	if err != nil {
		return false
	}

	riff := make([]byte, 4)
	size := make([]byte, 4)
	wave := make([]byte, 4)

	var n int
	reader := bufio.NewReader(f)
	n, err = reader.Read(riff)
	if err != nil || n != 4 {
		return false
	}
	n, err = reader.Read(size)
	if err != nil || n != 4 {
		return false
	}
	n, err = reader.Read(wave)
	if err != nil || n != 4 {
		return false
	}

	if string(riff) != "RIFF" || int64(binary.LittleEndian.Uint32(size)) != (fileSize-8) || string(wave) != "WAVE" {
		return false
	}

	return true
}

func FileIsImage(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	_, format, err := image.DecodeConfig(f)
	if err != nil {
		return false
	}

	return slices.Contains(supportedFormats, format)
}
