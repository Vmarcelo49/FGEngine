package graphics

import (
	"image/color"
	"slices"
	"sync"

	"fgengine/animation"
	"fgengine/collision"
	"fgengine/config"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	whitePixel *ebiten.Image
	once       sync.Once
	boxColors  = map[collision.BoxType]color.RGBA{
		collision.Collision: {R: 80, G: 80, B: 80, A: 32},
		collision.Hit:       {R: 100, G: 40, B: 40, A: 32},
		collision.Hurt:      {R: 40, G: 100, B: 40, A: 32},
	}
)

// whitePixel should never be unloaded and only is created once
func initWhitePixel() {
	once.Do(func() {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	})
}

// DrawBoxes draws all collision boxes with camera transformation applied
func DrawBoxes(renderable Renderable, screen *ebiten.Image, camera *Camera) {
	initWhitePixel()
	currentSprite := renderable.GetSprite()
	if currentSprite == nil {
		return
	}

	// consistent rendering order
	boxTypes := make([]collision.BoxType, 0, len(currentSprite.Boxes))
	for boxType := range currentSprite.Boxes {
		boxTypes = append(boxTypes, boxType)
	}
	slices.Sort(boxTypes)

	for _, boxType := range boxTypes {
		boxes := currentSprite.Boxes[boxType]
		for _, box := range boxes {
			options := createBoxImageOptionsWithCamera(renderable, box, boxType, camera)
			screen.DrawImage(whitePixel, options)
		}
	}

}

// DrawBoxesByType draws all collision boxes by type at exact world position (no camera)
func DrawBoxesByType(renderable Renderable, screen *ebiten.Image, boxtype collision.BoxType) {
	initWhitePixel()
	currentSprite := renderable.GetSprite()
	if currentSprite == nil {
		return
	}

	if boxes, ok := currentSprite.Boxes[boxtype]; ok {
		for _, box := range boxes {
			options := createBoxImageOptions(renderable, box, boxtype)
			screen.DrawImage(whitePixel, options)
		}
	}
}

func createBoxImageOptions(renderable Renderable, box types.Rect, boxType collision.BoxType) *ebiten.DrawImageOptions {
	boxImgOptions := &ebiten.DrawImageOptions{}

	position := calculateBoxScreenPosition(renderable, box)
	boxImgOptions.GeoM.Scale(box.W, box.H)
	boxImgOptions.GeoM.Translate(position.X, position.Y)

	// Set the color based on box type using ColorScale
	if color, exists := boxColors[boxType]; exists {
		boxImgOptions.ColorScale.ScaleWithColor(color)
	}

	return boxImgOptions
}

func createBoxImageOptionsWithCamera(renderable Renderable, box types.Rect, boxType collision.BoxType, camera *Camera) *ebiten.DrawImageOptions {
	boxImgOptions := &ebiten.DrawImageOptions{}

	boxImgOptions.GeoM.Scale(box.W, box.H)

	worldPos := renderable.GetPosition()
	boxWorldPos := types.Vector2{
		X: worldPos.X + box.X,
		Y: worldPos.Y + box.Y,
	}

	screenPos := camera.WorldToScreen(boxWorldPos)

	boxRenderable := &boxRenderableWrapper{
		position:         boxWorldPos,
		renderProperties: DefaultRenderProperties(),
	}

	applyCameraTransform(boxImgOptions, camera, boxRenderable, screenPos)

	if color, exists := boxColors[boxType]; exists {
		boxImgOptions.ColorScale.ScaleWithColor(color)
	}

	return boxImgOptions
}

// boxRenderableWrapper is a helper struct for box rendering with zoomAroundCenterOption
type boxRenderableWrapper struct {
	position         types.Vector2
	renderProperties RenderProperties
}

func (b *boxRenderableWrapper) GetPosition() types.Vector2 {
	return b.position
}

func (b *boxRenderableWrapper) GetSprite() *animation.Sprite {
	return nil // Boxes don't need sprites
}

func (b *boxRenderableWrapper) GetRenderProperties() RenderProperties {
	return b.renderProperties
}

func calculateBoxScreenPosition(renderable Renderable, box types.Rect) types.Vector2 {
	sprite := renderable.GetSprite()
	spriteScreenOriginX := float64(config.WindowWidth)/2 - (sprite.Rect.W / 2)
	spriteScreenOriginY := float64(config.WindowHeight)/2 - (sprite.Rect.H / 2)
	return types.Vector2{
		X: spriteScreenOriginX + box.X,
		Y: spriteScreenOriginY + box.Y,
	}
}
