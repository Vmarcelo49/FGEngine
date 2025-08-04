package graphics

import (
	"errors"
	"image/color"

	"FGEngine/character"
	"FGEngine/collision"
	"FGEngine/config"
	"FGEngine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	whitePixel *ebiten.Image
	boxColors  = map[collision.BoxType]color.RGBA{
		collision.Collision: {R: 80, G: 80, B: 80, A: 32},
		collision.Hit:       {R: 100, G: 40, B: 40, A: 32},
		collision.Hurt:      {R: 40, G: 100, B: 40, A: 32}}
)

// DrawBoxes draws all collision boxes for the current character's sprite on the screen.
// If the character is nil, it does nothing.
func DrawBoxes(character *character.Character, screen *ebiten.Image) error {
	if character == nil {
		return errors.New("cannot draw box: character is nil")
	}

	if character.CurrentSprite == nil {
		return errors.New("cannot draw box: character has no current sprite")
	}

	if whitePixel == nil {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	}

	for _, box := range character.GetAllBoxes() {
		options := createBoxImageOptions(character, box)
		screen.DrawImage(whitePixel, options)
	}
	return nil
}

// DrawBoxesByType draws boxes of a specific type, probably slower than DrawBoxes.
func DrawBoxesByType(character *character.Character, screen *ebiten.Image, boxtype collision.BoxType) error {
	if character == nil {
		return errors.New("cannot draw box: character is nil")
	}

	if character.CurrentSprite == nil {
		return errors.New("cannot draw box: character has no current sprite")
	}

	if whitePixel == nil {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	}

	switch boxtype {
	case collision.Collision:
		for _, boxRect := range character.CurrentSprite.CollisionBoxes {
			options := createBoxImageOptions(character, collision.Box{Rect: boxRect, BoxType: collision.Collision})
			screen.DrawImage(whitePixel, options)
		}
	case collision.Hit:
		for _, boxRect := range character.CurrentSprite.HitBoxes {
			options := createBoxImageOptions(character, collision.Box{Rect: boxRect, BoxType: collision.Hit})
			screen.DrawImage(whitePixel, options)
		}
	case collision.Hurt:
		for _, boxRect := range character.CurrentSprite.HurtBoxes {
			options := createBoxImageOptions(character, collision.Box{Rect: boxRect, BoxType: collision.Hurt})
			screen.DrawImage(whitePixel, options)
		}
	default:
		// If the box type is not recognized, we do nothing.
		return errors.New("cannot draw box: unrecognized box type")
	}
	return nil
}

func createBoxImageOptions(character *character.Character, box collision.Box) *ebiten.DrawImageOptions {
	boxImgOptions := &ebiten.DrawImageOptions{}

	zoom := config.GetZoom()
	position := calculateBoxScreenPosition(character, box, zoom)
	scale := calculateBoxScale(box, zoom)

	boxImgOptions.GeoM.Scale(scale.X, scale.Y) // X is width, Y is height
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

func calculateBoxScale(box collision.Box, zoom float64) types.Vector2 {
	scaleW := box.Rect.W * zoom
	scaleH := box.Rect.H * zoom

	if scaleW <= 0 {
		scaleW = 1 // Minimum scale
	}
	if scaleH <= 0 {
		scaleH = 1 // Minimum scale
	}

	return types.Vector2{X: scaleW, Y: scaleH}
}

func calculateBoxScreenPosition(character *character.Character, box collision.Box, zoom float64) types.Vector2 {
	sprite := character.CurrentSprite
	// Calculate sprite center on screen
	screenCenterX := float64(config.WindowWidth) / 2
	screenCenterY := float64(config.WindowHeight) / 2

	spriteScreenOriginX := screenCenterX - (sprite.SourceSize.W/2)*zoom
	spriteScreenOriginY := screenCenterY - (sprite.SourceSize.H/2)*zoom

	// Add box offset
	return types.Vector2{
		X: spriteScreenOriginX + box.Rect.X*zoom,
		Y: spriteScreenOriginY + box.Rect.Y*zoom,
	}
}
