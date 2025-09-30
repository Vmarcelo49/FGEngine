package animation

import (
	"fgengine/constants"
	"fgengine/types"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var StageImage *ebiten.Image

func DrawStaticColorStage(color color.RGBA, screen *ebiten.Image, screenPos types.Vector2) {
	if StageImage == nil {
		StageImage = ebiten.NewImage(int(constants.WorldWidth), int(constants.WorldHeight))
	}
	StageImage.Fill(color)

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(screenPos.X, screenPos.Y)

	screen.DrawImage(StageImage, options)
}

func DrawGridStage(gridPixels int, lineColor, bgColor color.RGBA, screen *ebiten.Image, screenPos types.Vector2) {
	if StageImage == nil {
		StageImage = ebiten.NewImage(int(constants.WorldWidth), int(constants.WorldHeight))
		StageImage.Fill(bgColor)
	}
	for x := 0; x < int(constants.WorldWidth); x += gridPixels {
		for y := 0; y < int(constants.WorldHeight); y += gridPixels {
			vector.StrokeLine(StageImage, float32(x), 0, float32(x), float32(constants.WorldHeight), 1, lineColor, false)
			vector.StrokeLine(StageImage, 0, float32(y), float32(constants.WorldWidth), float32(y), 1, lineColor, false)
		}
	}

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(screenPos.X, screenPos.Y)

	screen.DrawImage(StageImage, options)
}

func DrawStaticImageStage(img *ebiten.Image, screen *ebiten.Image, screenPos types.Vector2) {
	if img == nil {
		return
	}

	options := &ebiten.DrawImageOptions{}

	// screenPos is already calculated using WorldToScreen() from camera
	options.GeoM.Translate(screenPos.X, screenPos.Y)

	screen.DrawImage(img, options)
}
