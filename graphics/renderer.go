package graphics

import (
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

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

	applyCameraTransform(options, camera, renderable, screenPos)

	screen.DrawImage(entityImage, options)
}

// applyCameraTransform applies the camera transformation logic to drawing options
func applyCameraTransform(options *ebiten.DrawImageOptions, camera *Camera, renderable Renderable, screenPos types.Vector2) {
	if layoutMatchesCamSize(camera) {
		// Apply scaling around center of viewport
		if camera.Scaling != 0 && camera.Scaling != 1 {
			// Calculate center of viewport (in screen space)
			centerX := camera.Viewport.W / 2
			centerY := camera.Viewport.H / 2

			// Calculate relative position from center
			relativeX := screenPos.X - centerX
			relativeY := screenPos.Y - centerY

			// Scale the relative position
			scaledRelativeX := relativeX * camera.Scaling
			scaledRelativeY := relativeY * camera.Scaling

			// Calculate final position (back from center)
			finalX := scaledRelativeX + centerX
			finalY := scaledRelativeY + centerY

			options.GeoM.Scale(camera.Scaling, camera.Scaling)
			options.GeoM.Translate(finalX, finalY)
		} else {
			options.GeoM.Translate(screenPos.X, screenPos.Y)
		}
	} else {
		zoomAroundCenterOption(options, camera, renderable, screenPos)
	}
}

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

func DrawStaticWithScale(renderable Renderable, screen *ebiten.Image) {
	entityImage := loadImage(renderable)
	if entityImage == nil {
		return
	}
	options := &ebiten.DrawImageOptions{}
	pos := renderable.GetPosition()
	options.GeoM.Scale(renderable.GetRenderProperties().Scale.X, renderable.GetRenderProperties().Scale.Y)
	options.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(entityImage, options)
}
