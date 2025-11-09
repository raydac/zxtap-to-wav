//go:build !(js && wasm)

package main

import (
	"bytes"
	"flag"
	"fmt"
	wav "github.com/raydac/zxtap-wav"
	zxtape "github.com/raydac/zxtap-zxtape"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var fileInName string
var fileOutName string
var amplify bool
var gapBetweenFiles int
var silenceOnStart bool
var freq int

func init() {
	flag.StringVar(&fileInName, "i", "", "source TAP file")
	flag.StringVar(&fileOutName, "o", "", "target WAV file")
	flag.BoolVar(&amplify, "a", false, "amplify sound signal")
	flag.BoolVar(&silenceOnStart, "s", false, "add silence before the first file")
	flag.IntVar(&gapBetweenFiles, "g", 1, "time gap between sound blocks, in seconds")
	flag.IntVar(&freq, "f", 22050, "frequency of result wav, in Hz")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func header() {
	fmt.Printf(`
  TAP2WAV converter of .TAP files (a format for ZX-Spectrum emulators) into its .WAV image (PCM, mono).
  Project page : %s
        Author : %s
       Version : %s

`, __PROJECTURI__, __AUTHOR__, __VERSION__)
}

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

func loadTapFile(filePath string) ([]*zxtape.TapeBlock, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseTap(file)
}

func saveWav(tape []*zxtape.TapeBlock, filePath string, freq int) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var soundBuffer bytes.Buffer

	fmt.Print("Detected data blocks : ")

	for index, tape := range tape {
		if index > 0 || silenceOnStart {
			for i := 0; i < gapBetweenFiles; i++ {
				fmt.Print(".")
			}
			for i := 0; i < freq*gapBetweenFiles; i++ {
				soundBuffer.WriteByte(0x80)
			}
		}

		err = tape.SaveSoundData(amplify, &soundBuffer, freq)
		if err != nil {
			return err
		}

		if (*tape.Data)[0] < 128 {
			if len(*tape.Data) == 18 {
				switch (*tape.Data)[1] {
				case 0:
					fmt.Print("P")
				case 1:
					fmt.Print("N")
				case 2:
					fmt.Print("A")
				case 3:
					{
						if (*tape.Data)[12] == 0x00 && (*tape.Data)[13] == 0x1B && (*tape.Data)[14] == 0x00 && (*tape.Data)[15] == 0x40 {
							fmt.Print("$")
						} else {
							fmt.Print("C")
						}

					}
				default:
					fmt.Print("X")
				}
			} else {
				fmt.Print("u")
			}
		} else {
			fmt.Print("D")
		}
	}

	fmt.Print("\n")

	var data []byte = soundBuffer.Bytes()
	return wav.WriteWav(file, freq, &data)
}

func extractName(filename string) string {
	basename := path.Base(filename)
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

func sizeToHuman(path string) string {
	fi, err := os.Stat(path)
	if err != nil {
		return "<unknown>"
	}

	length := fi.Size()

	return strconv.FormatInt(length/1024, 10) + " Kb"
}

func main() {
	header()

	flag.Parse()

	if freq < 11025 {
		log.Fatal("Unexpected WAV frequency, must be >= 11025 Hz")
	}

	if len(fileInName) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(fileOutName) == 0 {
		fileOutName = filepath.Dir(fileInName) + string(os.PathSeparator) + extractName(fileInName) + ".wav"
	}

	parsedTape, err := loadTapFile(fileInName)
	if err != nil {
		log.Fatal(err)
	}

	err = saveWav(parsedTape, fileOutName, freq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The Result WAV file size is %s\n", sizeToHuman(fileOutName))
}
