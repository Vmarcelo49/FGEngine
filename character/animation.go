package character

type Animation struct {
	Name    string            `yaml:"name"`
	Sprites []*Sprite         `yaml:"sprites"`
	Prop    []FrameProperties `yaml:"properties"`
}
