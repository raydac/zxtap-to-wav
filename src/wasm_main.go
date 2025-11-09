//go:build js || wasm

package main

import (
	"syscall/js"
)

func ProcessFile(this js.Value, args []js.Value) interface{} {
	return nil
}

func main() {
	js.Global().Set("processFile", js.FuncOf(ProcessFile))
	println("WASM Go runtime started")
	select {}
}
