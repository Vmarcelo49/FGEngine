package editor

import (
	"fgengine/collision"
	"fgengine/types"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) handleBoxMouseEdit() {
	if g.editorManager.boxEditor == nil || g.editorManager.activeAnimation == nil {
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// Convert mouse screen coordinates to world coordinates accounting for scaling
	worldMousePos := g.getWorldMousePosition(float64(mouseX), float64(mouseY))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.editorManager.boxEditor.dragged {
			selectedBoxIndex, selectedBoxType := g.getBoxIndexUnderMouse(worldMousePos.X, worldMousePos.Y)
			if selectedBoxIndex >= 0 {
				g.editorManager.boxEditor.activeBoxIndex = selectedBoxIndex
				g.editorManager.boxEditor.activeBoxType = selectedBoxType
				g.editorManager.boxActionIndex = int(selectedBoxType) // Update UI dropdown
				g.editorManager.boxEditor.dragged = true
				g.editorManager.boxEditor.dragStartMousePos.X = worldMousePos.X
				g.editorManager.boxEditor.dragStartMousePos.Y = worldMousePos.Y

				if activeBox := g.getActiveBox(); activeBox != nil {
					g.editorManager.boxEditor.dragStartBoxPos.X = activeBox.X
					g.editorManager.boxEditor.dragStartBoxPos.Y = activeBox.Y
				}
			}
		} else {
			delta := types.Vector2{
				X: worldMousePos.X - g.editorManager.boxEditor.dragStartMousePos.X,
				Y: worldMousePos.Y - g.editorManager.boxEditor.dragStartMousePos.Y,
			}

			if activeBox := g.getActiveBox(); activeBox != nil {
				activeBox.X = g.editorManager.boxEditor.dragStartBoxPos.X + delta.X
				activeBox.Y = g.editorManager.boxEditor.dragStartBoxPos.Y + delta.Y

				g.syncCharacterActiveSprite()
			}
		}
	} else {
		// End dragging when left mouse button is released
		if g.editorManager.boxEditor.dragged {
			g.editorManager.boxEditor.dragged = false
			g.syncCharacterActiveSprite()
		}
	}
}

func (g *Game) getBoxIndexUnderMouse(worldX, worldY float64) (int, collision.BoxType) {
	if g.activeCharacter == nil || g.editorManager.boxEditor == nil {
		return -1, collision.Collision
	}

	// Get character's world position to offset the boxes
	characterPos := g.activeCharacter.GetPosition()
	point := types.Vector2{X: worldX, Y: worldY}

	// box priority order: Hit > Hurt > Collision
	boxTypes := []collision.BoxType{collision.Hit, collision.Hurt, collision.Collision}

	for _, boxType := range boxTypes {
		if boxes, exists := g.editorManager.boxEditor.boxes[boxType]; exists {
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
