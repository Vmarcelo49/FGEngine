package animation

import (
	"fgengine/constants"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TODO, move this to a proper structure to handle stages

var StageImage *ebiten.Image

func DrawStaticColorStage(color color.RGBA, screen *ebiten.Image) {
	if StageImage == nil {
		StageImage = ebiten.NewImage(int(constants.WorldWidth), int(constants.WorldHeight))
	}
	StageImage.Fill(color)
	screen.DrawImage(StageImage, nil)
}

func DrawGridStage(gridPixels int, lineColor, bgColor color.RGBA, screen *ebiten.Image) {
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
	screen.DrawImage(StageImage, nil)
}
