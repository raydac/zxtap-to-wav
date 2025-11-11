package zx

import (
	"io"
	"encoding/binary"
)

func ReadZxShort(in io.Reader) (int, error) {
    var value int16
	if err := binary.Read(in, binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return int(value), nil
}

func ReadZxArray(in io.Reader, n int) ([]byte, error) {
    buf := make([]byte, n)
	_, err := io.ReadFull(in, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ReadZxByte(in io.Reader) (byte, error) {
    var b [1]byte
	_, err := io.ReadFull(in, b[:])
	return b[0], err
}
