package graphics

import (
	"FGEngine/player"
)

// Placeholder

// TODO, this must be removed, we must implement a proper initialization
func checkDrawConditions(p *player.Player) bool {
	if p != nil && p.AnimationManager != nil {
		if p.AnimationManager.CurrentSprite != nil {
			return true
		}
	}
	return false
}
