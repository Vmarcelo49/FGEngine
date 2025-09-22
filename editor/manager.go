package editor

import (
	"fgengine/animation"

	_ "image/png"
)

type EditorManager struct {
	activeAnimation         *animation.Animation // should point to an animation in the active character
	frameCount              int
	frameIndex              int
	animationSelectionIndex int
	playingAnim             bool
	previousAnimationName   string

	boxEditor *BoxEditor

	// UI related
	logBuf             string
	logUpdated         bool
	logSubmitBuf       string
	choiceShowAllBoxes bool
	boxActionIndex     int
}

func (e *EditorManager) getCurrentSprite() *animation.Sprite {
	if e == nil || e.activeAnimation == nil {
		return nil
	}
	if e.frameIndex < 0 || e.frameIndex >= len(e.activeAnimation.Sprites) {
		return nil
	}
	return e.activeAnimation.Sprites[e.frameIndex]
}

func (e *EditorManager) SetActiveAnimation(anim *animation.Animation) {
	e.activeAnimation = anim
	e.frameIndex = 0
	e.playingAnim = false
	// Clear box editor when switching animations
	e.boxEditor = nil
}
