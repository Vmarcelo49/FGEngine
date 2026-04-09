package animation

import (
	"fgengine/types"
	"fmt"
)

type Animation struct {
	Name          string      `yaml:"-"`
	Sprites       []*Sprite   `yaml:"sprites"`
	FrameData     []FrameData `yaml:"framedata"`
	TotalDuration int         `yaml:"-"`
}

type Sprite struct {
	ImagePath string     `yaml:"imgPath"`
	Rect      types.Rect `yaml:"rect"`

	Anchor types.Vector2 `yaml:"anchor,omitempty"`
}

type AnimationPlayer struct {
	ActiveAnimation *Animation            `yaml:"-"`
	Animations      map[string]*Animation `yaml:"animations"`
	FrameIndex      int                   `yaml:"-"`
	ShouldLoop      bool                  `yaml:"-"`
	AnimationQueue  []string              `yaml:"-"` // names are probably smaller than full Animation structs

	FrameTimeLeft int `yaml:"-"`
}

func (ap *AnimationPlayer) ActiveSprite() *Sprite {
	if ap.ActiveAnimation == nil {
		return nil
	}
	frameData := ap.ActiveFrameData()
	if frameData == nil {
		return nil
	}
	return ap.ActiveAnimation.Sprites[frameData.SpriteIndex]
}

func (ap *AnimationPlayer) SetAnimation(name string) {
	if name == "" {
		return
	}
	if ap == nil || ap.Animations == nil {
		fmt.Println("Animation player has no animations map")
		return
	}

	anim, exists := ap.Animations[name]
	if !exists || anim == nil {

		fmt.Println(fmt.Sprintf("Animation '%s' not found", name))
		return
	}
	anim.Name = name
	ap.ActiveAnimation = anim
	ap.FrameIndex = 0
	if len(anim.FrameData) == 0 {
		ap.FrameTimeLeft = 0
		return
	}
	ap.FrameTimeLeft = anim.FrameData[0].Duration
}

func (ap *AnimationPlayer) Update() {
	if ap.ActiveAnimation == nil {
		return
	}

	// Don't update if animation has ended (non-looping)
	lastIndex := len(ap.ActiveAnimation.FrameData) - 1
	if !ap.ShouldLoop && ap.FrameIndex == lastIndex && ap.FrameTimeLeft <= 0 {
		return
	}

	ap.FrameTimeLeft--
	if ap.FrameTimeLeft > 0 {
		return
	}

	ap.FrameIndex++

	if ap.FrameIndex >= len(ap.ActiveAnimation.FrameData) {
		if ap.ShouldLoop {
			ap.FrameIndex = 0
		} else {
			ap.FrameIndex = lastIndex
			ap.FrameTimeLeft = 0
			return
		}
	}

	ap.FrameTimeLeft = ap.ActiveAnimation.FrameData[ap.FrameIndex].Duration
}

func (ap *AnimationPlayer) ActiveFrameData() *FrameData {
	if ap.ActiveAnimation == nil || len(ap.ActiveAnimation.FrameData) == 0 {
		return nil
	}
	return &ap.ActiveAnimation.FrameData[ap.FrameIndex]
}

func (ap *AnimationPlayer) ActiveAnimationName() string {
	if ap == nil || ap.ActiveAnimation == nil {
		return "none"
	}
	if ap.ActiveAnimation.Name != "" {
		return ap.ActiveAnimation.Name
	}

	for name, anim := range ap.Animations {
		if anim == ap.ActiveAnimation {
			ap.ActiveAnimation.Name = name
			return name
		}
	}

	return "none"
}

// IsFinished returns true if a non-looping animation has completed
func (ap *AnimationPlayer) IsFinished() bool {
	if ap.ActiveAnimation == nil || ap.ShouldLoop {
		return false
	}
	lastIndex := len(ap.ActiveAnimation.FrameData) - 1
	return ap.FrameIndex == lastIndex && ap.FrameTimeLeft <= 0
}

// Returns total duration in frames
func (a *Animation) Duration() int {
	if a.TotalDuration == 0 { // building this variable when called at least once, cus laziness to rewrite elsewhere
		for _, frameData := range a.FrameData {
			a.TotalDuration += frameData.Duration
		}
	}

	return a.TotalDuration
}

// Notes for future reference:

/*
Screenshake is better if the focus intensity on horizontal movement instead of vertical movement
*/
