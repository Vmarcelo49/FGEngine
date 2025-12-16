package animation

type AnimationPlayer struct {
	ActiveAnimation *Animation `yaml:"-"`
	FrameCounter    int        `yaml:"-"`
	ShouldLoop      bool       `yaml:"-"`
	AnimationQueue  []string   `yaml:"-"` // names are probably smaller than full Animation structs
}

// Returns the framedata and its index inside the framedata slice
func (ap *AnimationPlayer) GetActiveFrameData() (*FrameData, int) {
	if ap.ActiveAnimation == nil || len(ap.ActiveAnimation.FrameData) == 0 { // protections for the editor
		return nil, -1
	}

	var elapsed, totalDuration int

	for _, f := range ap.ActiveAnimation.FrameData {
		totalDuration += f.Duration // TODO, remove this loop, add a totalDuration field to animations
	}
	frameCounter := ap.FrameCounter
	if ap.ShouldLoop {
		frameCounter = ap.FrameCounter % totalDuration
	} else if ap.FrameCounter >= totalDuration {
		lastIdx := len(ap.ActiveAnimation.FrameData) - 1
		return &ap.ActiveAnimation.FrameData[lastIdx], lastIdx
	}
	for i := range ap.ActiveAnimation.FrameData {
		frame := &ap.ActiveAnimation.FrameData[i]
		elapsed += frame.Duration
		if frameCounter < elapsed {
			return frame, i
		}
	}
	lastIdx := len(ap.ActiveAnimation.FrameData) - 1
	return &ap.ActiveAnimation.FrameData[lastIdx], lastIdx
}

func (ap *AnimationPlayer) GetSpriteFromFrameCounter() *Sprite {
	if ap.ActiveAnimation == nil {
		return nil
	}
	frameData, _ := ap.GetActiveFrameData()
	if frameData == nil {
		return nil
	}
	return ap.ActiveAnimation.Sprites[frameData.SpriteIndex]
}
