package device

import (
	"fgengine/config"
	"fgengine/input"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type ControllerPosition int

const (
	UnAssigned ControllerPosition = iota
	P1Side
	P2Side
)

type Input struct {
	Owner         ControllerPosition
	ActiveButtons input.GameInput
	PrevButtons   input.GameInput
	ID            ebiten.GamepadID
	Mapping       InputMap
}

func GetPlayerInputs() [2]input.GameInput {
	inputs := [2]input.GameInput{input.NoInput, input.NoInput}
	for _, inpu := range GlobalInputs {
		if inpu.Owner == P1Side {
			inputs[0] |= inpu.ActiveButtons
		}
		if inpu.Owner == P2Side {
			inputs[1] |= inpu.ActiveButtons
		}
	}
	return inputs
}

// PollGamepads returns the combined GameInput for the specified gamepad IDs and the keyboard(if ID is -1). If no IDs are provided(nil is passed), it checks all connected gamepads.
// Raw per-device result — SOCD filtering applies after the per-player merge in UpdateGamepads (SPEC §5.2).
func PollGamepads(ids []ebiten.GamepadID) input.GameInput {
	var localInputs input.GameInput
	inputmap := NewDefaultInputMap()

	// If nil is passed, check all connected gamepads
	pollIDs := ids
	pollKeyboard := false
	if ids == nil {
		pollIDs = GamepadIDs
		pollKeyboard = true
	} else {
		// Check if -1 (keyboard) is among the requested IDs
		pollKeyboard = slices.Contains(ids, ebiten.GamepadID(-1))
	}

	if pollKeyboard {
		for gameInput, keys := range inputmap.KeyboardBindings {
			if slices.ContainsFunc(keys, ebiten.IsKeyPressed) {
				localInputs |= gameInput
			}
		}
	}

	for _, gamepadID := range pollIDs {
		if gamepadID == ebiten.GamepadID(-1) {
			continue // Skip keyboard marker in gamepad polling
		}
		for gameInput, buttons := range inputmap.GamepadButtons {
			for _, button := range buttons {
				if ebiten.IsStandardGamepadButtonPressed(gamepadID, button) {
					localInputs |= gameInput
					break
				}
			}
		}
		axisCount := ebiten.GamepadAxisCount(gamepadID)
		if axisCount >= 2 {
			// Left stick X axis (axis 0)
			xValue := ebiten.GamepadAxisValue(gamepadID, 0)
			if xValue > config.ActiveConfig.ControllerDeadzone {
				localInputs |= input.Right
			} else if xValue < -config.ActiveConfig.ControllerDeadzone {
				localInputs |= input.Left
			}

			// Left stick Y axis (axis 1)
			yValue := ebiten.GamepadAxisValue(gamepadID, 1)
			if yValue > config.ActiveConfig.ControllerDeadzone {
				localInputs |= input.Down
			} else if yValue < -config.ActiveConfig.ControllerDeadzone {
				localInputs |= input.Up
			}
		}
	}
	return localInputs
}

type keyboardBinding struct {
	gameInput input.GameInput
	keys      []ebiten.Key
}

var localKeyboardBindings = [2][]keyboardBinding{
	{
		{gameInput: input.Up, keys: []ebiten.Key{ebiten.KeyW}},
		{gameInput: input.Down, keys: []ebiten.Key{ebiten.KeyS}},
		{gameInput: input.Left, keys: []ebiten.Key{ebiten.KeyA}},
		{gameInput: input.Right, keys: []ebiten.Key{ebiten.KeyD}},
		{gameInput: input.A, keys: []ebiten.Key{ebiten.KeyU}},
		{gameInput: input.B, keys: []ebiten.Key{ebiten.KeyI}},
		{gameInput: input.C, keys: []ebiten.Key{ebiten.KeyO}},
		{gameInput: input.D, keys: []ebiten.Key{ebiten.KeyJ}},
	},
	{
		{gameInput: input.Up, keys: []ebiten.Key{ebiten.KeyUp}},
		{gameInput: input.Down, keys: []ebiten.Key{ebiten.KeyDown}},
		{gameInput: input.Left, keys: []ebiten.Key{ebiten.KeyLeft}},
		{gameInput: input.Right, keys: []ebiten.Key{ebiten.KeyRight}},
		{gameInput: input.A, keys: []ebiten.Key{ebiten.KeyN}},
		{gameInput: input.B, keys: []ebiten.Key{ebiten.KeyM}},
		{gameInput: input.C, keys: []ebiten.Key{ebiten.KeyComma}},
		{gameInput: input.D, keys: []ebiten.Key{ebiten.KeyPeriod}},
	},
}

// PollKeyboardPlayers returns independent P1 and P2 layouts for local play:
// WASD + U/I/O/J and arrows + N/M/comma/period.
func PollKeyboardPlayers() [2]input.GameInput {
	var players [2]input.GameInput
	for player, bindings := range localKeyboardBindings {
		for _, binding := range bindings {
			if slices.ContainsFunc(binding.keys, ebiten.IsKeyPressed) {
				players[player] |= binding.gameInput
			}
		}
		players[player] = input.ApplySOCD(players[player])
	}
	return players
}
