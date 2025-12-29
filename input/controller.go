package input

import (
	"fmt"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var GamepadIDs []ebiten.GamepadID

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

type ControllerState struct {
	ByPosition map[ControllerPosition][]ebiten.GamepadID
	PositionOf map[ebiten.GamepadID]ControllerPosition
}

func newControllerState() *ControllerState {
	cd := &ControllerState{
		ByPosition: make(map[ControllerPosition][]ebiten.GamepadID),
		PositionOf: make(map[ebiten.GamepadID]ControllerPosition),
	}
	// Initialize with current gamepads in Center
	for _, id := range GamepadIDs {
		cd.assign(id, Center)
	}
	// Add keyboard (-1) to Center
	cd.assign(ebiten.GamepadID(-1), Center)
	return cd
}

func (c *ControllerState) removeFrom(pos ControllerPosition, id ebiten.GamepadID) {
	list := c.ByPosition[pos]
	for i, v := range list {
		if v == id {
			c.ByPosition[pos] = append(list[:i], list[i+1:]...)
			return
		}
	}
}

func (c *ControllerState) assign(id ebiten.GamepadID, pos ControllerPosition) {
	if old, ok := c.PositionOf[id]; ok {
		if old == pos {
			// Already in desired position; ensure no duplicates
			c.removeFrom(pos, id)
		} else {
			c.removeFrom(old, id)
		}
	}
	c.PositionOf[id] = pos
	c.ByPosition[pos] = append(c.ByPosition[pos], id)
}

func (c *ControllerState) move(id ebiten.GamepadID, dir int) {
	order := []ControllerPosition{P1Side, Center, P2Side}
	cur, ok := c.PositionOf[id]
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

func (c *ControllerState) update() {
	// Build current IDs (gamepads + keyboard)
	tempIDs := append([]ebiten.GamepadID{}, GamepadIDs...)
	tempIDs = append(tempIDs, ebiten.GamepadID(-1)) // -1 is the keyboard

	// Ensure newly seen IDs have a default position
	for _, id := range tempIDs {
		if _, ok := c.PositionOf[id]; !ok {
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
