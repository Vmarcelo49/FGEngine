package graphics

import (
	"fgengine/animation"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

// Renderable represents any game entity that can be rendered
type Renderable interface {
	GetPosition() types.Vector2
	GetSprite() *animation.Sprite
}

// DrawStatic draws a renderable at its exact world position (no camera transformation)
func DrawStatic(renderable Renderable, screen *ebiten.Image) {
	entityImage := loadImage(renderable)
	if entityImage == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}
	pos := renderable.GetPosition()
	options.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(entityImage, options)
}

// DrawWithCamera draws a renderable with camera transformation applied
func DrawWithCamera(renderable Renderable, screen *ebiten.Image, camera *Camera) {
	if !camera.IsVisible(renderable) {
		return // Skip rendering if outside camera view
	}

	entityImage := loadImage(renderable)
	if entityImage == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}
	worldPos := renderable.GetPosition()
	screenPos := camera.WorldToScreen(worldPos)
	if camera.Scaling != 0 && camera.Scaling != 1 {
		options.GeoM.Scale(camera.Scaling, camera.Scaling)
	}
	options.GeoM.Translate(screenPos.X, screenPos.Y)
	screen.DrawImage(entityImage, options)
}

// DrawStaticWithScale draws a renderable at exact position with scaling (no camera)
func DrawStaticWithScale(renderable Renderable, screen *ebiten.Image, scaleX, scaleY float64) {
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
