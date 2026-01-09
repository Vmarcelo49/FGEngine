package input

import (
	"fgengine/config"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputMap struct {
	KeyboardBindings map[GameInput][]ebiten.Key
	GamepadButtons   map[GameInput][]ebiten.StandardGamepadButton
}

type InputManager struct {
	InputMap   *InputMap
	GamepadIDs []ebiten.GamepadID
}

func NewInputManager() *InputManager {
	return &InputManager{
		InputMap:   NewDefaultInputMap(),
		GamepadIDs: []ebiten.GamepadID{},
	}
}

// LoadKeyboardBinding lets callers set per-player keyboard controls while keeping default gamepad buttons.
func LoadKeyboardBinding(bindings map[GameInput][]ebiten.Key) *InputManager {
	defaultPad := NewDefaultInputMap().GamepadButtons
	return &InputManager{
		InputMap: &InputMap{
			KeyboardBindings: bindings,
			GamepadButtons:   defaultPad,
		},
		GamepadIDs: []ebiten.GamepadID{},
	}
}

// Poll aggregates keyboard and assigned gamepad inputs for this player.
func (im *InputManager) Poll() GameInput {
	var localInputs GameInput
	if im.InputMap == nil {
		return NoInput
	}

	for gameInput, keys := range im.InputMap.KeyboardBindings {
		if slices.ContainsFunc(keys, ebiten.IsKeyPressed) {
			localInputs |= gameInput
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
			xValue := ebiten.GamepadAxisValue(gamepadID, 0)
			if xValue > config.ControllerDeadzone {
				localInputs |= Right
			} else if xValue < -config.ControllerDeadzone {
				localInputs |= Left
			}

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

func NewDefaultInputMap() *InputMap {
	return &InputMap{
		KeyboardBindings: map[GameInput][]ebiten.Key{
			Up:    {ebiten.KeyW, ebiten.KeySpace, ebiten.KeyUp},
			Down:  {ebiten.KeyS, ebiten.KeyDown},
			Left:  {ebiten.KeyA, ebiten.KeyLeft},
			Right: {ebiten.KeyD, ebiten.KeyRight},
			A:     {ebiten.KeyU},
			B:     {ebiten.KeyI},
			C:     {ebiten.KeyO},
			D:     {ebiten.KeyJ},
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
