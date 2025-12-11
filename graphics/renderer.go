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

	ApplyCameraTransform(options, camera, renderable, screenPos)

	screen.DrawImage(entityImage, options)
}

// ApplyCameraTransform applies the camera transformation logic to drawing options
func ApplyCameraTransform(options *ebiten.DrawImageOptions, camera *Camera, renderable Renderable, screenPos types.Vector2) {
	// Sempre aplica escala do renderable primeiro
	renderProps := renderable.GetRenderProperties()
	options.GeoM.Scale(renderProps.Scale.X, renderProps.Scale.Y)

	// Calcula o centro do viewport uma vez
	centerX := camera.Viewport.W / 2
	centerY := camera.Viewport.H / 2

	// Aplica zoom da câmera se necessário
	if camera.Scaling != 0 && camera.Scaling != 1 {
		// Escala ao redor do centro do viewport
		options.GeoM.Translate(-centerX, -centerY)
		options.GeoM.Scale(camera.Scaling, camera.Scaling)
		options.GeoM.Translate(centerX, centerY)

		// scale again?
		if !layoutMatchesCamSize(camera) {
			options.GeoM.Scale(renderProps.Scale.X, renderProps.Scale.Y)
		}
	}

	// Aplica posição final
	options.GeoM.Translate(screenPos.X, screenPos.Y)
}

func ApplyCameraTransformNEW(options *ebiten.DrawImageOptions, camera *Camera, scaling types.Vector2, screenPos types.Vector2) {
	// Sempre aplica escala do renderable primeiro
	options.GeoM.Scale(scaling.X, scaling.Y)

	// Calcula o centro do viewport uma vez
	centerX := camera.Viewport.W / 2
	centerY := camera.Viewport.H / 2

	// Aplica zoom da câmera se necessário
	if camera.Scaling != 0 && camera.Scaling != 1 {
		// Escala ao redor do centro do viewport
		options.GeoM.Translate(-centerX, -centerY)
		options.GeoM.Scale(camera.Scaling, camera.Scaling)
		options.GeoM.Translate(centerX, centerY)

		// scale again?
		if !layoutMatchesCamSize(camera) {
			options.GeoM.Scale(scaling.X, scaling.Y)
		}
	}

	// Aplica posição final
	options.GeoM.Translate(screenPos.X, screenPos.Y)
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
