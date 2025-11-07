package animation

type Animation struct {
	Name      string      `yaml:"name,omitempty"`
	Sprites   []*Sprite   `yaml:"sprites"`
	FrameData []FrameData `yaml:"framedata"`
}

// Returns total duration in frames
func (a *Animation) Duration() int {
	var duration int
	for _, frameData := range a.FrameData {
		duration += frameData.Duration
	}
	return duration
}

// Notes for future reference:

/*
Screenshake is better if the focus intensity on horizontal movement instead of vertical movement
*/
