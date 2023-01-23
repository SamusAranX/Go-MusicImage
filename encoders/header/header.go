package header

import (
	"encoding/binary"
	"golang.org/x/exp/slices"

	"github.com/sigurn/crc8"
)

const (
	bitsSampleRateTableIndex uint32 = 4
	bitsGrooveSeparation     uint32 = 4
	bitsLabelSizeMultiplier  uint32 = 5
	bitsBitDepthTableIndex   uint32 = 2
	bitsUsesStereo           uint32 = 1
	bitsCRC8                 uint32 = 8

	maskSampleRateTableIndex uint32 = 1<<bitsSampleRateTableIndex - 1
	maskGrooveSeparation     uint32 = 1<<bitsGrooveSeparation - 1
	maskLabelSizeMultiplier  uint32 = 1<<bitsLabelSizeMultiplier - 1
	maskBitDepthTableIndex   uint32 = 1<<bitsBitDepthTableIndex - 1
	maskUsesStereo           uint32 = 1<<bitsUsesStereo - 1
	maskCRC8                 uint32 = 1<<bitsCRC8 - 1
)

var (
	knownSampleRates = []uint32{
		4000,
		4096,
		8000,
		8192,
		11025,
		16000,
		16384,
		22050,
		32000,
		32768,
		44100,
		48000,
		88200,
		96000,
		176400,
		192000,
	}

	knownBitDepths = []uint32{
		8,
		16,
		24,
		32,
	}

	fieldBits = []uint32{
		bitsSampleRateTableIndex,
		bitsGrooveSeparation,
		bitsLabelSizeMultiplier,
		bitsBitDepthTableIndex,
		bitsUsesStereo,
		bitsCRC8,
	}

	fieldMasks = []uint32{
		maskSampleRateTableIndex,
		maskGrooveSeparation,
		maskLabelSizeMultiplier,
		maskBitDepthTableIndex,
		maskUsesStereo,
		maskCRC8,
	}
)

type Header struct {
	SampleRateTableIndex uint32
	GrooveSeparation     uint32
	LabelSizeMultiplier  uint32
	BitDepthTableIndex   uint32
	UsesStereo           uint32
	CRC8                 uint32
}

func IsGoodSampleRate(sampleRate uint32) bool {
	return slices.Contains(knownSampleRates, sampleRate)
}

func IsGoodGrooveSeparation(grooveSeparation uint32) bool {
	return grooveSeparation > 0 && grooveSeparation <= 16
}

func IsGoodLabelDiameter(labelDiameter uint32) bool {
	return labelDiameter%16 == 0 && labelDiameter <= 31*16
}

func IsGoodBitDepth(bitDepth uint32) bool {
	return slices.Contains(knownBitDepths, bitDepth)
}

func (h Header) fieldsOrdered() []uint32 {
	return []uint32{
		h.SampleRateTableIndex,
		h.GrooveSeparation,
		h.LabelSizeMultiplier,
		h.BitDepthTableIndex,
		h.UsesStereo,
		h.CRC8,
	}
}

func (h Header) ToBytes(includeCRC bool) []byte {
	fields := h.fieldsOrdered()
	numFields := len(fields)

	if !includeCRC {
		// make the below loop omit the CRC field
		numFields--
	}

	intVal := fields[0]
	for i := 1; i < numFields; i++ {
		intVal <<= fieldBits[i]
		intVal |= fields[i] & fieldMasks[i]
	}

	if includeCRC {
		retBytes := make([]byte, 3)
		binary.BigEndian.PutUint32(retBytes, intVal)
		return retBytes
	} else {
		retBytes := make([]byte, 2)
		binary.BigEndian.PutUint32(retBytes, intVal)
		return retBytes
	}
}

func (h *Header) ParseInt(val uint32) {
	h.CRC8 = val & maskCRC8
	val >>= bitsCRC8

	h.UsesStereo = val & maskUsesStereo
	val >>= bitsUsesStereo

	h.BitDepthTableIndex = val & maskBitDepthTableIndex
	val >>= bitsBitDepthTableIndex

	h.LabelSizeMultiplier = val & maskLabelSizeMultiplier
	val >>= bitsLabelSizeMultiplier

	h.GrooveSeparation = val & maskGrooveSeparation
	val >>= bitsGrooveSeparation

	h.SampleRateTableIndex = val & maskSampleRateTableIndex
}

func (h Header) calculateCRC8() uint8 {
	table := crc8.MakeTable(crc8.CRC8)
	return crc8.Checksum(h.ToBytes(false), table)
}

func (h Header) VerifyCRC8() bool {
	crc := h.calculateCRC8()
	return h.CRC8 == uint32(crc)
}
