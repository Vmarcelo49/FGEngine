package graphics

// Placeholder

// TODO, this must be removed, we must implement a proper initialization
func checkDrawConditions(renderable Renderable) bool {
	if renderable != nil && renderable.GetAnimationComponent() != nil {
		return renderable.GetAnimationComponent().IsValid()
	}
	return false
}
