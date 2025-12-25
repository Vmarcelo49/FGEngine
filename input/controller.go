package input

import (
	"fgengine/graphics"
	"fmt"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// All gamepad IDs; use only controllers inside players for gameplay,
// others may control menus.
var GamepadIDs []ebiten.GamepadID

// checkForGamepadsConnections polls for just-connected and just-disconnected gamepads and
// updates the global GamepadIDs accordingly.
func checkForGamepadsConnections() {
	connectedGamepads := inpututil.AppendJustConnectedGamepadIDs(nil)

	for _, id := range connectedGamepads {
		msg := fmt.Sprintf("Gamepad connected: ID: %d, Name: %s", id, ebiten.GamepadName(id))
		log.Print(msg)
	}
	GamepadIDs = append(GamepadIDs, connectedGamepads...)
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

type ControllerDrawable struct {
	byPosition map[ControllerPosition][]ebiten.GamepadID
	positionOf map[ebiten.GamepadID]ControllerPosition
}

func newControllerDrawable() *ControllerDrawable {
	cd := &ControllerDrawable{
		byPosition: make(map[ControllerPosition][]ebiten.GamepadID),
		positionOf: make(map[ebiten.GamepadID]ControllerPosition),
	}
	// Initialize with current gamepads in Center
	for _, id := range GamepadIDs {
		cd.assign(id, Center)
	}
	// Add keyboard (-1) to Center
	cd.assign(ebiten.GamepadID(-1), Center)
	return cd
}

func (c *ControllerDrawable) removeFrom(pos ControllerPosition, id ebiten.GamepadID) {
	list := c.byPosition[pos]
	for i, v := range list {
		if v == id {
			c.byPosition[pos] = append(list[:i], list[i+1:]...)
			return
		}
	}
}

func (c *ControllerDrawable) assign(id ebiten.GamepadID, pos ControllerPosition) {
	if old, ok := c.positionOf[id]; ok {
		if old == pos {
			// Already in desired position; ensure no duplicates
			c.removeFrom(pos, id)
		} else {
			c.removeFrom(old, id)
		}
	}
	c.positionOf[id] = pos
	c.byPosition[pos] = append(c.byPosition[pos], id)
}

func (c *ControllerDrawable) move(id ebiten.GamepadID, dir int) {
	order := []ControllerPosition{P1Side, Center, P2Side}
	cur, ok := c.positionOf[id]
	if !ok {
		cur = Center
	}
	idx := 0
	for i, p := range order {
		if p == cur {
			idx = i
			break
		}
	}
	if dir < 0 && idx > 0 {
		idx--
	} else if dir > 0 && idx < len(order)-1 {
		idx++
	}
	c.assign(id, order[idx])
}

func (c *ControllerDrawable) update() {
	// Build current IDs (gamepads + keyboard)
	tempIDs := append([]ebiten.GamepadID{}, GamepadIDs...)
	tempIDs = append(tempIDs, ebiten.GamepadID(-1)) // -1 is the keyboard

	// Ensure newly seen IDs have a default position
	for _, id := range tempIDs {
		if _, ok := c.positionOf[id]; !ok {
			c.assign(id, Center)
		}
	}

	// Handle per-id left/right movement
	for _, id := range tempIDs {
		activeInput := LocalInputsFromIDS([]ebiten.GamepadID{id})
		left := activeInput.IsPressed(Left)
		right := activeInput.IsPressed(Right)
		if left && !right {
			c.move(id, -1)
		} else if right && !left {
			c.move(id, 1)
		}
	}
}

func (c *ControllerDrawable) Draw(screen *ebiten.Image, camera *graphics.Camera) { // camera used for layout and alignment
	baseY := camera.Viewport.H / 4
	gamepadImg := graphics.LoadImage("assets/common/gamepad.png")
	keyboardImg := graphics.LoadImage("assets/common/keyboard.png")

	positions := []ControllerPosition{P1Side, Center, P2Side}
	for _, pos := range positions {
		// Use gamepad width for column alignment
		x := columnX(camera, pos, gamepadImg.Bounds().Dx())
		ids := c.byPosition[pos]
		for i, id := range ids {
			img := gamepadImg
			if id == ebiten.GamepadID(-1) {
				img = keyboardImg
			}
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(x, baseY+float64(i*img.Bounds().Dy()))
			screen.DrawImage(img, opts)
		}
	}
}

func columnX(camera *graphics.Camera, pos ControllerPosition, iconWidth int) float64 {
	switch pos {
	case P1Side:
		return camera.Viewport.W/4 - float64(iconWidth/2)
	case Center:
		return camera.Viewport.W/2 - float64(iconWidth/2)
	case P2Side:
		return camera.Viewport.W*3/4 - float64(iconWidth/2)
	default:
		return 0
	}
}
