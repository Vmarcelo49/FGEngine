package scene

import (
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type ControllerScene struct {
	controllerIMGs map[ebiten.GamepadID]*ebiten.Image
}

func MakeControllerScene() Scene {
	cScene := &ControllerScene{
		controllerIMGs: make(map[ebiten.GamepadID]*ebiten.Image),
	}
	for _, id := range input.GamepadIDs {
		cScene.controllerIMGs[id] = graphics.LoadImage("assets/common/gamepad.png")
	}
	cScene.controllerIMGs[ebiten.GamepadID(-1)] = graphics.LoadImage("assets/common/keyboard.png")
	return cScene
}

func (c *ControllerScene) Update(inputs [2]input.GameInput) SceneStatus {
	input.UpdateGamepads()
	combinedInputs := input.CombinedInputs()

	for _, singleInput := range combinedInputs {
		if singleInput.Owner == input.P1Side {
			if singleInput.Buttons.IsPressed(input.Right) {
				singleInput.Owner = input.UnAssigned
				continue
			}
			if singleInput.Buttons.IsPressed(input.A) {
				return Scene1
			}
		}
		if singleInput.Owner == input.P2Side {
			if singleInput.Buttons.IsPressed(input.Left) {
				singleInput.Owner = input.UnAssigned
				continue
			}
			if singleInput.Buttons.IsPressed(input.A) {
				return Scene1
			}
		}
		if singleInput.Owner == input.UnAssigned {
			if singleInput.Buttons.IsPressed(input.Left) {
				singleInput.Owner = input.P1Side
			}
			if singleInput.Buttons.IsPressed(input.Right) {
				singleInput.Owner = input.P2Side
			}
		}
	}

	return SceneDontChange
}

func (c *ControllerScene) Draw(screen *ebiten.Image) {
	for id, img := range c.controllerIMGs {
		op := &ebiten.DrawImageOptions{}
		pos := input.UnAssigned
		for _, singleInput := range input.CombinedInputs() {
			if singleInput.ID == id {
				pos = singleInput.Owner
				break
			}
		}
		switch pos {
		case input.P1Side:
			leftSidePos := constants.CameraWidth/2 - float64(img.Bounds().Dx())/2 - 100
			op.GeoM.Translate(leftSidePos, 100)
		case input.P2Side:
			rightSidePos := constants.CameraWidth/2 + float64(img.Bounds().Dx())/2 + 100
			op.GeoM.Translate(rightSidePos, 100)
		default:
			op.GeoM.Translate(constants.CameraWidth/2-float64(img.Bounds().Dx())/2, 300)
		}
		screen.DrawImage(img, op)
	}
}
