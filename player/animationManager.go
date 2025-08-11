package player

import "FGEngine/character"

type AnimationManager struct {
	CurrentAnim                *character.Animation
	CurrentSprite              *character.SpriteEx
	FrameIndex                 uint
	SpriteIndex                uint
	ShouldLoopCurrentAnimation bool
	AnimationQueue             []string
	AnimationIndex             int
}

func (am *AnimationManager) Update() {
	if am.CurrentAnim == nil {
		return
	}

	am.FrameIndex++
	if am.FrameIndex >= am.CurrentSprite.Duration { // should advance to the next sprite
		am.FrameIndex = 0
		if am.SpriteIndex < uint(len(am.CurrentAnim.Sprites))-1 { // if in range of current animation len of sprites
			if am.SpriteIndex < uint(len(am.CurrentAnim.Sprites)-1) { // if in range of current animation len of sprites
				am.SpriteIndex++
				if am.SpriteIndex >= uint(len(am.CurrentAnim.Sprites)) { // if out of range of current animation len of sprites
					am.SpriteIndex = 0
				}
			}
		}
		if !am.ShouldLoopCurrentAnimation { // go to next animation in queue
			am.SpriteIndex = 0
		}
	}

	// Update the current sprite
	if am.CurrentAnim != nil {
		am.CurrentSprite = am.CurrentAnim.Sprites[am.SpriteIndex]
	}
}
