package graphics

import (
	"image/color"
	"sync"

	"FGEngine/animation"
	"FGEngine/collision"
	"FGEngine/config"
	"FGEngine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	whitePixel *ebiten.Image
	once       sync.Once
	boxColors  = map[collision.BoxType]color.RGBA{
		collision.Collision: {R: 80, G: 80, B: 80, A: 32},
		collision.Hit:       {R: 100, G: 40, B: 40, A: 32},
		collision.Hurt:      {R: 40, G: 100, B: 40, A: 32}}
)

// whitePixel should never be unloaded
func initWhitePixel() {
	once.Do(func() {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	})
}

// DrawBoxes draws all collision boxes for the current renderable entity's sprite on the screen.
// If entity or sprite data is invalid, the function returns early without drawing.
func DrawBoxes(renderable animation.Renderable, screen *ebiten.Image) {
	if checkDrawConditions(renderable) == false {
		return
	}

	initWhitePixel()

	for _, box := range renderable.GetAllBoxes() {
		options := createBoxImageOptions(renderable, box)
		screen.DrawImage(whitePixel, options)
	}
}

// DrawBoxesByType draws boxes of a specific type.
// If entity, sprite data, or box type is invalid, the function returns early without drawing.
func DrawBoxesByType(renderable animation.Renderable, screen *ebiten.Image, boxtype collision.BoxType) {
	if checkDrawConditions(renderable) == false {
		return
	}

	initWhitePixel()
	currentSprite := renderable.GetAnimationComponent().GetCurrentSprite()

	switch boxtype {
	case collision.Collision:
		for _, boxRect := range currentSprite.CollisionBoxes {
			options := createBoxImageOptions(renderable, collision.Box{Rect: boxRect, BoxType: collision.Collision})
			screen.DrawImage(whitePixel, options)
		}
	case collision.Hit:
		for _, boxRect := range currentSprite.HitBoxes {
			options := createBoxImageOptions(renderable, collision.Box{Rect: boxRect, BoxType: collision.Hit})
			screen.DrawImage(whitePixel, options)
		}
	case collision.Hurt:
		for _, boxRect := range currentSprite.HurtBoxes {
			options := createBoxImageOptions(renderable, collision.Box{Rect: boxRect, BoxType: collision.Hurt})
			screen.DrawImage(whitePixel, options)
		}
	default:
		return
	}
}

func createBoxImageOptions(renderable animation.Renderable, box collision.Box) *ebiten.DrawImageOptions {
	boxImgOptions := &ebiten.DrawImageOptions{}

	position := calculateBoxScreenPosition(renderable, box)

	boxImgOptions.GeoM.Translate(position.X, position.Y)

	// Set the color based on box type using ColorScale
	switch box.BoxType {
	case collision.Collision:
		boxImgOptions.ColorScale.ScaleWithColor(boxColors[collision.Collision])
	case collision.Hit:
		boxImgOptions.ColorScale.ScaleWithColor(boxColors[collision.Hit])
	case collision.Hurt:
		boxImgOptions.ColorScale.ScaleWithColor(boxColors[collision.Hurt])
	}

	return boxImgOptions
}

func calculateBoxScreenPosition(renderable animation.Renderable, box collision.Box) types.Vector2 {
	sprite := renderable.GetAnimationComponent().GetCurrentSprite()
	screenCenterX := float64(config.WindowWidth) / 2
	screenCenterY := float64(config.WindowHeight) / 2

	spriteScreenOriginX := screenCenterX - (sprite.SourceSize.W / 2)
	spriteScreenOriginY := screenCenterY - (sprite.SourceSize.H / 2)

	return types.Vector2{
		X: spriteScreenOriginX + box.Rect.X,
		Y: spriteScreenOriginY + box.Rect.Y,
	}
}
