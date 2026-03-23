//go:build !js && !wasm
// +build !js,!wasm

package editorimgui

import (
	"fgengine/collision"
	"fgengine/types"
	_ "image/png"
)

type uiVariables struct {
	frameDataIndex       int
	playingAnim          bool
	logBuf               string
	newAnimationName     string
	renameAnimationName  string
	frameDurationInput   string
	boxDropdownTypeIndex int

	enableMouseInput *bool
	activeBoxType    collision.BoxType
	activeBoxIndex   int

	dragged           bool
	dragStartMousePos types.Vector2
	dragStartBoxPos   types.Vector2
}
