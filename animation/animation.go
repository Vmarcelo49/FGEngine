package animation

type Animation struct {
	Name      string      `yaml:"name,omitempty"`
	Sprites   []*Sprite   `yaml:"sprites"`
	FrameData []FrameData `yaml:"framedata"`
}

// Returns total duration in frames
func (a *Animation) Duration() int {
	var duration int
	for _, fdata := range a.FrameData {
		duration += fdata.Duration
	}
	return duration
}
