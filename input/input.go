package input

import (
	"fgengine/config"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameInput uint8

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

func (gi GameInput) IsPressed(input GameInput) bool {
	return gi&input != 0
}

type InputMap struct {
	KeyboardBindings map[GameInput][]ebiten.Key
	GamepadButtons   map[GameInput][]ebiten.StandardGamepadButton
}

func MakeDefaultInputMap() *InputMap {
	return &InputMap{
		KeyboardBindings: map[GameInput][]ebiten.Key{
			Up:    {ebiten.KeyW, ebiten.KeySpace, ebiten.KeyUp},
			Down:  {ebiten.KeyS, ebiten.KeyDown},
			Left:  {ebiten.KeyA, ebiten.KeyLeft},
			Right: {ebiten.KeyD, ebiten.KeyRight},
			A:     {ebiten.KeyU},
			B:     {ebiten.KeyI},
			C:     {ebiten.KeyO},
			D:     {ebiten.KeyK},
		},
		GamepadButtons: map[GameInput][]ebiten.StandardGamepadButton{
			Up:    {ebiten.StandardGamepadButtonLeftTop},
			Down:  {ebiten.StandardGamepadButtonLeftBottom},
			Left:  {ebiten.StandardGamepadButtonLeftLeft},
			Right: {ebiten.StandardGamepadButtonLeftRight},
			A:     {ebiten.StandardGamepadButtonRightLeft},
			B:     {ebiten.StandardGamepadButtonRightTop},
			C:     {ebiten.StandardGamepadButtonRightRight},
			D:     {ebiten.StandardGamepadButtonRightBottom},
		},
	}
}

type InputManager struct {
	InputMap   *InputMap
	GamepadIDs []ebiten.GamepadID
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
func (im *InputManager) GetLocalInputs() GameInput { // TODO, refactor to only check inputs of assigned gamepads and keyboard
	var localInputs GameInput

	for gameInput, keys := range im.InputMap.KeyboardBindings {
		if slices.ContainsFunc(keys, ebiten.IsKeyPressed) {
			localInputs |= gameInput // Once we find one pressed button for this input, we don't need to check the others
		}
	}

	for _, gamepadID := range im.GamepadIDs {
		for gameInput, buttons := range im.InputMap.GamepadButtons {
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
