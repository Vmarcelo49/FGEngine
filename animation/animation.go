package animation

import (
	"fgengine/types"
	"fmt"
)

type Animation struct {
	Name          string      `toml:"-"`
	Sprites       []*Sprite   `toml:"sprites"`
	FrameData     []FrameData `toml:"framedata"`
	TotalDuration int         `toml:"-"`

	LoopFrames *LoopFrame `toml:"loopFrames,omitempty"`
}

type LoopFrame struct {
	Start int `toml:"start"`
	End   int `toml:"end"`
}

func (ap *AnimationPlayer) Update(intentAnimation string, stunFrames int) {
	if ap.ActiveAnimation == nil || len(ap.ActiveAnimation.FrameData) == 0 {
		return
	}

	// Stun hold (SPEC §7.6): while stunned, hold within the active
	// animation's loopFrames range and never advance past it. start==end
	// holds a single frame. Without loopFrames, fall through to the
	// normal update below.
	if stunFrames > 0 {
		if lf := ap.ActiveAnimation.LoopFrames; lf != nil {
			last := len(ap.ActiveAnimation.FrameData) - 1
			start, end := lf.Start, lf.End
			if start < 0 {
				start = 0
			}
			if end > last {
				end = last
			}
			if start > end {
				start = end
			}
			ap.FrameTimeLeft--
			if ap.FrameTimeLeft > 0 {
				return
			}
			if ap.FrameIndex < start {
				ap.FrameIndex = start
			} else if ap.FrameIndex >= end {
				ap.FrameIndex = start
			} else {
				ap.FrameIndex++
				if ap.FrameIndex > end {
					ap.FrameIndex = start
				}
			}
			ap.FrameTimeLeft = ap.ActiveAnimation.FrameData[ap.FrameIndex].Duration
			return
		}
	}

	ap.FrameTimeLeft--
	if ap.FrameTimeLeft > 0 {
		return
	}

	ap.FrameIndex++

	// end of animation reached, stop at the last frame, should never happen because there is a fallback to idle, also helps not to put wrong values into the frameindex and point to nil frames.
	if ap.FrameIndex >= len(ap.ActiveAnimation.FrameData) {
		ap.FrameIndex = len(ap.ActiveAnimation.FrameData) - 1
		ap.FrameTimeLeft = 0
		return
	}

	loopFrames := ap.ActiveAnimation.LoopFrames

	if loopFrames != nil && loopFrames.Start != loopFrames.End {
		holding := ap.ActiveAnimation.Name == "idle" ||
			ap.ActiveAnimation.Name == intentAnimation

		if holding {
			// loop
			if ap.FrameIndex > loopFrames.End {
				ap.FrameIndex = loopFrames.Start
			}
		} else {
			// released the button
			if ap.FrameIndex >= loopFrames.Start &&
				ap.FrameIndex <= loopFrames.End {

				ap.FrameIndex = loopFrames.End + 1
				// A loop range reaching the last frame (e.g. a fully-looped
				// reaction whose stun just expired) has no trailing frames:
				// clamp to the last frame so it finishes instead of
				// indexing past the end.
				if ap.FrameIndex >= len(ap.ActiveAnimation.FrameData) {
					ap.FrameIndex = len(ap.ActiveAnimation.FrameData) - 1
				}
			}
		}
	}

	ap.FrameTimeLeft = ap.ActiveAnimation.FrameData[ap.FrameIndex].Duration
}

type Sprite struct {
	ImagePath string     `toml:"imgPath"`
	Rect      types.Rect `toml:"rect"`

	Anchor types.Vector2 `toml:"anchor,omitempty"`
}

type AnimationPlayer struct {
	ActiveAnimation *Animation            `toml:"-"`
	Animations      map[string]*Animation `toml:"-"`
	FrameIndex      int                   `toml:"-"`
	AnimationQueue  []string              `toml:"-"` // names are probably smaller than full Animation structs

	FrameTimeLeft int `toml:"-"`
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
		fmt.Printf("Missing animation %s\n", name)
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
