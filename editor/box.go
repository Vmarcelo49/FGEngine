package editor

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/types"
	"fmt"

	"github.com/ebitengine/debugui"
)

type BoxEditor struct {
	boxes map[collision.BoxType][]types.Rect

	activeBoxType  collision.BoxType
	activeBoxIndex int
	// mouse input related
	dragged           bool
	dragStartMousePos types.Vector2
	dragStartBoxPos   types.Vector2
}

func (g *Game) getActiveBoxindex() int {
	if g.editorManager.boxEditor == nil {
		return -1
	}
	return g.editorManager.boxEditor.activeBoxIndex
}

func (g *Game) getActiveBox() *types.Rect {
	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil || g.editorManager.boxEditor == nil {
		return nil
	}

	boxType := g.editorManager.boxEditor.activeBoxType
	boxIndex := g.editorManager.boxEditor.activeBoxIndex

	if boxIndex < 0 || boxIndex >= len(sprite.Boxes[boxType]) {
		return nil
	}

	return &sprite.Boxes[boxType][boxIndex]
}

func (g *Game) boxEditor(ctx *debugui.Context) {
	ctx.Header("Box Editor", true, func() {
		sprite := g.editorManager.getCurrentSprite()
		if sprite == nil {
			ctx.Text("No sprite selected")
			return
		}

		if g.editorManager.boxEditor == nil {
			g.loadBoxEditor(sprite)
		}

		ctx.SetGridLayout([]int{-1}, nil)
		ctx.Text(fmt.Sprintf("Boxes - Collision: %d, Hurt: %d, Hit: %d", len(sprite.Boxes[collision.Collision]), len(sprite.Boxes[collision.Hurt]), len(sprite.Boxes[collision.Hit])))

		// Show current box type selection
		currentBoxType := collision.BoxType(g.editorManager.boxActionIndex)
		ctx.Text(fmt.Sprintf("Current Type: %s", currentBoxType.String()))

		if g.editorManager.boxEditor.activeBoxIndex >= 0 {
			boxTypeStr := g.editorManager.boxEditor.activeBoxType.String()
			ctx.Text(fmt.Sprintf("Selected: %s Box", boxTypeStr))

			ctx.SetGridLayout([]int{-1, -1}, nil)
			// index selection
			ctx.Text("Index:")
			activeIndex := g.getActiveBoxindex()
			ctx.NumberField(&activeIndex, 1).On(func() {
				boxes := sprite.Boxes[g.editorManager.boxEditor.activeBoxType]
				if len(boxes) > 0 && activeIndex >= 0 && activeIndex < len(boxes) {
					g.editorManager.boxEditor.activeBoxIndex = activeIndex
				}
			})

			ctx.Button("Clear Selection").On(func() {
				g.clearBoxSelection()
			})
			ctx.Button("Delete Box").On(func() {
				g.deleteSelectedBox()
			})
		} else {
			ctx.Text("No box selected")
			ctx.SetGridLayout([]int{-1}, nil)
			ctx.Button("Clear Selection").On(func() {
				g.clearBoxSelection()
			})
		}

		ctx.SetGridLayout([]int{-1}, nil)
		ctx.Text("Create New Box:")
		ctx.SetGridLayout([]int{-2, -1}, nil)
		ctx.Dropdown(&g.editorManager.boxActionIndex, []string{collision.Collision.String(), collision.Hit.String(), collision.Hurt.String()}).On(func() {
			g.clearBoxSelection()
			g.editorManager.boxEditor.activeBoxType = collision.BoxType(g.editorManager.boxActionIndex)
		})
		ctx.Button("Add Box").On(func() {
			g.addBox()
		})

		if g.editorManager.boxEditor != nil {
			activeBox := g.getActiveBox()
			if activeBox != nil {
				ctx.SetGridLayout([]int{-1}, nil)
				ctx.Text("Box Properties:")
				ctx.SetGridLayout([]int{-1, -1}, nil)
				ctx.Text("X:")
				ctx.NumberFieldF(&activeBox.X, 1.0, 1).On(func() {
					g.syncCharacterActiveSprite()
				})
				ctx.Text("Y:")
				ctx.NumberFieldF(&activeBox.Y, 1.0, 1).On(func() {
					g.syncCharacterActiveSprite()
				})
				ctx.Text("Width:")
				ctx.NumberFieldF(&activeBox.W, 1.0, 1).On(func() {
					g.syncCharacterActiveSprite()
				})
				ctx.Text("Height:")
				ctx.NumberFieldF(&activeBox.H, 1.0, 1).On(func() {
					g.syncCharacterActiveSprite()
				})
			}
		}
	})
}

