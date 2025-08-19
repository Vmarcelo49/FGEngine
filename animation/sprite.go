package animation

import "fgengine/types"

// TODO, check if sprite can be simplified and other info be sent do frame properties
type Sprite struct {
	ImagePath      string       `yaml:"imgPath"`
	Duration       uint         `yaml:"duration"`
	SourceSize     types.Rect   `yaml:"sourceSize"`
	Anchor         types.Rect   `yaml:"anchorPoint"`
	CollisionBoxes []types.Rect `yaml:"collisionBoxes"`
	HurtBoxes      []types.Rect `yaml:"hurtBoxes"`
	HitBoxes       []types.Rect `yaml:"hitBoxes"`
}
