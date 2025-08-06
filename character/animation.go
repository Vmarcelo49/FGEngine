package character

type Animation struct {
	Name    string       `yaml:"name"`
	Sprites []*SpriteEx  `yaml:"sprites"`
	Prop    []Properties `yaml:"properties"`
}

type AnimationManager struct {
	Animations                 map[string]*Animation `yaml:"animations"`
	CurrentAnim                *Animation
	CurrentSprite              *SpriteEx
	FrameIndex                 uint // current frame in the current animation, mainly interacts with SpriteEx.Duration
	SpriteIndex                uint // current sprite in the current animation
	ShouldLoopCurrentAnimation bool // whether the current animation should loop

	AnimationQueue []string // queue of animations to play
	AnimationIndex int      // index of the current animation in the queue
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
