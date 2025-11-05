package animation

import (
	"fgengine/types"
)

type Sprite struct {
	ImagePath string     `yaml:"imgPath"`
	Rect      types.Rect `yaml:"rect"`
}
