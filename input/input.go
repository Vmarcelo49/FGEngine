package input

import (
	"FGEngine/config"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameInput uint8

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

type InputMap struct { // TODO, check if we want to make a slice of buttons for each input type, then we can have more flexibility
	KeyboardBindings map[GameInput]ebiten.Key
	GamepadButtons   map[GameInput]ebiten.StandardGamepadButton
}

type InputManager struct {
	InputMap   *InputMap
	GamepadIDs []ebiten.GamepadID
}

func MakeDefaultInputMap() *InputMap {
	return &InputMap{
		KeyboardBindings: map[GameInput]ebiten.Key{
			Up:    ebiten.KeyW, // ebiten.KeySpace,
			Down:  ebiten.KeyS,
			Left:  ebiten.KeyA,
			Right: ebiten.KeyD,
			A:     ebiten.KeyU,
			B:     ebiten.KeyI,
			C:     ebiten.KeyO,
			D:     ebiten.KeyK,
		},
		GamepadButtons: map[GameInput]ebiten.StandardGamepadButton{
			Up:    ebiten.StandardGamepadButtonLeftTop,
			Down:  ebiten.StandardGamepadButtonLeftBottom,
			Left:  ebiten.StandardGamepadButtonLeftLeft,
			Right: ebiten.StandardGamepadButtonLeftRight,
			A:     ebiten.StandardGamepadButtonRightLeft,
			B:     ebiten.StandardGamepadButtonRightTop,
			C:     ebiten.StandardGamepadButtonRightRight,
			D:     ebiten.StandardGamepadButtonRightBottom,
		},
	}
}

func NewInputManager() *InputManager {
	return &InputManager{
		InputMap:   MakeDefaultInputMap(),
		GamepadIDs: []ebiten.GamepadID{},
	}
}

func (im *InputManager) UpdateGamepadList() {
	im.GamepadIDs = ebiten.AppendGamepadIDs(im.GamepadIDs[:0])
}

// GetLocalInputs retrieves the current local input state.
func (im *InputManager) GetLocalInputs() GameInput {
	var localInputs GameInput

	for gameInput, key := range im.InputMap.KeyboardBindings {
		if ebiten.IsKeyPressed(key) {
			localInputs |= gameInput
		}
	}

	for _, gamepadID := range im.GamepadIDs {
		for gameInput, button := range im.InputMap.GamepadButtons {
			if ebiten.IsStandardGamepadButtonPressed(gamepadID, button) {
				localInputs |= gameInput
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

	return localInputs
}

func (gi GameInput) IsPressed(input GameInput) bool {
	return gi&input != 0
}
