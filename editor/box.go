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

	activeBoxType collision.BoxType
	activeBox     *types.Rect
	// mouse input related
	dragged           bool
	dragStartMousePos types.Vector2
	dragStartBoxPos   types.Vector2
}

func (g *Game) getActiveBoxindex() int {
	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil || g.editorManager.boxEditor == nil || g.editorManager.boxEditor.activeBox == nil {
		return -1
	}

	boxType := g.editorManager.boxEditor.activeBoxType
	boxes := sprite.Boxes[boxType]
	for i, box := range boxes {
		if box.X == g.editorManager.boxEditor.activeBox.X &&
			box.Y == g.editorManager.boxEditor.activeBox.Y &&
			box.W == g.editorManager.boxEditor.activeBox.W &&
			box.H == g.editorManager.boxEditor.activeBox.H {
			return i
		}
	}
	return -1
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
		ctx.Checkbox(&g.editorManager.choiceShowAllBoxes, "Show All Boxes")

		// Show current box type selection
		currentBoxType := collision.BoxType(g.editorManager.boxActionIndex)
		ctx.Text(fmt.Sprintf("Current Type: %s", currentBoxType.String()))

		if g.editorManager.boxEditor.activeBox != nil {
			boxTypeStr := g.editorManager.boxEditor.activeBoxType.String()
			ctx.Text(fmt.Sprintf("Selected: %s Box", boxTypeStr))

			ctx.SetGridLayout([]int{-1, -1}, nil)
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

		if g.editorManager.boxEditor != nil && g.editorManager.boxEditor.activeBox != nil {
			ctx.SetGridLayout([]int{-1}, nil)
			ctx.Text("Box Properties:")
			ctx.SetGridLayout([]int{-1, -1}, nil)
			ctx.Text("X:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.X, 1.0, 1).On(func() {
				g.syncCurrentSpriteToCharacter()
			})
			ctx.Text("Y:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.Y, 1.0, 1).On(func() {
				g.syncCurrentSpriteToCharacter()
			})
			ctx.Text("Width:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.W, 1.0, 1).On(func() {
				g.syncCurrentSpriteToCharacter()
			})
			ctx.Text("Height:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.H, 1.0, 1).On(func() {
				g.syncCurrentSpriteToCharacter()
			})
		}
	})
}

func (g *Game) loadBoxEditor(sprite *animation.Sprite) {
	// Ensure sprite.Boxes is initialized
	if sprite.Boxes == nil {
		sprite.Boxes = make(map[collision.BoxType][]types.Rect)
	}

	g.editorManager.boxEditor = &BoxEditor{
		boxes:         sprite.Boxes,                                      // Work directly with sprite's boxes
		activeBoxType: collision.BoxType(g.editorManager.boxActionIndex), // Initialize with current dropdown selection
	}

	// Ensure character's ActiveSprite is synchronized when loading box editor
	g.syncCharacterActiveSprite()
}

func (g *Game) clearBoxSelection() {
	if g.editorManager.boxEditor != nil {
		g.editorManager.boxEditor.activeBox = nil
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

	// Ensure sprite.Boxes is initialized
	if sprite.Boxes == nil {
		sprite.Boxes = make(map[collision.BoxType][]types.Rect)
	}

	// Point box editor directly to the current sprite's boxes and clear selection
	g.editorManager.boxEditor.boxes = sprite.Boxes
	g.clearBoxSelection()

	// Ensure character's ActiveSprite is synchronized
	g.syncCharacterActiveSprite()
}

func (g *Game) deleteSelectedBox() {
	if g.editorManager.boxEditor.activeBox == nil {
		return
	}

	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil {
		return
	}

	activeBox := g.editorManager.boxEditor.activeBox
	boxType := g.editorManager.boxEditor.activeBoxType

	// Since g.editorManager.boxEditor.boxes points to sprite.Boxes,
	// we only need to remove from one place
	for i, box := range sprite.Boxes[boxType] {
		if box.X == activeBox.X && box.Y == activeBox.Y &&
			box.W == activeBox.W && box.H == activeBox.H {
			// Remove from sprite boxes (which is the same as editor boxes)
			sprite.Boxes[boxType] = append(
				sprite.Boxes[boxType][:i],
				sprite.Boxes[boxType][i+1:]...)
			break
		}
	}

	g.clearBoxSelection()
	// Ensure character's ActiveSprite is synchronized after box deletion
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
	if g.editorManager.boxEditor.activeBox != nil {
		g.syncCharacterActiveSprite()
	}

	// Set the box type from the dropdown selection
	g.editorManager.boxEditor.activeBoxType = collision.BoxType(g.editorManager.boxActionIndex)

	selectedBox := g.addBoxOfType(sprite, g.editorManager.boxEditor.activeBoxType)
	g.editorManager.boxEditor.activeBox = selectedBox

	// Ensure character's ActiveSprite is synchronized after box addition
	g.syncCharacterActiveSprite()

	g.writeLog(fmt.Sprintf("Added %s box", g.editorManager.boxEditor.activeBoxType.String()))
}

// syncCharacterActiveSprite ensures the character's ActiveSprite points to the currently edited sprite
// This is necessary for the box renderer to display the updated boxes correctly
func (g *Game) syncCharacterActiveSprite() {
	if g.activeCharacter != nil {
		currentSprite := g.editorManager.getCurrentSprite()
		if currentSprite != nil {
			g.activeCharacter.ActiveSprite = currentSprite
		}
	}
}

// syncCurrentSpriteToCharacter is a helper method that combines the common pattern:
// getCurrentSprite() -> nil check -> syncCharacterActiveSprite()
func (g *Game) syncCurrentSpriteToCharacter() {
	sprite := g.editorManager.getCurrentSprite()
	if sprite != nil {
		g.syncCharacterActiveSprite()
	}
}

func (g *Game) addBoxOfType(currentFrame *animation.Sprite, typeOfBox collision.BoxType) *types.Rect {
	newRect := types.Rect{X: 0, Y: 0, W: 50, H: 50}

	if currentFrame.Boxes[typeOfBox] == nil {
		currentFrame.Boxes[typeOfBox] = []types.Rect{}
	}

	// Since g.editorManager.boxEditor.boxes points to currentFrame.Boxes,
	// we only need to append to the sprite boxes
	currentFrame.Boxes[typeOfBox] = append(currentFrame.Boxes[typeOfBox], newRect)

	// Return pointer to the box in the sprite's slice (the one being edited)
	spriteSlice := currentFrame.Boxes[typeOfBox]
	return &spriteSlice[len(spriteSlice)-1]
}
