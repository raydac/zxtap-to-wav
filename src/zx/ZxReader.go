package zx

import (
	"io"
)

func ReadZxShort(in io.Reader) (int, error) {
	arr := make([]byte, 2)
	_, err := io.ReadAtLeast(in, arr, 2)
	if err != nil {
		return 0, err
	}
	return int(arr[1])<<8 | int(arr[0]), nil
}

func ReadZxByte(in io.Reader) (byte, error) {
	arr := make([]byte, 1)
	_, err := io.ReadAtLeast(in, arr, 1)
	if err != nil {
		return 0, err
	}
	return arr[0], nil
}
