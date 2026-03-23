//go:build !js && !wasm
// +build !js,!wasm

package main

import "fgengine/editorimgui"

func main() {
	editorimgui.Run()
}
