package zxtape

import (
	"bytes"
	"testing"
)

func TestReadTapeBlock(t *testing.T) {
	block, err := ReadTapeBlock(bytes.NewReader([]byte{7, 0, 1, 2, 3, 4, 5, 6, 99}))

	if err != nil {
		t.Error("Unexpected IO error", err)
	}

	if block == nil {
		t.Error("Block is nil")
	}

	if len(*block.Data) != 6 {
		t.Error("Wrong data length", len(*block.Data))
	}

	if block.Checksum != 99 {
		t.Error("Wrong checksum value", block.Checksum)
	}
}
