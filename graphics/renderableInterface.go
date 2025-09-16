package graphics

import (
	"fgengine/animation"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

// Renderable represents any game entity that can be rendered and animated
type Renderable interface {
	GetID() int // maybe useful for debugging?
	GetPosition() types.Vector2
	GetSprite() *animation.Sprite
}

// DrawRenderable draws any renderable entity (player, projectile, etc.)
func DrawRenderable(renderable Renderable, screen *ebiten.Image) {
	entityImage := loadImage(renderable)
	if entityImage == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}
	pos := renderable.GetPosition()
	options.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(entityImage, options)
}

func DrawRenderableWithScale(renderable Renderable, screen *ebiten.Image, scaleX, scaleY float64) {
	entityImage := loadImage(renderable)
	if entityImage == nil {
		return
	}
	options := &ebiten.DrawImageOptions{}
	pos := renderable.GetPosition()
	if scaleX == 0 {
		scaleX = 1
	}
	if scaleY == 0 {
		scaleY = 1
	}
	options.GeoM.Scale(scaleX, scaleY)
	options.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(entityImage, options)
}

// TODO, add sprite priority to control the z-index
// TODO, child sprites for projectiles and other entities
