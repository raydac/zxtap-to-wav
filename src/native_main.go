//go:build !(js && wasm)

package main

import (
	"flag"
	"fmt"
	wav "github.com/raydac/zxtap-wav"
	zxtape "github.com/raydac/zxtap-zxtape"
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

func loadTapFile(filePath string, consumer *func(string)) ([]*zxtape.TapeBlock, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseTap(file, consumer)
}

func saveWav(tape []*zxtape.TapeBlock, filePath string, freq int, consumer *func(string)) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data, errPrepare = PrepareWav(tape, amplify, gapBetweenFiles, silenceOnStart, freq, consumer)
	if errPrepare != nil {
		return errPrepare
	}

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

	consoleOutputLn := func(s string) {
		fmt.Println(s)
	}

	consoleOutput := func(s string) {
		fmt.Print(s)
	}

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

	parsedTape, err := loadTapFile(fileInName, &consoleOutputLn)
	if err != nil {
		log.Fatal(err)
	}

	err = saveWav(parsedTape, fileOutName, freq, &consoleOutput)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The Result WAV file size is %s\n", sizeToHuman(fileOutName))
}
