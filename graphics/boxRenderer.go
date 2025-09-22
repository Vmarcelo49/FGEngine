package graphics

import (
	"image/color"
	"sync"

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

	// Object pool for DrawImageOptions to reduce allocations
	drawOptionsPool = sync.Pool{
		New: func() any {
			return &ebiten.DrawImageOptions{}
		},
	}

	// Cached screen center calculations
	screenCenterX = float64(config.WindowWidth) / 2 // TODO, check if this should be replaced with world coordinates
	screenCenterY = float64(config.WindowHeight) / 2
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

	for boxType, boxes := range currentSprite.Boxes {
		for _, box := range boxes {
			options := createBoxImageOptionsWithCamera(renderable, box, boxType, camera)
			screen.DrawImage(whitePixel, options)
			// Return the options to the pool after use
			drawOptionsPool.Put(options)
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
			// Return the options to the pool after use
			drawOptionsPool.Put(options)
		}
	}
}

func createBoxImageOptions(renderable Renderable, box types.Rect, boxType collision.BoxType) *ebiten.DrawImageOptions {
	boxImgOptions := drawOptionsPool.Get().(*ebiten.DrawImageOptions)

	// Reset options to clean state
	boxImgOptions.GeoM.Reset()
	boxImgOptions.ColorScale.Reset()

	position := calculateBoxScreenPosition(renderable, box)
	boxImgOptions.GeoM.Translate(position.X, position.Y)

	// Set the color based on box type using ColorScale
	if color, exists := boxColors[boxType]; exists {
		boxImgOptions.ColorScale.ScaleWithColor(color)
	}

	return boxImgOptions
}

func createBoxImageOptionsWithCamera(renderable Renderable, box types.Rect, boxType collision.BoxType, camera *Camera) *ebiten.DrawImageOptions {
	boxImgOptions := drawOptionsPool.Get().(*ebiten.DrawImageOptions)

	// Reset options to clean state
	boxImgOptions.GeoM.Reset()
	boxImgOptions.ColorScale.Reset()

	// Calculate world position of the box
	worldPos := renderable.GetPosition()
	boxWorldPos := types.Vector2{
		X: worldPos.X + box.X,
		Y: worldPos.Y + box.Y,
	}

	// Transform to screen coordinates using camera
	screenPos := camera.WorldToScreen(boxWorldPos)
	boxImgOptions.GeoM.Scale(box.W, box.H)
	boxImgOptions.GeoM.Scale(camera.Scaling, camera.Scaling)
	boxImgOptions.GeoM.Translate(screenPos.X, screenPos.Y)

	// Set the color based on box type using ColorScale
	if color, exists := boxColors[boxType]; exists {
		boxImgOptions.ColorScale.ScaleWithColor(color)
	}

	return boxImgOptions
}

func calculateBoxScreenPosition(renderable Renderable, box types.Rect) types.Vector2 {
	sprite := renderable.GetSprite()
	spriteScreenOriginX := screenCenterX - (sprite.Rect.W / 2)
	spriteScreenOriginY := screenCenterY - (sprite.Rect.H / 2)
	return types.Vector2{
		X: spriteScreenOriginX + box.X,
		Y: spriteScreenOriginY + box.Y,
	}
}
