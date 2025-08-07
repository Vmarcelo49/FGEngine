package graphics

import (
	"FGEngine/character"
)

// Placeholder

func checkDrawConditions(character *character.Character) bool {
	var isValid bool
	if character != nil {
		if character.CurrentSprite != nil {
			isValid = true
		}
	}
	return isValid
}
