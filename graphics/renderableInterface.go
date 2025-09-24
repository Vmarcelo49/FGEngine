package graphics

import (
	"fgengine/animation"
	"fgengine/types"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Renderable interface - all renderable objects must implement these
type Renderable interface {
	GetPosition() types.Vector2
	GetSprite() *animation.Sprite
	GetRenderProperties() RenderProperties
}

// RenderProperties provides rendering properties for all renderable objects
type RenderProperties struct {
	Scale    types.Vector2 // 1.0 = normal size
	Rotation float64       // in radians
	Layer    int           // Higher numbers render on top (0 = default)
	ColorMod color.RGBA    // Color modulation (white = no change)
}

// DefaultRenderProperties returns default rendering properties
func DefaultRenderProperties() RenderProperties {
	return RenderProperties{
		Scale:    types.Vector2{X: 1.0, Y: 1.0},
		Rotation: 0.0,
		Layer:    0,
		ColorMod: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White = no change
	}
}

func applyScalingTransform(options *ebiten.DrawImageOptions, screenPos types.Vector2, camera *Camera, scaleX, scaleY float64) {
	options.GeoM.Scale(scaleX*camera.Scaling, scaleY*camera.Scaling)
	options.GeoM.Translate(screenPos.X, screenPos.Y)
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
