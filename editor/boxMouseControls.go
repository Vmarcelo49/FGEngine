package editor

import (
	"fgengine/collision"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

// handleMouseInput processes mouse input for box selection and dragging
// Uses left mouse button only to avoid conflicts with camera controls (middle/right mouse)
func (g *Game) handleMouseInput() {
	if g.editorManager.boxEditor == nil || g.editorManager.activeAnimation == nil {
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()
	worldMouseX, worldMouseY := g.screenToWorldPos(float64(mouseX), float64(mouseY))

	// Only handle left mouse button for box editing to avoid conflict with camera controls
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.editorManager.boxEditor.dragged {
			// Start dragging: find box under mouse and set up drag state
			selectedBox := g.getBoxUnderMouse(worldMouseX, worldMouseY)
			if selectedBox != nil {
				g.editorManager.boxEditor.activeBox = selectedBox
				g.editorManager.boxEditor.dragged = true
				g.editorManager.boxEditor.dragStartMousePos.X = worldMouseX
				g.editorManager.boxEditor.dragStartMousePos.Y = worldMouseY
				g.editorManager.boxEditor.dragStartBoxPos.X = selectedBox.X
				g.editorManager.boxEditor.dragStartBoxPos.Y = selectedBox.Y
			}
		} else {
			// Continue dragging: update box position based on mouse delta
			deltaX := worldMouseX - g.editorManager.boxEditor.dragStartMousePos.X
			deltaY := worldMouseY - g.editorManager.boxEditor.dragStartMousePos.Y

			if g.editorManager.boxEditor.activeBox != nil {
				g.editorManager.boxEditor.activeBox.X = g.editorManager.boxEditor.dragStartBoxPos.X + deltaX
				g.editorManager.boxEditor.activeBox.Y = g.editorManager.boxEditor.dragStartBoxPos.Y + deltaY

				// Update the frame boxes in the sprite data
				sprite := g.editorManager.getCurrentSprite()
				if sprite != nil {
					g.updateFrameBoxes(sprite)
				}
			}
		}
	} else {
		// End dragging when left mouse button is released
		if g.editorManager.boxEditor.dragged {
			g.editorManager.boxEditor.dragged = false
		}
	}
}

// screenToWorldPos converts screen coordinates to world coordinates using the camera system
func (g *Game) screenToWorldPos(screenX, screenY float64) (float64, float64) {
	screenPos := types.Vector2{X: screenX, Y: screenY}
	worldPos := g.camera.ScreenToWorld(screenPos)
	return worldPos.X, worldPos.Y
}

// getBoxUnderMouse returns the box under the mouse cursor, if any
// Accounts for character world position when checking box collision
// Sets the active box type when a box is found
func (g *Game) getBoxUnderMouse(worldX, worldY float64) *types.Rect {
	if g.activeCharacter == nil || g.editorManager.boxEditor == nil {
		return nil
	}

	// Get character's world position to offset the boxes
	characterPos := g.activeCharacter.GetPosition()
	point := types.Vector2{X: worldX, Y: worldY}

	// Priority order for box selection: Hit > Hurt > Collision
	// This way smaller hit boxes are selected before larger hurt boxes
	boxTypes := []collision.BoxType{collision.Hit, collision.Hurt, collision.Collision}

	for _, boxType := range boxTypes {
		if boxes, exists := g.editorManager.boxEditor.boxes[boxType]; exists {
			for i, box := range boxes {
				// Transform box to world coordinates by adding character position
				worldBoxX := characterPos.X + box.X
				worldBoxY := characterPos.Y + box.Y

				if worldBoxX <= point.X && point.X <= worldBoxX+box.W &&
					worldBoxY <= point.Y && point.Y <= worldBoxY+box.H {
					// Set the active box type when a box is found
					g.editorManager.boxEditor.activeBoxType = boxType
					g.editorManager.boxActionIndex = int(boxType) // Update UI dropdown
					return &g.editorManager.boxEditor.boxes[boxType][i]
				}
			}
		}
	}

	return nil
}
