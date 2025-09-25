package config

import (
	"fgengine/constants"
)

// Only user-configurable settings should be here
var (
	WindowHeight, WindowWidth, LayoutSizeW, LayoutSizeH int
	ControllerDeadzone                                  float64
)

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
	LayoutSizeW = int(constants.CameraWidth)
	LayoutSizeH = int(constants.CameraHeight)
	ControllerDeadzone = 0.3
}
