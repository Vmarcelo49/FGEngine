package animation

import (
	"fgengine/collision"
	"fgengine/types"
)

type Sprite struct {
	ImagePath string                             `yaml:"imgPath"`
	Rect      types.Rect                         `yaml:"rect"` // Position and Size
	Boxes     map[collision.BoxType][]types.Rect `yaml:"boxes"`
}

func (s *Sprite) Pos() types.Vector2 {
	return types.Vector2{X: s.Rect.X, Y: s.Rect.Y}
}

func (s *Sprite) Size() types.Vector2 {
	return types.Vector2{X: s.Rect.W, Y: s.Rect.H}
}
