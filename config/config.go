package config

import (
	"fgengine/constants"
	"fgengine/types"
)

// Only user-configurable settings should be here
var (
	WindowHeight, WindowWidth, LayoutSizeW, LayoutSizeH int
	ControllerDeadzone                                  float64
)

// GetWindowRect returns a rect representing the current window dimensions
func GetWindowRect() types.Rect {
	return types.Rect{X: 0, Y: 0, W: float64(WindowWidth), H: float64(WindowHeight)}
}

// GetLayoutRect returns a rect representing the current layout dimensions
func GetLayoutRect() types.Rect {
	return types.Rect{X: 0, Y: 0, W: float64(LayoutSizeW), H: float64(LayoutSizeH)}
}

func SetEditorConfig() {
	WindowWidth = 1600
	WindowHeight = 900
	LayoutSizeW = WindowWidth
	LayoutSizeH = WindowHeight
	ControllerDeadzone = 0.3
}

func InitDefaultConfig() {
	WindowWidth = 1600
	WindowHeight = 900
	LayoutSizeW = int(constants.Camera.W)
	LayoutSizeH = int(constants.Camera.H)
	ControllerDeadzone = 0.3
}
