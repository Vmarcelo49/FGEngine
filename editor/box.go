package editor

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/types"
	"fmt"

	"github.com/ebitengine/debugui"
)

type BoxEditor struct {
	collisionBoxes []collision.Box
	hurtBoxes      []collision.Box
	hitBoxes       []collision.Box
	activeBoxType  collision.BoxType
	activeBoxIndex int

	activeBox *collision.Box

	// mouse input related
	dragged           bool
	dragStartMousePos struct{ X, Y float64 }
	dragStartBoxPos   struct{ X, Y float64 }
}

func (g *Game) GetActiveBoxType() collision.BoxType {
	if g.editorManager.boxEditor != nil && g.editorManager.boxEditor.activeBox != nil {
		return g.editorManager.boxEditor.activeBox.BoxType
	}
	return collision.Collision // theres no default or unknown box type, 0 is Collision
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
		ctx.Text(fmt.Sprintf("Boxes - Collision: %d, Hurt: %d, Hit: %d", len(sprite.CollisionBoxes), len(sprite.HurtBoxes), len(sprite.HitBoxes)))
		ctx.Checkbox(&g.editorManager.choiceShowAllBoxes, "Show All Boxes")
		if g.editorManager.boxEditor.activeBox != nil {
			boxTypeStr := g.GetActiveBoxType().String()
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
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.Rect.X, 1.0, 1)
			ctx.Text("Y:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.Rect.Y, 1.0, 1)
			ctx.Text("Width:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.Rect.W, 1.0, 1)
			ctx.Text("Height:")
			ctx.NumberFieldF(&g.editorManager.boxEditor.activeBox.Rect.H, 1.0, 1)
		}
	})
}

// loads boxes into the BoxEditor from a sprite
func (g *Game) loadBoxEditor(sprite *animation.Sprite) {
	if sprite == nil {
		return
	}
	g.editorManager.boxEditor = &BoxEditor{}

	for _, box := range sprite.CollisionBoxes {
		g.editorManager.boxEditor.collisionBoxes = append(g.editorManager.boxEditor.collisionBoxes, collision.Box{Rect: box, BoxType: collision.Collision})
	}
	for _, box := range sprite.HitBoxes {
		g.editorManager.boxEditor.hitBoxes = append(g.editorManager.boxEditor.hitBoxes, collision.Box{Rect: box, BoxType: collision.Hit})
	}
	for _, box := range sprite.HurtBoxes {
		g.editorManager.boxEditor.hurtBoxes = append(g.editorManager.boxEditor.hurtBoxes, collision.Box{Rect: box, BoxType: collision.Hurt})
	}

}

func (g *Game) clearBoxSelection() {
	g.editorManager.boxEditor.activeBox = nil
}

func (g *Game) deleteSelectedBox() {
	if g.editorManager.boxEditor == nil || g.editorManager.boxEditor.activeBox == nil {
		return
	}

	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil {
		return
	}

	activeBox := g.editorManager.boxEditor.activeBox
	activeIndex := g.editorManager.boxEditor.activeBoxIndex

	switch activeBox.BoxType {
	case collision.Collision:
		// Delete from editor's collision boxes
		for i, box := range g.editorManager.boxEditor.collisionBoxes {
			if box == *activeBox {
				g.editorManager.boxEditor.collisionBoxes = append(g.editorManager.boxEditor.collisionBoxes[:i], g.editorManager.boxEditor.collisionBoxes[i+1:]...)
				break
			}
		}
		// Delete from sprite's collision boxes
		if activeIndex >= 0 && activeIndex < len(sprite.CollisionBoxes) {
			sprite.CollisionBoxes = append(sprite.CollisionBoxes[:activeIndex], sprite.CollisionBoxes[activeIndex+1:]...)
		}
	case collision.Hit:
		// Delete from editor's hit boxes
		for i, box := range g.editorManager.boxEditor.hitBoxes {
			if box == *activeBox {
				g.editorManager.boxEditor.hitBoxes = append(g.editorManager.boxEditor.hitBoxes[:i], g.editorManager.boxEditor.hitBoxes[i+1:]...)
				break
			}
		}
		// Delete from sprite's hit boxes
		if activeIndex >= 0 && activeIndex < len(sprite.HitBoxes) {
			sprite.HitBoxes = append(sprite.HitBoxes[:activeIndex], sprite.HitBoxes[activeIndex+1:]...)
		}
	case collision.Hurt:
		// Delete from editor's hurt boxes
		for i, box := range g.editorManager.boxEditor.hurtBoxes {
			if box == *activeBox {
				g.editorManager.boxEditor.hurtBoxes = append(g.editorManager.boxEditor.hurtBoxes[:i], g.editorManager.boxEditor.hurtBoxes[i+1:]...)
				break
			}
		}
		// Delete from sprite's hurt boxes
		if activeIndex >= 0 && activeIndex < len(sprite.HurtBoxes) {
			sprite.HurtBoxes = append(sprite.HurtBoxes[:activeIndex], sprite.HurtBoxes[activeIndex+1:]...)
		}
	default:
		return
	}

	// Invalidate the box cache since we deleted a box
	sprite.InvalidateBoxCache()
	g.clearBoxSelection()
}

