package input

import (
	"fgengine/config"
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
	Owner       ControllerPosition
	Buttons     GameInput
	PrevButtons GameInput
	ID          ebiten.GamepadID
	Mapping     InputMap
}

func GetPlayerInputs() [2]GameInput {
	inputs := [2]GameInput{NoInput, NoInput}
	for _, inpu := range GlobalInputs {
		if inpu.Owner == P1Side {
			inputs[0] |= inpu.Buttons
		}
		if inpu.Owner == P2Side {
			inputs[1] |= inpu.Buttons
		}
	}
	return inputs
}

// PollGamepads returns the combined GameInput for the specified gamepad IDs and the keyboard(if ID is -1). If no IDs are provided(nil is passed), it checks all connected gamepads.
func PollGamepads(ids []ebiten.GamepadID) GameInput {
	var localInputs GameInput
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
			if xValue > config.ControllerDeadzone {
				localInputs |= Right
			} else if xValue < -config.ControllerDeadzone {
				localInputs |= Left
			}

			// Left stick Y axis (axis 1)
			yValue := ebiten.GamepadAxisValue(gamepadID, 1)
			if yValue > config.ControllerDeadzone {
				localInputs |= Down
			} else if yValue < -config.ControllerDeadzone {
				localInputs |= Up
			}
		}
	}
	checkSOCD(&localInputs)
	return localInputs
}
