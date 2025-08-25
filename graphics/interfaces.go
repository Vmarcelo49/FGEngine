package graphics

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/types"
)

// Renderable represents any game entity that can be rendered and animated
type Renderable interface {
	GetID() int
	GetPosition() types.Vector2
	GetAllBoxes() []collision.Box
	GetSprite() *animation.Sprite
}
