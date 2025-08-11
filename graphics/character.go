package graphics

import (
	"FGEngine/player"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawPlayer(p *player.Player, screen *ebiten.Image) {
	if checkDrawConditions(p) == false {
		return
	}

	characterImage := loadPlayerImage(p)
	if characterImage == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(p.State.Position.X, p.State.Position.Y)
	screen.DrawImage(characterImage, options)
}

// TODO, add sprite priority to control the z-index
// TODO, child sprites for projectiles and other entities
