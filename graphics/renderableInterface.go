package graphics

import (
	"fgengine/animation"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderable interface {
	GetPosition() types.Vector2
	GetSprite() *animation.Sprite
}

func applyScalingTransform(options *ebiten.DrawImageOptions, screenPos types.Vector2, camera *Camera, scaleX, scaleY float64) {
	screenCenterX := camera.Viewport.W / 2
	screenCenterY := camera.Viewport.H / 2

	relativeX := screenPos.X - screenCenterX
	relativeY := screenPos.Y - screenCenterY

	scaledRelativeX := relativeX * camera.Scaling
	scaledRelativeY := relativeY * camera.Scaling

	finalX := scaledRelativeX + screenCenterX
	finalY := scaledRelativeY + screenCenterY

	options.GeoM.Scale(scaleX*camera.Scaling, scaleY*camera.Scaling)
	options.GeoM.Translate(finalX, finalY)
}

func applyBasicTransform(options *ebiten.DrawImageOptions, screenPos types.Vector2, scaleX, scaleY float64) {
	options.GeoM.Scale(scaleX, scaleY)
	options.GeoM.Translate(screenPos.X, screenPos.Y)
}

func DrawStatic(renderable Renderable, screen *ebiten.Image) {
	entityImage := loadImage(renderable)
	if entityImage == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}
	pos := renderable.GetPosition()
	applyBasicTransform(options, pos, 1.0, 1.0)
	screen.DrawImage(entityImage, options)
}

func Draw(renderable Renderable, screen *ebiten.Image, camera *Camera) {
	if !camera.IsVisible(renderable) {
		return
	}

	entityImage := loadImage(renderable)
	if entityImage == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}
	worldPos := renderable.GetPosition()
	screenPos := camera.WorldToScreen(worldPos)

	if camera.Scaling != 0 && camera.Scaling != 1 {
		applyScalingTransform(options, screenPos, camera, 1.0, 1.0)
	} else {
		applyBasicTransform(options, screenPos, 1.0, 1.0)
	}

	screen.DrawImage(entityImage, options)
}

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
	applyBasicTransform(options, pos, scaleX, scaleY)
	screen.DrawImage(entityImage, options)
}
