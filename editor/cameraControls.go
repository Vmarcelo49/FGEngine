package editor

import (
	"fgengine/input"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CameraMoveSpeed = 5.0
)

func (g *Game) handleCameraInput() {
	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if !g.mouse.isDragging {
			// Start dragging
			g.mouse.isDragging = true
			g.mouse.lastMouseX = x
			g.mouse.lastMouseY = y
		} else {
			deltaX := float64(g.mouse.lastMouseX - x)
			deltaY := float64(g.mouse.lastMouseY - y)

			// Adjust delta for camera scaling
			if g.camera.Scaling != 0 && g.camera.Scaling != 1 {
				deltaX /= g.camera.Scaling
				deltaY /= g.camera.Scaling
			}

			// mouse movement to camera (inverted because we want to "grab" the world)
			g.camera.SetPosition(types.Vector2{
				X: g.camera.Viewport.X + deltaX,
				Y: g.camera.Viewport.Y + deltaY,
			})

			g.mouse.lastMouseX = x
			g.mouse.lastMouseY = y
		}
	} else {
		g.mouse.isDragging = false
	}

	// keyboard input for camera movement
	inputs := input.LocalInputsFromIDS([]ebiten.GamepadID{ebiten.GamepadID(-1)})

	var cameraMove types.Vector2
	if inputs.IsPressed(input.Left) {
		cameraMove.X -= CameraMoveSpeed
	}
	if inputs.IsPressed(input.Right) {
		cameraMove.X += CameraMoveSpeed
	}
	if inputs.IsPressed(input.Up) {
		cameraMove.Y -= CameraMoveSpeed
	}
	if inputs.IsPressed(input.Down) {
		cameraMove.Y += CameraMoveSpeed
	}

	if cameraMove.X != 0 || cameraMove.Y != 0 {
		// Adjust camera movement speed based on scaling
		if g.camera.Scaling != 0 && g.camera.Scaling != 1 {
			cameraMove.X /= g.camera.Scaling
			cameraMove.Y /= g.camera.Scaling
		}

		g.camera.SetPosition(types.Vector2{
			X: g.camera.Viewport.X + cameraMove.X,
			Y: g.camera.Viewport.Y + cameraMove.Y,
		})
	}
}
