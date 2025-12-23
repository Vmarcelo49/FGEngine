package input

import (
	"fgengine/constants"
	"fgengine/graphics"
	"fmt"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// All gamepadIDS, use only controllers inside players for the gameplay, the others for controlling menus
var GamepadIDs []ebiten.GamepadID

func CheckForGamePads() {
	connectedGamePads := inpututil.AppendJustConnectedGamepadIDs(nil)

	for _, id := range connectedGamePads {
		msg := fmt.Sprintf("Gamepad connected: ID: %d, Name: %s", id, ebiten.GamepadName(id))
		log.Print(msg)
		GamepadIDs = append(GamepadIDs, connectedGamePads...)
		log.Printf("There are %d controllers connected", len(GamepadIDs))
	}

	for i := len(GamepadIDs) - 1; i >= 0; i-- {
		if inpututil.IsGamepadJustDisconnected(GamepadIDs[i]) {
			log.Printf("Gamepad disconnected: ID: %d", GamepadIDs[i])
			GamepadIDs = slices.Delete(GamepadIDs, i, i+1)
			log.Printf("There are %d controllers connected", len(GamepadIDs))
		}
	}

}

type ControllerPosition int

const (
	P1Side ControllerPosition = iota
	Center
	P2Side
)

type ControllerConfig struct {
	position   ControllerPosition
	gamepadIDs []ebiten.GamepadID
}

type ControllerDrawable struct {
	controllers []ControllerConfig
}

func (c *ControllerDrawable) Draw(screen *ebiten.Image, camera *graphics.Camera) {
	controllerSelectionY := constants.Camera.H / 4
	gamepadImg := graphics.LoadImage("assets/common/gamepad.png")
	keyboardImg := graphics.LoadImage("assets/common/keyboard.png")

	var verticalSum float64

	for _, gameController := range c.controllers {
		x := getControllerXPosition(gameController, gamepadImg.Bounds().Dx())
		drawController(screen, gamepadImg, x, controllerSelectionY+verticalSum)
		verticalSum += float64(gamepadImg.Bounds().Dy())
	}
}

func getControllerXPosition(controller ControllerConfig, size int) int {
	switch controller.position {
	case P1Side:
		return int(constants.Camera.W/4) - size/2
	case Center:
		return int(constants.Camera.W/2) - size/2
	case P2Side:
		return int(constants.Camera.W*3/4) - size/2
	default:
		return 0
	}
}
