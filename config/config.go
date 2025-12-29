package config

import (
	"fgengine/constants"
	"fgengine/language"
	"fgengine/types"
)

// Only user-configurable settings should be here
var (
	WindowHeight, WindowWidth, LayoutSizeW, LayoutSizeH int // TODO, remove Layout variables, either move them to constants or remove them entirely and use only Window variables
	ControllerDeadzone                                  float64
	Language                                            language.Lang
)

// WindowRect returns a rect representing the current window dimensions
func WindowRect() types.Rect {
	return types.Rect{X: 0, Y: 0, W: float64(WindowWidth), H: float64(WindowHeight)}
}

// LayoutRect returns a rect representing the current layout dimensions
func LayoutRect() types.Rect {
	return types.Rect{X: 0, Y: 0, W: float64(LayoutSizeW), H: float64(LayoutSizeH)}
}

func SetEditorConfig() {
	initDefaultConfig()
	LayoutSizeW = WindowWidth
	LayoutSizeH = WindowHeight

}

func InitGameConfig() {
	initDefaultConfig()
	LayoutSizeW = int(constants.Camera.W)
	LayoutSizeH = int(constants.Camera.H)

}

func initDefaultConfig() {
	WindowWidth = 1600
	WindowHeight = 900
	ControllerDeadzone = 0.3
	Language = "EN"
}
