package editor

import (
	"fgengine/animation"

	_ "image/png"
)

type EditorManager struct {
	activeAnimation         *animation.Animation // should point to an animation in the active character
	frameCount              int
	frameIndex              int
	frameCounter            int // Counter for animation playback timing
	animationSelectionIndex int
	playingAnim             bool
	previousAnimationName   string

	boxEditor *BoxEditor

	// UI related
	logBuf         string
	logUpdated     bool
	logSubmitBuf   string
	boxActionIndex int
}

func (e *EditorManager) getCurrentSprite() *animation.Sprite {
	if e == nil || e.activeAnimation == nil {
		return nil
	}
	if e.frameIndex < 0 || e.frameIndex >= len(e.activeAnimation.FrameData) {
		return nil
	}

	spriteIndex := e.activeAnimation.FrameData[e.frameIndex].SpriteIndex
	if spriteIndex < 0 || spriteIndex >= len(e.activeAnimation.Sprites) {
		return nil
	}

	return e.activeAnimation.Sprites[spriteIndex]
}

func (e *EditorManager) setActiveAnimation(anim *animation.Animation) {
	e.activeAnimation = anim
	e.frameIndex = 0
	e.frameCounter = 0 // Reset frame counter when switching animations
	e.playingAnim = false
	// Clear box editor when switching animations
	e.boxEditor = nil
}