func (g *Game) addBox() {
	if g.editorManager.boxEditor == nil || g.editorManager.activeAnimation == nil {
		return
	}

	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil {
		return
	}

	// any edits are saved to the frame before adding a new box
	g.updateFrameBoxes(sprite)

	typeOfBox := g.editorManager.boxEditor.activeBoxType

	// Use the helper function to add the box and get a pointer to it
	selectedBox, selectedIndex := g.addBoxOfType(sprite, typeOfBox)

	// Update selection state
	g.editorManager.boxEditor.activeBox = selectedBox
	g.editorManager.boxEditor.activeBoxIndex = selectedIndex
	g.editorManager.boxEditor.activeBoxType = typeOfBox
}

// updateFrameBoxes updates the current frame's box arrays with the modified box data
func (g *Game) updateFrameBoxes(sprite *animation.Sprite) {
	if sprite == nil || g.editorManager.boxEditor.activeBox == nil {
		return
	}

	g.updateBoxSlice(g.editorManager.boxEditor.activeBox.BoxType, sprite, g.editorManager.boxEditor.activeBoxIndex, g.editorManager.boxEditor.activeBox.Rect)
}

// updateBoxSlice updates a specific box in the appropriate slice based on box type
func (g *Game) updateBoxSlice(boxType collision.BoxType, sprite *animation.Sprite, index int, rect types.Rect) {
	switch boxType {
	case collision.Collision:
		if index >= 0 && index < len(sprite.CollisionBoxes) {
			sprite.CollisionBoxes[index] = rect
		}
	case collision.Hit:
		if index >= 0 && index < len(sprite.HitBoxes) {
			sprite.HitBoxes[index] = rect
		}
	case collision.Hurt:
		if index >= 0 && index < len(sprite.HurtBoxes) {
			sprite.HurtBoxes[index] = rect
		}
	}
	// Invalidate the box cache since we modified the boxes
	sprite.InvalidateBoxCache()
}

func (g *Game) addBoxOfType(currentFrame *animation.Sprite, typeOfBox collision.BoxType) (*collision.Box, int) {
	newRect := types.Rect{X: 0, Y: 0, W: 50, H: 50}

	var box *collision.Box
	var index int

	switch typeOfBox {
	case collision.Collision:
		box, index = g.appendBoxToSlices(&currentFrame.CollisionBoxes, &g.editorManager.boxEditor.collisionBoxes, newRect, typeOfBox)
	case collision.Hit:
		box, index = g.appendBoxToSlices(&currentFrame.HitBoxes, &g.editorManager.boxEditor.hitBoxes, newRect, typeOfBox)
	case collision.Hurt:
		box, index = g.appendBoxToSlices(&currentFrame.HurtBoxes, &g.editorManager.boxEditor.hurtBoxes, newRect, typeOfBox)
	default:
		return nil, -1
	}

	// Invalidate the box cache since we added a new box
	currentFrame.InvalidateBoxCache()
	return box, index
}

func (g *Game) appendBoxToSlices(frameSlice *[]types.Rect, editorSlice *[]collision.Box, newRect types.Rect, boxType collision.BoxType) (*collision.Box, int) {
	*frameSlice = append(*frameSlice, newRect)
	newBox := collision.Box{Rect: newRect, BoxType: boxType}
	*editorSlice = append(*editorSlice, newBox)
	newBoxIndex := len(*editorSlice) - 1
	return &(*editorSlice)[newBoxIndex], newBoxIndex
}
