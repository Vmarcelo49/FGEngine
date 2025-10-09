//go:build js && wasm
// +build js,wasm

package filepicker

import (
	"fmt"
)

// WebAssemblyFilePicker placeholder implementation for future HTML integration
type WebAssemblyFilePicker struct{}

func newPlatformFilePicker() FilePicker {
	return &WebAssemblyFilePicker{}
}

// LoadFile TODO: implement using syscall/js and HTML File API
func (w *WebAssemblyFilePicker) LoadFile(filter FileFilter) (string, error) {
	return "", fmt.Errorf("WebAssembly file picker not implemented")
}

// SaveFile TODO: implement using syscall/js and File System Access API
func (w *WebAssemblyFilePicker) SaveFile(filter FileFilter) (string, error) {
	return "", fmt.Errorf("WebAssembly file save not implemented")
}
