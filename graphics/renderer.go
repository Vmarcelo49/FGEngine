package graphics

import (
	"FGEngine/player"
)

// Placeholder

func checkDrawConditions(p *player.Player) bool {
	if p != nil && p.State.AnimationManager != nil {
		if p.State.AnimationManager.CurrentSprite != nil {
			return true
		}
	}
	return false
}
