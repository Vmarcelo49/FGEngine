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
	worldMousePos := g.camera.ScreenToWorld(types.Vector2{X: float64(mouseX), Y: float64(mouseY)})

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

func (g *Game) drawMouseCrosshair(screen *ebiten.Image) {
	crosshairSize := float32(10.0)
	crosshairThickness := float32(1)

	mouseX, mouseY := ebiten.CursorPosition()
	worldMousePos := g.camera.ScreenToWorld(types.Vector2{X: float64(mouseX), Y: float64(mouseY)})
	worldMouseX, worldMouseY := worldMousePos.X, worldMousePos.Y

	screenPos := g.camera.WorldToScreen(types.Vector2{X: worldMouseX, Y: worldMouseY})
	centerX := float32(screenPos.X)
	centerY := float32(screenPos.Y)

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
