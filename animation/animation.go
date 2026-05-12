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

	LoopFrames *LoopFrame `yaml:"loopFrames,omitempty"`
}

type LoopFrame struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
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
	if frameData.SpriteIndex < 0 || frameData.SpriteIndex >= len(ap.ActiveAnimation.Sprites) {
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
		fmt.Printf("Animation '%s' not found", name)
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

func (ap *AnimationPlayer) Update(intentAnimation string) {
	if ap.ActiveAnimation == nil {
		return
	}

	ap.FrameTimeLeft--
	if ap.FrameTimeLeft > 0 {
		return
	}

	ap.FrameIndex++

	// end of animation reached, stop at the last frame, should never happen because there is a fallback to idle, but we dont want to underflow the index.
	if ap.FrameIndex >= len(ap.ActiveAnimation.FrameData) {
		ap.FrameIndex = len(ap.ActiveAnimation.FrameData) - 1
		ap.FrameTimeLeft = 0
		return
	}

	// loopframes[0] is the index to loop back to, loopframes[1] is the last frame of the loop (inclusive)
	if ap.ActiveAnimation.LoopFrames != nil { // basically if the animation has loop frames
		// and is idle or the same animation as the intent animation, loop
		if ap.ActiveAnimation.Name == "idle" || ap.ActiveAnimation.Name == intentAnimation {
			if ap.FrameIndex > ap.ActiveAnimation.LoopFrames.End {
				ap.FrameIndex = ap.ActiveAnimation.LoopFrames.Start
			}
		}
		if ap.ActiveAnimation.Name != "idle" && ap.ActiveAnimation.Name != intentAnimation {
			ap.FrameIndex = ap.ActiveAnimation.LoopFrames.End // try to switch to recovery frames
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
	if ap.ActiveAnimation == nil {
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
