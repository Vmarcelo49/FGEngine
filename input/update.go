package input

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var GamepadIDs []ebiten.GamepadID // all connected gamepads, updated by UpdateGamepads()

var GlobalInputs []*Input

// UpdateGamepads checks for newly connected or disconnected gamepads and updates the GamepadIDs slice accordingly. It also logs these events.
func UpdateGamepads() [2]GameInput {
	for _, id := range inpututil.AppendJustConnectedGamepadIDs(nil) {
		log.Printf("Gamepad connected: ID: %d, Name: %s", id, ebiten.GamepadName(id))
	}
	for _, id := range GamepadIDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			log.Printf("Gamepad disconnected: ID: %d", id)
		}
	}
	GamepadIDs = ebiten.AppendGamepadIDs(GamepadIDs[:0])
	GamepadIDs = append(GamepadIDs, ebiten.GamepadID(-1)) // Add -1 for keyboard input

	if len(GlobalInputs) < len(GamepadIDs) {
		for i := len(GlobalInputs); i < len(GamepadIDs); i++ {
			GlobalInputs = append(GlobalInputs, &Input{
				ID:      GamepadIDs[i],
				Mapping: *NewDefaultInputMap(),
			})
		}
	} else if len(GlobalInputs) > len(GamepadIDs) {
		GlobalInputs = GlobalInputs[:len(GamepadIDs)]
	}
	var p1IDs []ebiten.GamepadID
	var p2IDs []ebiten.GamepadID
	for _, i := range GlobalInputs {
		i.PrevButtons = i.Buttons
		i.Buttons = PollGamepads([]ebiten.GamepadID{i.ID})
		if i.Owner == P1Side {
			p1IDs = append(p1IDs, i.ID)
		}
		if i.Owner == P2Side {
			p2IDs = append(p2IDs, i.ID)
		}
	}
	inputs := [2]GameInput{PollGamepads(p1IDs), PollGamepads(p2IDs)}
	return inputs
}
