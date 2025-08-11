package character

type Animation struct {
	Name    string            `yaml:"name"`
	Sprites []*SpriteEx       `yaml:"sprites"`
	Prop    []FrameProperties `yaml:"properties"`
}
