package graphics

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/types"
)

// Renderable represents any game entity that can be rendered and animated
type Renderable interface {
	GetPosition() types.Vector2
	GetAnimationComponent() *animation.AnimationManager
	GetAllBoxes() []collision.Box
	GetID() int
}
