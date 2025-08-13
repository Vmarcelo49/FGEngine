package animation

import "fgengine/types"

type Sprite struct {
	ImagePath      string       `yaml:"imgPath"`
	Duration       uint         `yaml:"duration"`
	SourceSize     types.Rect   `yaml:"sourceSize"`
	Anchor         types.Rect   `yaml:"anchorPoint"`
	CollisionBoxes []types.Rect `yaml:"collisionBoxes"`
	HurtBoxes      []types.Rect `yaml:"hurtBoxes"`
	HitBoxes       []types.Rect `yaml:"hitBoxes"`
}
