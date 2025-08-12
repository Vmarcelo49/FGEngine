package animation

import "FGEngine/character"

type AnimationComponent struct {
	ID                         int
	CurrentAnim                *character.Animation
	CurrentSprite              *character.Sprite
	FrameIndex                 uint
	SpriteIndex                uint
	ShouldLoopCurrentAnimation bool
	AnimationQueue             []string
	AnimationIndex             int
	CharacterRef               *character.Character
}

// NewAnimationComponent creates a new animation component with a reference to the character
func NewAnimationComponent(id int, char *character.Character) *AnimationComponent {
	return &AnimationComponent{
		ID:           id,
		CharacterRef: char,
	}
}

func (ac *AnimationComponent) Update() {
	if ac.CurrentAnim == nil {
		return
	}

	ac.FrameIndex++
	if ac.FrameIndex >= ac.CurrentSprite.Duration {
		ac.FrameIndex = 0
		if ac.SpriteIndex < uint(len(ac.CurrentAnim.Sprites))-1 {
			ac.SpriteIndex++
			if ac.SpriteIndex >= uint(len(ac.CurrentAnim.Sprites)) {
				ac.SpriteIndex = 0
			}
		}
		if !ac.ShouldLoopCurrentAnimation {
			ac.SpriteIndex = 0
		}
	}

	if ac.CurrentAnim != nil {
		ac.CurrentSprite = ac.CurrentAnim.Sprites[ac.SpriteIndex]
	}
}

func (ac *AnimationComponent) SetAnimation(animName string) {
	if ac.CharacterRef == nil {
		return
	}

	for _, anim := range ac.CharacterRef.Animations {
		if anim.Name == animName {
			ac.CurrentAnim = anim
			ac.CurrentSprite = anim.Sprites[0]
			ac.SpriteIndex = 0
			ac.FrameIndex = 0
			return
		}
	}
}

func (ac *AnimationComponent) GetCurrentSprite() *character.Sprite {
	return ac.CurrentSprite
}

func (ac *AnimationComponent) IsValid() bool {
	return ac.CurrentSprite != nil && ac.CurrentAnim != nil
}
