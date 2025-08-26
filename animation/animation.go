package animation

// Animation represents animation data loaded from character files
type Animation struct {
	Name    string            `yaml:"name"`
	Sprites []*Sprite         `yaml:"sprites"`
	Prop    []FrameProperties `yaml:"properties"`
}
