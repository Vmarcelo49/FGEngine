//go:build !js || !wasm
// +build !js !wasm

package filepicker

import (
	"fmt"
	"strings"

	"github.com/sqweek/dialog"
)

type DesktopFilePicker struct{}

func newPlatformFilePicker() FilePicker {
	return &DesktopFilePicker{}
}

func (d *DesktopFilePicker) LoadFile(filter FileFilter) (string, error) {
	dialogBuilder := dialog.File()

	if filter.Description != "" && len(filter.Extensions) > 0 {
		extensionsStr := strings.Join(filter.Extensions, ", ")
		dialogBuilder = dialogBuilder.Filter(filter.Description, extensionsStr)
	}

	path, err := dialogBuilder.Load()
	if err != nil {
		return "", fmt.Errorf("file selection cancelled or failed: %w", err)
	}

	return path, nil
}

func (d *DesktopFilePicker) LoadFiles(filter FileFilter) ([]string, error) {
	// TODO: Use native multi-select dialog when available in dialog package
	// For now, call LoadFile once and return as slice
	path, err := d.LoadFile(filter)
	if err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func (d *DesktopFilePicker) SaveFile(filter FileFilter) (string, error) {
	dialogBuilder := dialog.File()

	if filter.Description != "" && len(filter.Extensions) > 0 {
		extensionsStr := strings.Join(filter.Extensions, ", ")
		dialogBuilder = dialogBuilder.Filter(filter.Description, extensionsStr)
	}

	path, err := dialogBuilder.Save()
	if err != nil {
		return "", fmt.Errorf("file save cancelled or failed: %w", err)
	}

	return path, nil
}
