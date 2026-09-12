package device

import (
	"fgengine/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputMap struct {
	KeyboardBindings map[input.GameInput][]ebiten.Key
	GamepadButtons   map[input.GameInput][]ebiten.StandardGamepadButton
}

func NewDefaultInputMap() *InputMap {
	return &InputMap{
		KeyboardBindings: map[input.GameInput][]ebiten.Key{
			input.Up:    {ebiten.KeyW, ebiten.KeySpace, ebiten.KeyUp},
			input.Down:  {ebiten.KeyS, ebiten.KeyDown},
			input.Left:  {ebiten.KeyA, ebiten.KeyLeft},
			input.Right: {ebiten.KeyD, ebiten.KeyRight},
			input.A:     {ebiten.KeyU},
			input.B:     {ebiten.KeyI},
			input.C:     {ebiten.KeyO},
			input.D:     {ebiten.KeyJ},
		},
		GamepadButtons: map[input.GameInput][]ebiten.StandardGamepadButton{
			input.Up:    {ebiten.StandardGamepadButtonLeftTop},
			input.Down:  {ebiten.StandardGamepadButtonLeftBottom},
			input.Left:  {ebiten.StandardGamepadButtonLeftLeft},
			input.Right: {ebiten.StandardGamepadButtonLeftRight},
			input.A:     {ebiten.StandardGamepadButtonRightLeft},
			input.B:     {ebiten.StandardGamepadButtonRightTop},
			input.C:     {ebiten.StandardGamepadButtonRightRight},
			input.D:     {ebiten.StandardGamepadButtonRightBottom},
		},
	}
}
