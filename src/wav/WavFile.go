package wav

import (
	"encoding/binary"
	"io"
)

type WavChunkHeader struct {
	ID   [4]uint8
	Size uint32
}

type WavFormat struct {
	ID            [4]uint8
	Size          uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

func writeHeader(writer io.Writer, length uint32) error {
	header := WavChunkHeader{}
	copy(header.ID[:], []uint8("RIFF")[0:4])
	header.Size = length

	if err := binary.Write(writer, binary.LittleEndian, header); err != nil {
		return err
	}

	return binary.Write(writer, binary.LittleEndian, []uint8("WAVE"))
}

func writeObj(writer io.Writer, obj interface{}) error {
	if err := binary.Write(writer, binary.LittleEndian, obj); err != nil {
		return err
	}
	return nil
}

func WriteWav(writer io.Writer, freq int, sndData *[]byte) error {
	if err := writeHeader(writer, uint32(36+len(*sndData))); err != nil {
		return err
	}

	wavFormat := WavFormat{}
	copy(wavFormat.ID[:], []uint8("fmt ")[0:4])
	wavFormat.Size = 16
	wavFormat.AudioFormat = 1
	wavFormat.NumChannels = 1
	wavFormat.SampleRate = uint32(freq)
	wavFormat.ByteRate = uint32(freq)
	wavFormat.BlockAlign = 1
	wavFormat.BitsPerSample = 8

	wavData := WavChunkHeader{}
	copy(wavData.ID[:], []uint8("data")[0:4])
	wavData.Size = uint32(len(*sndData))

	if err := writeObj(writer, wavFormat); err != nil {
		return err
	}

	if err := writeObj(writer, wavData); err != nil {
		return err
	}

	if _, err := writer.Write(*sndData); err != nil {
		return err
	}

	return nil
}
