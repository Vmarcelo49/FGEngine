package editor

import (
	"fgengine/collision"
	"fgengine/types"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) handleBoxMouseEdit() {
	if !*g.uiVariables.enableMouseInput {
		return
	}
	if g.character == nil {
		return
	}
	if len(g.character.Animations) == 0 {
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// Convert mouse screen coordinates to world coordinates accounting for scaling
	worldMousePos := g.getWorldMousePosition(float64(mouseX), float64(mouseY))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.uiVariables.dragged {
			selectedBoxIndex, selectedBoxType := g.getBoxIndexUnderMouse(worldMousePos.X, worldMousePos.Y)
			if selectedBoxIndex >= 0 {
				g.uiVariables.activeBoxIndex = selectedBoxIndex
				g.uiVariables.activeBoxType = selectedBoxType
				g.uiVariables.boxDropdownTypeIndex = int(selectedBoxType) // Update UI dropdown
				g.uiVariables.dragged = true
				g.uiVariables.dragStartMousePos.X = worldMousePos.X
				g.uiVariables.dragStartMousePos.Y = worldMousePos.Y

				if activeBox := g.getActiveBox(); activeBox != nil {
					g.uiVariables.dragStartBoxPos.X = activeBox.X
					g.uiVariables.dragStartBoxPos.Y = activeBox.Y
				}
			}
		} else {
			delta := types.Vector2{
				X: worldMousePos.X - g.uiVariables.dragStartMousePos.X,
				Y: worldMousePos.Y - g.uiVariables.dragStartMousePos.Y,
			}

			if activeBox := g.getActiveBox(); activeBox != nil {
				activeBox.X = g.uiVariables.dragStartBoxPos.X + delta.X
				activeBox.Y = g.uiVariables.dragStartBoxPos.Y + delta.Y
			}
		}
	} else {
		// End dragging when left mouse button is released
		if g.uiVariables.dragged {
			g.uiVariables.dragged = false
		}
	}
}

func (g *Game) getBoxIndexUnderMouse(worldX, worldY float64) (int, collision.BoxType) {
	if g.character == nil {
		return -1, collision.Collision
	}

	// Get the active frame data
	frameData := g.character.AnimationPlayer.GetActiveFrameData()
	if frameData == nil {
		return -1, collision.Collision
	}

	// Get character's world position to offset the boxes
	characterPos := g.character.GetPosition()
	point := types.Vector2{X: worldX, Y: worldY}

	// box priority order: Hit > Hurt > Collision
	boxTypes := []collision.BoxType{collision.Hit, collision.Hurt, collision.Collision}

	for _, boxType := range boxTypes {
		if boxes, exists := frameData.Boxes[boxType]; exists {
			for i, box := range boxes {
				// Transform box to world coordinates by adding character position
				worldBoxX := characterPos.X + box.X
				worldBoxY := characterPos.Y + box.Y

				if worldBoxX <= point.X && point.X <= worldBoxX+box.W &&
					worldBoxY <= point.Y && point.Y <= worldBoxY+box.H {
					// Return the index and box type when a box is found
					return i, boxType
				}
			}
		}
	}

	return -1, collision.Collision
}

// getWorldMousePosition converts screen mouse coordinates to world coordinates accounting for camera scaling
func (g *Game) getWorldMousePosition(screenX, screenY float64) types.Vector2 {
	if g.camera.Scaling != 0 && g.camera.Scaling != 1 {
		centerX := g.camera.Viewport.W / 2
		centerY := g.camera.Viewport.H / 2

		relativeX := screenX - centerX
		relativeY := screenY - centerY

		unscaledRelativeX := relativeX / g.camera.Scaling
		unscaledRelativeY := relativeY / g.camera.Scaling

		adjustedScreenX := unscaledRelativeX + centerX
		adjustedScreenY := unscaledRelativeY + centerY

		return types.Vector2{
			X: adjustedScreenX + g.camera.Viewport.X,
			Y: adjustedScreenY + g.camera.Viewport.Y,
		}
	} else {
		return types.Vector2{
			X: screenX + g.camera.Viewport.X,
			Y: screenY + g.camera.Viewport.Y,
		}
	}
}

func (g *Game) drawMouseCrosshair(screen *ebiten.Image) {
	crosshairSize := float32(10.0)
	crosshairThickness := float32(1)

	mouseX, mouseY := ebiten.CursorPosition()

	// Draw crosshair directly at mouse screen position (no coordinate conversion needed)
	centerX := float32(mouseX)
	centerY := float32(mouseY)

	crosshairColor := color.RGBA{R: 0, G: 255, B: 255, A: 255} // Cyan

	// horizontal line
	vector.FillRect(screen,
		centerX-crosshairSize, centerY-crosshairThickness/2,
		crosshairSize*2, crosshairThickness,
		crosshairColor, false)

	// vertical line
	vector.FillRect(screen,
		centerX-crosshairThickness/2, centerY-crosshairSize,
		crosshairThickness, crosshairSize*2,
		crosshairColor, false)
}
