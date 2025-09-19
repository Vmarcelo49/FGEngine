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
	if sprite == nil {
		return -1
	}

	boxes := sprite.Boxes[g.editorManager.boxEditor.activeBoxType]
	for i, box := range boxes {
		if box == *g.editorManager.boxEditor.activeBox {
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
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.X, 1.0, 1)
			ctx.Text("Y:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.Y, 1.0, 1)
			ctx.Text("Width:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.W, 1.0, 1)
			ctx.Text("Height:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.H, 1.0, 1)
		}
	})
}

func (g *Game) loadBoxEditor(sprite *animation.Sprite) {
	g.editorManager.boxEditor = &BoxEditor{
		boxes: make(map[collision.BoxType][]types.Rect),
	}
	copyBoxes(sprite.Boxes, g.editorManager.boxEditor.boxes)
}

func (g *Game) clearBoxSelection() {
	g.editorManager.boxEditor.activeBox = nil
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

	for i, box := range g.editorManager.boxEditor.boxes[boxType] {
		if box == *activeBox {
			g.editorManager.boxEditor.boxes[boxType] = append(
				g.editorManager.boxEditor.boxes[boxType][:i],
				g.editorManager.boxEditor.boxes[boxType][i+1:]...)

			if i < len(sprite.Boxes[boxType]) {
				sprite.Boxes[boxType] = append(
					sprite.Boxes[boxType][:i],
					sprite.Boxes[boxType][i+1:]...)
			}
			break
		}
	}

	g.clearBoxSelection()
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

	// any edits are saved to the frame before adding a new box
	g.updateFrameBoxes(sprite)
	selectedBox := g.addBoxOfType(sprite, g.editorManager.boxEditor.activeBoxType)

	g.editorManager.boxEditor.activeBox = selectedBox
}

func (g *Game) updateFrameBoxes(sprite *animation.Sprite) {
	g.updateBoxSlice(g.editorManager.boxEditor.activeBoxType, sprite, g.getActiveBoxindex(), *g.editorManager.boxEditor.activeBox)
}

func (g *Game) updateBoxSlice(boxType collision.BoxType, sprite *animation.Sprite, index int, rect types.Rect) {
	if g.editorManager.boxEditor.boxes == nil {
		g.editorManager.boxEditor.boxes = make(map[collision.BoxType][]types.Rect)
	}

	if index >= 0 && index < len(g.editorManager.boxEditor.boxes[boxType]) {
		g.editorManager.boxEditor.boxes[boxType][index] = rect

		if index < len(sprite.Boxes[boxType]) {
			sprite.Boxes[boxType][index] = rect
		}
	}
}

func (g *Game) addBoxOfType(currentFrame *animation.Sprite, typeOfBox collision.BoxType) *types.Rect {
	newRect := types.Rect{X: 0, Y: 0, W: 50, H: 50}

	if currentFrame.Boxes[typeOfBox] == nil {
		currentFrame.Boxes[typeOfBox] = []types.Rect{}
	}
	if g.editorManager.boxEditor.boxes[typeOfBox] == nil {
		g.editorManager.boxEditor.boxes[typeOfBox] = []types.Rect{}
	}

	currentFrame.Boxes[typeOfBox] = append(currentFrame.Boxes[typeOfBox], newRect)
	g.editorManager.boxEditor.boxes[typeOfBox] = append(g.editorManager.boxEditor.boxes[typeOfBox], newRect)

	// pointer to the box in the editor slice (the one being edited)
	editorSlice := g.editorManager.boxEditor.boxes[typeOfBox]
	return &editorSlice[len(editorSlice)-1]
}
