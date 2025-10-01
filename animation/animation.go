package animation

type Animation struct {
	Name    string            `yaml:"name,omitempty"`
	Sprites []*Sprite         `yaml:"sprites"`
	Prop    []FrameProperties `yaml:"properties"`
}

// Returns total duration in frames
func (a *Animation) Duration() int {
	var duration int
	for _, prop := range a.Prop {
		duration += prop.Duration
	}
	return duration
}
