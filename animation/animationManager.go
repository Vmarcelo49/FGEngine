package animation

// Animation represents animation data loaded from character files
type Animation struct {
	Name    string            `yaml:"name"`
	Sprites []*Sprite         `yaml:"sprites"`
	Prop    []FrameProperties `yaml:"properties"`
}

type AnimationSystem struct {
	ID                         int
	CurrentAnim                *Animation
	CurrentSprite              *Sprite
	FrameIndex                 uint
	SpriteIndex                uint
	ShouldLoopCurrentAnimation bool
	AnimationQueue             []string
	AnimationIndex             int
	AnimationCache             map[string]*Animation
}

// NewAnimationSystem creates a new animation manager for a character
func NewAnimationSystem(id int) *AnimationSystem {
	return &AnimationSystem{
		ID:                         id,
		ShouldLoopCurrentAnimation: true, // Default to looping animations
		AnimationCache:             make(map[string]*Animation),
	}
}

func (as *AnimationSystem) Update() {
	if as.CurrentAnim == nil || len(as.CurrentAnim.Sprites) == 0 {
		return
	}

	as.FrameIndex++
	if as.FrameIndex >= as.CurrentSprite.Duration {
		as.FrameIndex = 0
		as.SpriteIndex++

		if as.SpriteIndex >= uint(len(as.CurrentAnim.Sprites)) {
			if as.ShouldLoopCurrentAnimation {
				as.SpriteIndex = 0 // Loop back to first sprite
			} else {
				as.SpriteIndex = uint(len(as.CurrentAnim.Sprites)) - 1 // Stay on last sprite
				// Animation has finished - could trigger events here later
			}
		}
	}

	// Update current sprite reference
	if as.SpriteIndex < uint(len(as.CurrentAnim.Sprites)) {
		as.CurrentSprite = as.CurrentAnim.Sprites[as.SpriteIndex]
	}
}

func (as *AnimationSystem) SetAnimation(animName string) {
	if as.Character == nil {
		return
	}

	// Check cache first
	if anim, exists := as.AnimationCache[animName]; exists {
		as.setAnimation(anim)
		return
	}

	// Get animation data directly from character
	if anim, exists := as.Character.Animations[animName]; exists {
		as.AnimationCache[animName] = anim
		as.setAnimation(anim)
	}
}

func (as *AnimationSystem) setAnimation(anim *Animation) {
	// Only change if it's a different animation
	if as.CurrentAnim == nil || as.CurrentAnim.Name != anim.Name {
		as.CurrentAnim = anim
		as.SpriteIndex = 0
		as.FrameIndex = 0

		// Set default looping behavior - most animations should loop
		// This could be made data-driven later
		as.ShouldLoopCurrentAnimation = true

		// Safely set the current sprite
		if len(anim.Sprites) > 0 {
			as.CurrentSprite = anim.Sprites[0]
		}
	}
}

func (as *AnimationSystem) GetCurrentSprite() *Sprite {
	return as.CurrentSprite
}

func (as *AnimationSystem) IsValid() bool {
	return as.CurrentSprite != nil && as.CurrentAnim != nil
}

// SetLooping sets whether the current animation should loop
func (as *AnimationSystem) SetLooping(shouldLoop bool) {
	as.ShouldLoopCurrentAnimation = shouldLoop
}

// IsAnimationFinished returns true if a non-looping animation has finished
func (as *AnimationSystem) IsAnimationFinished() bool {
	if as.CurrentAnim == nil || as.ShouldLoopCurrentAnimation {
		return false
	}
	return as.SpriteIndex >= uint(len(as.CurrentAnim.Sprites))-1 && as.FrameIndex >= as.CurrentSprite.Duration-1
}

// GetCurrentAnimationName returns the name of the current animation, or empty string if none
func (as *AnimationSystem) GetCurrentAnimationName() string {
	if as.CurrentAnim == nil {
		return ""
	}
	return as.CurrentAnim.Name
}

// RestartCurrentAnimation resets the current animation to the beginning
func (as *AnimationSystem) RestartCurrentAnimation() {
	if as.CurrentAnim != nil && len(as.CurrentAnim.Sprites) > 0 {
		as.SpriteIndex = 0
		as.FrameIndex = 0
		as.CurrentSprite = as.CurrentAnim.Sprites[0]
	}
}
