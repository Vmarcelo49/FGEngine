package editor

import (
	"fgengine/collision"
	"fgengine/types"
	_ "image/png"
)

type uiVariables struct {
	frameDataIndex          int
	animationSelectionIndex int
	playingAnim             bool
	logBuf                  string
	logUpdated              bool
	logSubmitBuf            string
	boxDropdownTypeIndex    int

	// Box editor
	enableMouseInput *bool
	activeBoxType    collision.BoxType
	activeBoxIndex   int
	// mouse input related
	dragged           bool
	dragStartMousePos types.Vector2
	dragStartBoxPos   types.Vector2
}
