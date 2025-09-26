package animation

// Animation represents animation data loaded from character files
type Animation struct {
	Name    string            `yaml:"name"`
	Sprites []*Sprite         `yaml:"sprites"`
	Prop    []FrameProperties `yaml:"properties"`
}

// Duration returns the total duration of the animation in frames
func (a *Animation) Duration() int {
	var duration int
	for _, prop := range a.Prop {
		duration += prop.Duration
	}
	return duration
}
