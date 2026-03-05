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
	Owner   ControllerPosition
	Buttons GameInput
	ID      ebiten.GamepadID
	Mapping InputMap
}

type GameInput byte

const (
	NoInput GameInput = 0
)

const (
	Up GameInput = 1 << iota
	Down
	Left
	Right
	A
	B
	C
	D
)

func (gi GameInput) String() string {
	if gi == NoInput {
		return "NoInput"
	}
	str := ""
	if gi&Up != 0 {
		str += "Up "
	}
	if gi&Down != 0 {
		str += "Down "
	}
	if gi&Left != 0 {
		str += "Left "
	}
	if gi&Right != 0 {
		str += "Right "
	}
	if gi&A != 0 {
		str += "A "
	}
	if gi&B != 0 {
		str += "B "
	}
	if gi&C != 0 {
		str += "C "
	}
	if gi&D != 0 {
		str += "D "
	}
	return str
}

func (gi GameInput) IsPressed(input GameInput) bool {
	return gi&input != 0
}

func (im *InputManager) UpdateGamepadList() {
	im.GamepadIDs = ebiten.AppendGamepadIDs(im.GamepadIDs[:0])
}

func CombinedInputs() []Input {
	var inputs []Input
	for _, id := range GamepadIDs {
		inputs = append(inputs, Input{
			Owner:   UnAssigned,
			Buttons: PollGamepads([]ebiten.GamepadID{id}),
			ID:      id,
			Mapping: *NewDefaultInputMap(),
		})
	}
	inputs = append(inputs, Input{
		Owner:   UnAssigned,
		Buttons: PollGamepads([]ebiten.GamepadID{-1}), // Poll all gamepads + keyboard
		ID:      -1,                                   // -1 for keyboard
		Mapping: *NewDefaultInputMap(),
	})
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

func (im *InputManager) AssignGamepadID(id ebiten.GamepadID) {
	im.GamepadIDs = append(im.GamepadIDs, id)
}

func checkSOCD(input *GameInput) {
	if input.IsPressed(Left) && input.IsPressed(Right) {
		*input &^= (Left | Right)
	}

	if input.IsPressed(Up) && input.IsPressed(Down) {
		*input &^= (Up | Down)
	}
}
