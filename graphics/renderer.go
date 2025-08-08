package graphics

import (
	"FGEngine/character"
)

// Placeholder

func checkDrawConditions(character *character.Character) bool {
	if character != nil {
		if character.CurrentSprite != nil {
			return true
		}
	}
	return false
}
