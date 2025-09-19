package editor

import (
	"fgengine/config"
	"fgengine/input"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CameraMoveSpeed = 5.0
)

func (g *Game) handleCameraInput() {
	// Handle mouse drag for camera movement
	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if !g.isDragging {
			// Start dragging
			g.isDragging = true
			g.lastMouseX = x
			g.lastMouseY = y
		} else {
			// Continue dragging
			deltaX := float64(g.lastMouseX - x)
			deltaY := float64(g.lastMouseY - y)

			// Apply mouse movement to camera (inverted because we want to "grab" the world)
			g.camera.SetPosition(types.Vector2{
				X: g.camera.Viewport.X + deltaX,
				Y: g.camera.Viewport.Y + deltaY,
			})

			g.lastMouseX = x
			g.lastMouseY = y
		}
	} else {
		g.isDragging = false
	}

	// Handle keyboard input for camera movement
	inputs := g.inputManager.GetLocalInputs()

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

	// Apply camera movement
	if cameraMove.X != 0 || cameraMove.Y != 0 {
		g.camera.SetPosition(types.Vector2{
			X: g.camera.Viewport.X + cameraMove.X,
			Y: g.camera.Viewport.Y + cameraMove.Y,
		})
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.camera.SetPosition(types.Vector2{X: -float64(config.WindowWidth) / 2, Y: -float64(config.WindowHeight) / 2})
	}
}
