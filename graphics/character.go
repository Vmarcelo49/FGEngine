package graphics

import (
	"FGEngine/character"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawCharacter(character *character.Character, screen *ebiten.Image) {
	if checkDrawConditions(character) == false {
		return
	}

	characterImage := loadCharacterImage(character)
	if characterImage == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(character.Position.X, character.Position.Y)
	screen.DrawImage(characterImage, options)
}

// TODO, add sprite priority to control the z-index
// TODO, child sprites for projectiles and other entities
