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
	FrameIndex                 uint // refers to the current frame in the current animation
	SpriteIndex                uint // refers to the number of the current animation sprite
	ShouldLoopCurrentAnimation bool // whether the current animation should loop

	AnimationQueue []string // queue of animations to play
	AnimationIndex int      // index of the current animation in the queue
}

func (am *AnimationManager) Update() {
	if am.CurrentAnim == nil {
		return
	}

	am.FrameIndex++
	if am.FrameIndex >= uint(len(am.CurrentAnim.Sprites)) {
		if am.ShouldLoopCurrentAnimation {
			am.FrameIndex = 0
		} else {
			am.CurrentAnim = nil
		}
	}

	// Update the current sprite
	if am.CurrentAnim != nil {
		am.CurrentSprite = am.CurrentAnim.Sprites[am.FrameIndex]
	}
}
