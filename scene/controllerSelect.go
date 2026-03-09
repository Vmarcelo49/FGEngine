package scene

import (
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type controllerEntry struct {
	ID  ebiten.GamepadID
	Img *ebiten.Image
}

type ControllerScene struct {
	controllerIMGs []controllerEntry
}

func MakeControllerScene() Scene {
	cScene := &ControllerScene{
		controllerIMGs: make([]controllerEntry, 0, 2),
	}
	for _, id := range input.GamepadIDs {
		cScene.controllerIMGs = append(cScene.controllerIMGs, controllerEntry{ID: id, Img: graphics.LoadImage("assets/common/gamepad.png")})
	}
	cScene.controllerIMGs = append(cScene.controllerIMGs, controllerEntry{ID: ebiten.GamepadID(-1), Img: graphics.LoadImage("assets/common/keyboard.png")})
	return cScene
}

func (c *ControllerScene) Update(inputs [2]input.GameInput) SceneStatus {
	for _, id := range input.GamepadIDs {
		found := false
		for _, entry := range c.controllerIMGs {
			if entry.ID == id {
				found = true
				break
			}
		}
		if !found {
			img := graphics.LoadImage("assets/common/gamepad.png")
			if id == ebiten.GamepadID(-1) {
				img = graphics.LoadImage("assets/common/keyboard.png")
			}
			c.controllerIMGs = append(c.controllerIMGs, controllerEntry{ID: id, Img: img})
		}
	}
	for _, singleInput := range input.GlobalInputs {
		cur := singleInput.Buttons
		prev := singleInput.PrevButtons
		if singleInput.Owner == input.P1Side {
			if input.JustPressed(cur, prev, input.Right) {
				singleInput.Owner = input.UnAssigned
				continue
			}
			if input.JustPressed(cur, prev, input.A) {
				return Scene1
			}
		}
		if singleInput.Owner == input.P2Side {
			if input.JustPressed(cur, prev, input.Left) {
				singleInput.Owner = input.UnAssigned
				continue
			}
			if input.JustPressed(cur, prev, input.A) {
				return Scene1
			}
		}
		if singleInput.Owner == input.UnAssigned {
			if input.JustPressed(cur, prev, input.Left) {
				singleInput.Owner = input.P1Side
			}
			if input.JustPressed(cur, prev, input.Right) {
				singleInput.Owner = input.P2Side
			}
		}
	}

	return SceneDontChange
}

func (c *ControllerScene) Draw(screen *ebiten.Image) {
	p1Count, p2Count, unassignedCount := 0, 0, 0
	for _, entry := range c.controllerIMGs {
		op := &ebiten.DrawImageOptions{}
		img := entry.Img
		pos := input.UnAssigned
		for _, singleInput := range input.GlobalInputs {
			if singleInput.ID == entry.ID {
				pos = singleInput.Owner
				break
			}
		}
		imgH := float64(img.Bounds().Dy())
		spacing := imgH + 10
		switch pos {
		case input.P1Side:
			leftSidePos := constants.CameraWidth/2 - float64(img.Bounds().Dx())/2 - 150
			op.GeoM.Translate(leftSidePos, 100+spacing*float64(p1Count))
			p1Count++
		case input.P2Side:
			rightSidePos := constants.CameraWidth/2 - float64(img.Bounds().Dx())/2 + 150
			op.GeoM.Translate(rightSidePos, 100+spacing*float64(p2Count))
			p2Count++
		default:
			op.GeoM.Translate(constants.CameraWidth/2-float64(img.Bounds().Dx())/2, 100+spacing*float64(unassignedCount))
			unassignedCount++
		}
		screen.DrawImage(img, op)
	}
}
