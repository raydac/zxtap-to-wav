//go:build js || wasm

package main

import (
	"syscall/js"
)

func ProcessTapArray(this js.Value, args []js.Value) any {
	if len(args) < 5 {
		return js.ValueOf(map[string]any{
			"data":  nil,
			"error": "missing arguments",
		})
	}

	src := args[0]

	if src.Type() != js.TypeObject {
		return js.ValueOf(map[string]any{
			"data":  nil,
			"error": "first argument must be Uint8Array",
		})
	}

	length := src.Get("length").Int()
	data := make([]byte, length)

	amplify := args[1].Bool()
	gapBetweenFiles := args[2].Int()
	if gapBetweenFiles < 0 {
		return js.ValueOf(map[string]any{
			"data":  nil,
			"error": "gap between files must be non-negative",
		})
	}
	silenceOnStart := args[3].Bool()
	freq := args[4].Int()
	if freq < 0 {
		return js.ValueOf(map[string]any{
			"data":  nil,
			"error": "frequency must be non-negative",
		})
	}

	parsedTape, err := ParseTap(BytesToReadCloser(data))
	if err != nil {
		return js.ValueOf(map[string]any{
			"data":  nil,
			"error": "error during parse tap data",
		})
	}

	consoleOutput := func(s string) {
		js.Global().Get("console").Call("log", s)
	}

	var wavData, errPrepare = PrepareWav(parsedTape, amplify, gapBetweenFiles, silenceOnStart, freq, &consoleOutput)

	if errPrepare != nil {
		return js.ValueOf(map[string]any{
			"data":  nil,
			"error": "error during prepare wav",
		})
	}

	uint8Array := js.Global().Get("Uint8Array").New(len(wavData))
	js.CopyBytesToJS(uint8Array, wavData)

	return js.ValueOf(map[string]any{
		"data":  uint8Array,
		"error": nil,
	})
}

func main() {
	js.Global().Set("processTapArray", js.FuncOf(ProcessTapArray))
	println("WASM Go runtime started")
	select {}
}
