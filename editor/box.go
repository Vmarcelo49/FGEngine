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

	boxTypeSelectionDropdown collision.BoxType // needed for box creation
	activeBox                *collision.Box
	// mouse input related
	dragged           bool
	dragStartMousePos types.Vector2
	dragStartBoxPos   types.Vector2
}

func (g *Game) getActiveBox() *collision.Box {
	if g.editorManager != nil && g.editorManager.boxEditor != nil {
		return g.editorManager.boxEditor.activeBox
	}
	return nil
}

func (g *Game) getActiveBoxType() collision.BoxType {
	if activeBox := g.getActiveBox(); activeBox != nil {
		return activeBox.BoxType
	}
	return collision.Collision // theres no default or unknown box type, 0 is Collision
}

func (g *Game) getActiveBoxindex() int {
	activeBox := g.getActiveBox()
	if activeBox == nil {
		return -1
	}

	switch g.getActiveBoxType() {
	case collision.Collision:
		for i, box := range g.editorManager.boxEditor.collisionBoxes {
			if box == *activeBox {
				return i
			}
		}
	case collision.Hit:
		for i, box := range g.editorManager.boxEditor.hitBoxes {
			if box == *activeBox {
				return i
			}
		}
	case collision.Hurt:
		for i, box := range g.editorManager.boxEditor.hurtBoxes {
			if box == *activeBox {
				return i
			}
		}
	}
	return -1
}

func (g *Game) boxEditor(ctx *debugui.Context) {
	if g.editorManager == nil {
		ctx.Text("Editor manager not initialized")
		return
	}

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
			boxTypeStr := g.getActiveBoxType().String()
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
			g.editorManager.boxEditor.boxTypeSelectionDropdown = collision.BoxType(g.editorManager.boxActionIndex)
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
	g.editorManager.boxEditor = &BoxEditor{
		collisionBoxes: make([]collision.Box, 0, len(sprite.CollisionBoxes)),
		hurtBoxes:      make([]collision.Box, 0, len(sprite.HurtBoxes)),
		hitBoxes:       make([]collision.Box, 0, len(sprite.HitBoxes)),
	}

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
	if g.editorManager != nil && g.editorManager.boxEditor != nil {
		g.editorManager.boxEditor.activeBox = nil
	}
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
	activeIndex := g.getActiveBoxindex()

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
	if g.editorManager == nil || g.editorManager.boxEditor == nil || g.editorManager.activeAnimation == nil {
		return
	}

	sprite := g.editorManager.getCurrentSprite()
	if sprite == nil {
		return
	}

	// any edits are saved to the frame before adding a new box
	g.updateFrameBoxes(sprite)
	selectedBox := g.addBoxOfType(sprite, g.editorManager.boxEditor.boxTypeSelectionDropdown)

	g.editorManager.boxEditor.activeBox = selectedBox
}

// updateFrameBoxes updates the current frame's box arrays with the modified box data
func (g *Game) updateFrameBoxes(sprite *animation.Sprite) {
	if g.editorManager == nil || g.editorManager.boxEditor == nil || g.editorManager.boxEditor.activeBox == nil || sprite == nil {
		return
	}

	g.updateBoxSlice(g.editorManager.boxEditor.activeBox.BoxType, sprite, g.getActiveBoxindex(), g.editorManager.boxEditor.activeBox.Rect)
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

func (g *Game) addBoxOfType(currentFrame *animation.Sprite, typeOfBox collision.BoxType) *collision.Box {
	newRect := types.Rect{X: 0, Y: 0, W: 50, H: 50}

	var box *collision.Box

	switch typeOfBox {
	case collision.Collision:
		box = g.appendBoxToSlices(&currentFrame.CollisionBoxes, &g.editorManager.boxEditor.collisionBoxes, newRect, typeOfBox)
	case collision.Hit:
		box = g.appendBoxToSlices(&currentFrame.HitBoxes, &g.editorManager.boxEditor.hitBoxes, newRect, typeOfBox)
	case collision.Hurt:
		box = g.appendBoxToSlices(&currentFrame.HurtBoxes, &g.editorManager.boxEditor.hurtBoxes, newRect, typeOfBox)
	default:
		return nil
	}

	// Invalidate the box cache since we added a new box
	currentFrame.InvalidateBoxCache()
	return box
}

func (g *Game) appendBoxToSlices(frameSlice *[]types.Rect, editorSlice *[]collision.Box, newRect types.Rect, boxType collision.BoxType) *collision.Box {
	*frameSlice = append(*frameSlice, newRect)
	newBox := collision.Box{Rect: newRect, BoxType: boxType}
	*editorSlice = append(*editorSlice, newBox)
	return &(*editorSlice)[len(*editorSlice)-1]
}
