package editor

import (
	"fgengine/collision"
	"fgengine/config"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

// handleMouseInput processes mouse input for box selection and dragging
func (g *Game) handleMouseInput() {
	if g.editorManager.boxEditor == nil || g.editorManager.activeAnimation == nil {
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()
	worldMouseX, worldMouseY := g.screenToWorldPos(float64(mouseX), float64(mouseY))
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.editorManager.boxEditor.dragged {
			selectedBox := g.getBoxUnderMouse(worldMouseX, worldMouseY)
			if selectedBox != nil {
				g.editorManager.boxEditor.activeBox = selectedBox
				g.editorManager.boxEditor.dragged = true
				g.editorManager.boxEditor.dragStartMousePos.X = worldMouseX
				g.editorManager.boxEditor.dragStartMousePos.Y = worldMouseY
				g.editorManager.boxEditor.dragStartBoxPos.X = selectedBox.X
				g.editorManager.boxEditor.dragStartBoxPos.Y = selectedBox.Y

				//g.editorManager.boxEditor.activeBoxType = selectedBox.type
			}
		} else {
			deltaX := worldMouseX - g.editorManager.boxEditor.dragStartMousePos.X
			deltaY := worldMouseY - g.editorManager.boxEditor.dragStartMousePos.Y

			if g.editorManager.boxEditor.activeBox != nil {
				g.editorManager.boxEditor.activeBox.X = g.editorManager.boxEditor.dragStartBoxPos.X + deltaX
				g.editorManager.boxEditor.activeBox.Y = g.editorManager.boxEditor.dragStartBoxPos.Y + deltaY

				sprite := g.editorManager.getCurrentSprite()
				if sprite != nil {
					g.updateFrameBoxes(sprite)
				}
			}
		}
	} else {
		if g.editorManager.boxEditor.dragged {
			g.editorManager.boxEditor.dragged = false
		}
	}
}

// screenToWorldPos converts screen coordinates to world coordinates
func (g *Game) screenToWorldPos(screenX, screenY float64) (float64, float64) {

	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil {
		return screenX / Zoom, screenY / Zoom
	}

	spriteW := sprite.Rect.W
	spriteH := sprite.Rect.H

	spriteScreenAnchorX := float64(config.WindowWidth)/2 - spriteW/2*Zoom
	spriteScreenAnchorY := float64(config.WindowHeight)/2 - spriteH/2*Zoom

	mouseXRelativeToSpriteAnchorScreen := screenX - spriteScreenAnchorX
	mouseYRelativeToSpriteAnchorScreen := screenY - spriteScreenAnchorY

	worldX := mouseXRelativeToSpriteAnchorScreen / Zoom
	worldY := mouseYRelativeToSpriteAnchorScreen / Zoom

	return worldX, worldY
}

// getBoxUnderMouse returns the box under the mouse cursor, if any
func (g *Game) getBoxUnderMouse(worldX, worldY float64) *types.Rect {
	point := types.Vector2{X: worldX, Y: worldY}

	for i, box := range g.editorManager.boxEditor.boxes[collision.Collision] {
		if box.Contains(point.X, point.Y) {
			return &g.editorManager.boxEditor.boxes[collision.Collision][i]
		}
	}

	for i, box := range g.editorManager.boxEditor.boxes[collision.Hit] {
		if box.Contains(point.X, point.Y) {
			return &g.editorManager.boxEditor.boxes[collision.Hit][i]
		}
	}

	for i, box := range g.editorManager.boxEditor.boxes[collision.Hurt] {
		if box.Contains(point.X, point.Y) {
			return &g.editorManager.boxEditor.boxes[collision.Hurt][i]
		}
	}

	return nil
}
