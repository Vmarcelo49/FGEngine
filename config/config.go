package config

import (
	"fgengine/language"

	"github.com/hajimehoshi/ebiten/v2"
)

// Only user-configurable settings should be here
var (
	WindowHeight, WindowWidth int
	ControllerDeadzone        float64
	Language                  language.Lang
)

func InitGameConfig() {
	initDefaultConfig()
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("FG Engine")

}

func initDefaultConfig() {
	WindowWidth = 1600
	WindowHeight = 900
	ControllerDeadzone = 0.3
	Language = "EN"
}
