package editor

import (
	_ "image/png"
)

type uiVariables struct {
	frameDataIndex          int
	animationSelectionIndex int
	playingAnim             bool
	logBuf                  string
	logUpdated              bool
	logSubmitBuf            string
	boxActionIndex          int
}