func (g *Game) loadBoxEditor(sprite *animation.Sprite) {
	if sprite.Boxes == nil {
		sprite.Boxes = make(map[collision.BoxType][]types.Rect)
	}

	g.editorManager.boxEditor = &BoxEditor{
		boxes:          sprite.Boxes,
		activeBoxType:  collision.BoxType(g.editorManager.boxActionIndex),
		activeBoxIndex: -1,
	}

	g.syncCharacterActiveSprite()
}

func (g *Game) clearBoxSelection() {
	if g.editorManager.boxEditor != nil {
		g.editorManager.boxEditor.activeBoxIndex = -1
	}
}

// refreshBoxEditor updates the box editor when the frame changes
func (g *Game) refreshBoxEditor() {
	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil {
		return
	}

	// If no box editor exists, create it
	if g.editorManager.boxEditor == nil {
		g.loadBoxEditor(sprite)
		return
	}

	if sprite.Boxes == nil {
		sprite.Boxes = make(map[collision.BoxType][]types.Rect)
	}

	// Point box editor directly to the current sprite's boxes and clear selection
	g.editorManager.boxEditor.boxes = sprite.Boxes
	g.clearBoxSelection()

	g.syncCharacterActiveSprite()
}

func (g *Game) deleteSelectedBox() {
	activeBox := g.getActiveBox()
	if activeBox == nil {
		return
	}

	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil {
		return
	}

	boxType := g.editorManager.boxEditor.activeBoxType
	boxIndex := g.editorManager.boxEditor.activeBoxIndex

	// Remove the box from the slice
	sprite.Boxes[boxType] = append(
		sprite.Boxes[boxType][:boxIndex],
		sprite.Boxes[boxType][boxIndex+1:]...)

	g.clearBoxSelection()
	g.syncCharacterActiveSprite()
}

func (g *Game) addBox() {
	if g.editorManager.activeAnimation == nil {
		g.writeLog("Cannot add box: No active animation")
		return
	}

	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil {
		g.writeLog("Cannot add box: No active sprite")
		return
	}

	if g.editorManager.boxEditor == nil {
		g.loadBoxEditor(sprite)
	}

	// Save any pending edits before adding a new box
	activeBox := g.getActiveBox()
	if activeBox != nil {
		g.syncCharacterActiveSprite()
	}

	// Set the box type from the dropdown selection
	g.editorManager.boxEditor.activeBoxType = collision.BoxType(g.editorManager.boxActionIndex)

	selectedBoxIndex := g.addBoxOfType(sprite, g.editorManager.boxEditor.activeBoxType)
	g.editorManager.boxEditor.activeBoxIndex = selectedBoxIndex

	g.syncCharacterActiveSprite()

	g.writeLog(fmt.Sprintf("Added %s box", g.editorManager.boxEditor.activeBoxType.String()))
}

func (g *Game) syncCharacterActiveSprite() {
	sprite := g.editorManager.getCurrentSprite()
	if sprite != nil {
		if g.activeCharacter != nil {
			currentSprite := g.editorManager.getCurrentSprite()
			if currentSprite != nil {
				g.activeCharacter.ActiveSprite = currentSprite
			}
		}
	}
}

func (g *Game) addBoxOfType(currentFrame *animation.Sprite, typeOfBox collision.BoxType) int {
	newRect := types.Rect{X: 0, Y: 0, W: 50, H: 50}

	if currentFrame.Boxes[typeOfBox] == nil {
		currentFrame.Boxes[typeOfBox] = []types.Rect{}
	}
	currentFrame.Boxes[typeOfBox] = append(currentFrame.Boxes[typeOfBox], newRect)

	// Return the index of the newly added box
	return len(currentFrame.Boxes[typeOfBox]) - 1
}
