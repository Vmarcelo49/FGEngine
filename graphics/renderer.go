package graphics

import (
	"FGEngine/animation"
)

// Placeholder

// TODO, this must be removed, we must implement a proper initialization
func checkDrawConditions(renderable animation.Renderable) bool {
	if renderable != nil && renderable.GetAnimationComponent() != nil {
		return renderable.GetAnimationComponent().IsValid()
	}
	return false
}
