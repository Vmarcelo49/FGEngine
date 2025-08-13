package animation

// Animation represents animation data loaded from character files
type Animation struct {
	Name    string            `yaml:"name"`
	Sprites []*Sprite         `yaml:"sprites"`
	Prop    []FrameProperties `yaml:"properties"`
}

type AnimationManager struct {
	ID                         int
	CurrentAnim                *Animation
	CurrentSprite              *Sprite
	FrameIndex                 uint
	SpriteIndex                uint
	ShouldLoopCurrentAnimation bool
	AnimationQueue             []string
	AnimationIndex             int
	Character                  *Character
	AnimationCache             map[string]*Animation
}

// CreateAnimationManager creates a new animation manager for a character
func CreateAnimationManager(id int, char *Character) *AnimationManager {
	return &AnimationManager{
		ID:                         id,
		Character:                  char,
		ShouldLoopCurrentAnimation: true, // Default to looping animations
		FrameIndex:                 0,
		SpriteIndex:                0,
		AnimationQueue:             make([]string, 0),
		AnimationIndex:             0,
		AnimationCache:             make(map[string]*Animation),
	}
}

func (am *AnimationManager) Update() {
	if am.CurrentAnim == nil || len(am.CurrentAnim.Sprites) == 0 {
		return
	}

	am.FrameIndex++
	if am.FrameIndex >= am.CurrentSprite.Duration {
		am.FrameIndex = 0
		am.SpriteIndex++

		if am.SpriteIndex >= uint(len(am.CurrentAnim.Sprites)) {
			if am.ShouldLoopCurrentAnimation {
				am.SpriteIndex = 0 // Loop back to first sprite
			} else {
				am.SpriteIndex = uint(len(am.CurrentAnim.Sprites)) - 1 // Stay on last sprite
				// Animation has finished - could trigger events here later
			}
		}
	}

	// Update current sprite reference
	if am.SpriteIndex < uint(len(am.CurrentAnim.Sprites)) {
		am.CurrentSprite = am.CurrentAnim.Sprites[am.SpriteIndex]
	}
}

func (am *AnimationManager) SetAnimation(animName string) {
	if am.Character == nil {
		return
	}

	// Check cache first
	if anim, exists := am.AnimationCache[animName]; exists {
		am.setCurrentAnimation(anim)
		return
	}

	// Get animation data directly from character
	if anim, exists := am.Character.Animations[animName]; exists {
		am.AnimationCache[animName] = anim
		am.setCurrentAnimation(anim)
	}
}

// setCurrentAnimation is a helper to set the current animation
func (am *AnimationManager) setCurrentAnimation(anim *Animation) {
	// Only change if it's a different animation
	if am.CurrentAnim == nil || am.CurrentAnim.Name != anim.Name {
		am.CurrentAnim = anim
		am.SpriteIndex = 0
		am.FrameIndex = 0

		// Set default looping behavior - most animations should loop
		// This could be made data-driven later
		am.ShouldLoopCurrentAnimation = true

		// Safely set the current sprite
		if len(anim.Sprites) > 0 {
			am.CurrentSprite = anim.Sprites[0]
		}
	}
}

func (am *AnimationManager) GetCurrentSprite() *Sprite {
	return am.CurrentSprite
}

func (am *AnimationManager) IsValid() bool {
	return am.CurrentSprite != nil && am.CurrentAnim != nil
}

// SetLooping sets whether the current animation should loop
func (am *AnimationManager) SetLooping(shouldLoop bool) {
	am.ShouldLoopCurrentAnimation = shouldLoop
}

// IsAnimationFinished returns true if a non-looping animation has finished
func (am *AnimationManager) IsAnimationFinished() bool {
	if am.CurrentAnim == nil || am.ShouldLoopCurrentAnimation {
		return false
	}
	return am.SpriteIndex >= uint(len(am.CurrentAnim.Sprites))-1 && am.FrameIndex >= am.CurrentSprite.Duration-1
}

// GetCurrentAnimationName returns the name of the current animation, or empty string if none
func (am *AnimationManager) GetCurrentAnimationName() string {
	if am.CurrentAnim == nil {
		return ""
	}
	return am.CurrentAnim.Name
}

// RestartCurrentAnimation resets the current animation to the beginning
func (am *AnimationManager) RestartCurrentAnimation() {
	if am.CurrentAnim != nil && len(am.CurrentAnim.Sprites) > 0 {
		am.SpriteIndex = 0
		am.FrameIndex = 0
		am.CurrentSprite = am.CurrentAnim.Sprites[0]
	}
}
