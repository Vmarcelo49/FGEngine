package animation

import (
	"fgengine/collision"
	"fgengine/types"
)

type Sprite struct {
	ImagePath string                             `yaml:"imgPath"`
	Rect      types.Rect                         `yaml:"rect"`
	Boxes     map[collision.BoxType][]types.Rect `yaml:"boxes,omitempty"`
}
