package graphics

import (
	"fgengine/animation"
	"fgengine/types"
	"image/color"
)

type Renderable interface {
	GetPosition() types.Vector2
	GetSprite() *animation.Sprite
	GetRenderProperties() RenderProperties
}

type RenderProperties struct {
	Scale    types.Vector2 // 1.0 = normal size
	Rotation float64       // in radians
	Layer    int           // Higher numbers render on top (0 = default)
	ColorMod color.RGBA    // Color modulation (white = no change)
}

func DefaultRenderProperties() RenderProperties {
	return RenderProperties{
		Scale:    types.Vector2{X: 1.0, Y: 1.0},
		Rotation: 0.0,
		Layer:    0,
		ColorMod: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White = no change
	}
}
