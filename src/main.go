package main

import (
	"bytes"
	zxtape "github.com/raydac/zxtap-zxtape"
	"io"
)

const __AUTHOR__ = "Igor Maznitsa (http://www.igormaznitsa.com)"
const __VERSION__ = "1.0.4"
const __PROJECTURI__ = "https://github.com/raydac/zxtap-to-wav"

func ParseTap(tapReader io.Reader) ([]*zxtape.TapeBlock, error) {
	var result []*zxtape.TapeBlock

	for {
		block, err := zxtape.ReadTapeBlock(tapReader)
		if err == nil {
			if block != nil {
				result = append(result, block)
			}
		} else {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}

	return result, nil
}

func BytesToReadCloser(data []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(data))
}

func PrepareWav(tape []*zxtape.TapeBlock, amplify bool, gapBetweenFiles int, silenceOnStart bool, freq int, consumer *func(string)) ([]byte, error) {
	var soundBuffer bytes.Buffer

	if consumer != nil && *consumer != nil {
		(*consumer)("Detected data blocks : ")
	}

	for index, tape := range tape {
		if index > 0 || silenceOnStart {
			for i := 0; i < gapBetweenFiles; i++ {
				if consumer != nil && *consumer != nil {
					(*consumer)(".")
				}
			}
			for i := 0; i < freq*gapBetweenFiles; i++ {
				soundBuffer.WriteByte(0x80)
			}
		}

		err := tape.SaveSoundData(amplify, &soundBuffer, freq)
		if err != nil {
			return nil, err
		}

		var label string

		if (*tape.Data)[0] < 128 {
			if len(*tape.Data) == 18 {
				switch (*tape.Data)[1] {
				case 0:
					label = "P"
				case 1:
					label = "N"
				case 2:
					label = "A"
				case 3:
					{
						if (*tape.Data)[12] == 0x00 && (*tape.Data)[13] == 0x1B && (*tape.Data)[14] == 0x00 && (*tape.Data)[15] == 0x40 {
							label = "$"
						} else {
							label = "C"
						}

					}
				default:
					label = "X"
				}
			} else {
				label = "u"
			}
		} else {
			label = "D"
		}

		if consumer != nil && *consumer != nil {
			(*consumer)(label)
		}
	}

	if consumer != nil && *consumer != nil {
		(*consumer)("\n")
	}

	return soundBuffer.Bytes(), nil
}
