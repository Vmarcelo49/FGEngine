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
	checkGamepadConnections()
	// Rebuild GlobalInputs by device ID to preserve owners while keeping
	// device list in sync with the latest connected IDs.
	byID := make(map[ebiten.GamepadID]*Input, len(GlobalInputs))
	for _, in := range GlobalInputs {
		byID[in.ID] = in
	}

	syncedInputs := make([]*Input, 0, len(GamepadIDs))
	for _, id := range GamepadIDs {
		if existing, ok := byID[id]; ok {
			syncedInputs = append(syncedInputs, existing)
			continue
		}
		syncedInputs = append(syncedInputs, &Input{
			ID:      id,
			Mapping: *NewDefaultInputMap(),
		})
	}
	GlobalInputs = syncedInputs

	p1IDs := make([]ebiten.GamepadID, 0, len(GlobalInputs))
	p2IDs := make([]ebiten.GamepadID, 0, len(GlobalInputs))
	for _, i := range GlobalInputs {
		i.PrevButtons = i.ActiveButtons
		i.ActiveButtons = PollGamepads([]ebiten.GamepadID{i.ID})
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

func checkGamepadConnections() {
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
}
