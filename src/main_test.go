package main

import (
	"bytes"
	"testing"
)

func TestParseTap(t *testing.T) {
	blocks, err := ParseTap(bytes.NewReader([]byte{7, 0, 1, 2, 3, 4, 5, 6, 99, 4, 0, 3, 5, 6, 32}))
	if err != nil {
		t.Error("Unexpected error", err)
	}

	if len(blocks) != 2 {
		t.Error("Unexpected number of blocks", len(blocks))
	}
}
