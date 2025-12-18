package animation

import "fgengine/types"

type Animation struct {
	Name          string      `yaml:"name,omitempty"`
	Sprites       []*Sprite   `yaml:"sprites"`
	FrameData     []FrameData `yaml:"framedata"`
	TotalDuration int
}

type Sprite struct {
	ImagePath string     `yaml:"imgPath"`
	Rect      types.Rect `yaml:"rect"`
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

type AnimationPlayer struct {
	ActiveAnimation *Animation `yaml:"-"`
	FrameIndex      int        `yaml:"-"`
	ShouldLoop      bool       `yaml:"-"`
	AnimationQueue  []string   `yaml:"-"` // names are probably smaller than full Animation structs

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

// IsFinished returns true if a non-looping animation has completed
func (ap *AnimationPlayer) IsFinished() bool {
	if ap.ActiveAnimation == nil || ap.ShouldLoop {
		return false
	}
	lastIndex := len(ap.ActiveAnimation.FrameData) - 1
	return ap.FrameIndex == lastIndex && ap.FrameTimeLeft <= 0
}

// Notes for future reference:

/*
Screenshake is better if the focus intensity on horizontal movement instead of vertical movement
*/
