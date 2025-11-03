package animation

type AnimationPlayer struct {
	ActiveAnimation *Animation `yaml:"-"`
	FrameCounter    int        `yaml:"-"`
	ShouldLoop      bool       `yaml:"-"`
	AnimationQueue  []string   `yaml:"-"` // names are probably smaller than full Animation structs
}

func (ap *AnimationPlayer) GetSpriteFromFrameCounter() *Sprite {
	if ap == nil || ap.ActiveAnimation == nil {
		return nil
	}
	frameData := ap.GetActiveFrameData()
	if frameData == nil {
		return nil
	}
	return ap.ActiveAnimation.Sprites[frameData.SpriteIndex]
}

func (ap *AnimationPlayer) GetActiveFrameData() *FrameData {
	if ap == nil || ap.ActiveAnimation == nil || len(ap.ActiveAnimation.FrameData) == 0 { // protections for the editor
		return nil
	}

	var elapsed, totalDuration int

	for _, f := range ap.ActiveAnimation.FrameData {
		totalDuration += f.Duration
	}
	frameCounter := ap.FrameCounter % totalDuration // Loop around if exceeds total duration, if looping

	for i := range ap.ActiveAnimation.FrameData {
		frame := &ap.ActiveAnimation.FrameData[i]
		elapsed += frame.Duration
		if frameCounter < elapsed {
			return frame
		}
	}
	return &ap.ActiveAnimation.FrameData[len(ap.ActiveAnimation.FrameData)-1]
}
