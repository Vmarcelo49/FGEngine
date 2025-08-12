package graphics

import (
	"FGEngine/animation"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawRenderable draws any renderable entity (player, projectile, etc.)
func DrawRenderable(renderable animation.Renderable, screen *ebiten.Image) {
	if checkDrawConditions(renderable) == false {
		fmt.Println("Entity is not in a drawable state")
		return
	}

	entityImage := loadRenderableImage(renderable)
	if entityImage == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}
	pos := renderable.GetPosition()
	options.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(entityImage, options)
}

// TODO, add sprite priority to control the z-index
// TODO, child sprites for projectiles and other entities
