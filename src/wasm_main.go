//go:build js || wasm

package main

import (
	wav "github.com/raydac/zxtap-wav"
	"strconv"
	"syscall/js"
	"bytes"
)

func ProcessTapArray(this js.Value, args []js.Value) interface{} {
	if len(args) < 5 {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "missing arguments",
		})
	}

	src := args[0]

	if src.Type() != js.TypeObject {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "first argument must be Uint8Array",
		})
	}

	consoleOutput := func(s string) {
		js.Global().Get("console").Call("log", s)
	}

	length := src.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, src)

	consoleOutput("Received " + strconv.Itoa(len(data)) + " bytes")

	amplify := args[1].Bool()
	gapBetweenFiles := args[2].Int()
	if gapBetweenFiles < 0 {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "gap between files must be non-negative",
		})
	}
	silenceOnStart := args[3].Bool()
	freq := args[4].Int()
	if freq < 0 {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "frequency must be non-negative",
		})
	}

	parsedTape, err := ParseTap(BytesToReader(data),  &consoleOutput)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "error during parse tap data: " + err.Error(),
		})
	}
	consoleOutput("Tap parsed to " + strconv.Itoa(len(parsedTape)) + " part(s)")
	consoleOutput("Making WAV " + strconv.Itoa(freq) + "Hz")

	var wavData, errPrepare = PrepareWav(parsedTape, amplify, gapBetweenFiles, silenceOnStart, freq, &consoleOutput)

	if errPrepare != nil {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "error during prepare wav :" + errPrepare.Error(),
		})
	}

	var buf bytes.Buffer
	writeWavError := wav.WriteWav(&buf, freq, &wavData)
	if writeWavError != nil {
		return js.ValueOf(map[string]interface{}{
			"data":  nil,
			"error": "error during form wav data : " + writeWavError.Error(),
		})
	}
	resultWavData := buf.Bytes()

	consoleOutput("Generated WAV data " + strconv.Itoa(len(resultWavData)) + " bytes")

	uint8Array := js.Global().Get("Uint8Array").New(len(resultWavData))
	js.CopyBytesToJS(uint8Array, resultWavData)

	consoleOutput("Converted into Uint8Array " + strconv.Itoa(uint8Array.Get("length").Int()) + " bytes")

	return js.ValueOf(map[string]interface{}{
		"data":  uint8Array,
		"error": nil,
	})
}

func main() {
	js.Global().Set("ProcessTapArray", js.FuncOf(ProcessTapArray))
	println("WASM Go runtime started")
	select {}
}
