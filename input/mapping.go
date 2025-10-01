package input

import "github.com/hajimehoshi/ebiten/v2"

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
