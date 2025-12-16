package animation

type Animation struct {
	Name          string      `yaml:"name,omitempty"`
	Sprites       []*Sprite   `yaml:"sprites"`
	FrameData     []FrameData `yaml:"framedata"`
	TotalDuration int
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
