package ui

import (
	"fgengine/graphics"
	"fgengine/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type ControllerUI struct {
	State       *input.ControllerState
	gamepadImg  *ebiten.Image
	keyboardImg *ebiten.Image
}

func NewControllerUI(state *input.ControllerState) *ControllerUI {
	return &ControllerUI{
		State:       state,
		gamepadImg:  graphics.LoadImage("assets/common/gamepad.png"),
		keyboardImg: graphics.LoadImage("assets/common/keyboard.png"),
	}
}

func (c *ControllerUI) Draw(screen *ebiten.Image, camera *graphics.Camera) {
	baseY := camera.Viewport.H / 4

	positions := []input.ControllerPosition{input.P1Side, input.Center, input.P2Side}
	for _, pos := range positions {
		x := columnX(camera, pos, c.gamepadImg.Bounds().Dx())
		ids := c.State.ByPosition[pos]
		for i, id := range ids {
			img := c.gamepadImg
			if id == ebiten.GamepadID(-1) {
				img = c.keyboardImg
			}
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(x, baseY+float64(i*img.Bounds().Dy()))
			screen.DrawImage(img, opts)
		}
	}
}

func columnX(camera *graphics.Camera, pos input.ControllerPosition, iconWidth int) float64 {
	switch pos {
	case input.P1Side:
		return camera.Viewport.W/4 - float64(iconWidth/2)
	case input.Center:
		return camera.Viewport.W/2 - float64(iconWidth/2)
	case input.P2Side:
		return camera.Viewport.W*3/4 - float64(iconWidth/2)
	default:
		return 0
	}
}
