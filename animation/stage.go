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
		StageImage = ebiten.NewImage(int(constants.World.W), int(constants.World.H))
	}
	StageImage.Fill(color)

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(screenPos.X, screenPos.Y)

	screen.DrawImage(StageImage, options)
}

func DrawGridStage(gridPixels int, lineColor, bgColor color.RGBA, screen *ebiten.Image, screenPos types.Vector2) {
	if StageImage == nil {
		StageImage = ebiten.NewImage(int(constants.World.W), int(constants.World.H))
		StageImage.Fill(bgColor)
	}
	for x := 0; x < int(constants.World.W); x += gridPixels {
		for y := 0; y < int(constants.World.H); y += gridPixels {
			vector.StrokeLine(StageImage, float32(x), 0, float32(x), float32(constants.World.H), 1, lineColor, false)
			vector.StrokeLine(StageImage, 0, float32(y), float32(constants.World.W), float32(y), 1, lineColor, false)
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
