package editor

import (
	"fgengine/collision"
	"fgengine/types"
	"fmt"

	"github.com/ebitengine/debugui"
)

func (g *Game) getActiveBox() *types.Rect {
	frameData := g.character.AnimationPlayer.GetActiveFrameData()
	if frameData == nil {
		return nil
	}

	boxType := g.uiVariables.activeBoxType
	boxIndex := g.uiVariables.activeBoxIndex

	return &frameData.Boxes[boxType][boxIndex]
}

func (g *Game) boxEditor(ctx *debugui.Context) {
	ctx.Header("Box Editor", true, func() {
		frameData := g.character.AnimationPlayer.GetActiveFrameData()
		if frameData == nil {
			ctx.Text("No frame data available")
			return
		}

		ctx.SetGridLayout([]int{-1}, nil)
		ctx.Text(fmt.Sprintf("Boxes - Collision: %d, Hurt: %d, Hit: %d", len(frameData.Boxes[collision.Collision]), len(frameData.Boxes[collision.Hurt]), len(frameData.Boxes[collision.Hit])))

		// Show current box type selection
		currentBoxType := collision.BoxType(g.uiVariables.boxDropdownIndex)
		ctx.Text(fmt.Sprintf("Current Type: %s", currentBoxType.String()))

		if g.uiVariables.activeBoxIndex >= 0 {
			ctx.Text(fmt.Sprintf("Selected: %s Box", g.uiVariables.activeBoxType.String()))

			ctx.SetGridLayout([]int{-1, -1}, nil)
			// index selection
			ctx.Text("Index:")
			ctx.NumberField(&g.uiVariables.activeBoxIndex, 1).On(func() {
				boxes := frameData.Boxes[g.uiVariables.activeBoxType]
				if len(boxes) > 0 && g.uiVariables.activeBoxIndex >= 0 && g.uiVariables.activeBoxIndex < len(boxes) {
					g.uiVariables.activeBoxIndex = g.uiVariables.activeBoxIndex
				}
			})

			ctx.Button("Clear Selection").On(func() {
				g.uiVariables.activeBoxIndex = -1
			})
			ctx.Button("Delete Box").On(func() {
				g.deleteSelectedBox()
			})
		} else {
			ctx.Text("No box selected")
			ctx.SetGridLayout([]int{-1}, nil)
			ctx.Button("Clear Selection").On(func() {
				g.uiVariables.activeBoxIndex = -1
			})
		}

		ctx.SetGridLayout([]int{-1}, nil)
		ctx.Text("Create New Box:")
		ctx.SetGridLayout([]int{-2, -1}, nil)
		ctx.Dropdown(&g.uiVariables.boxDropdownIndex, []string{collision.Collision.String(), collision.Hit.String(), collision.Hurt.String()}).On(func() {
			g.clearBoxSelection()
			g.uiVariables.boxEditor.activeBoxType = collision.BoxType(g.uiVariables.boxDropdownIndex)
		})
		ctx.Button("Add Box").On(func() {
			g.addBox()
		})

		if g.uiVariables.boxEditor != nil {
			activeBox := g.getActiveBox()
			if activeBox != nil {
				ctx.SetGridLayout([]int{-1}, nil)
				ctx.Text("Box Properties:")
				ctx.SetGridLayout([]int{-1, -1}, nil)
				ctx.Text("X:")
				ctx.NumberFieldF(&activeBox.X, 1.0, 1)
				ctx.Text("Y:")
				ctx.NumberFieldF(&activeBox.Y, 1.0, 1)
				ctx.Text("Width:")
				ctx.NumberFieldF(&activeBox.W, 1.0, 1)
				ctx.Text("Height:")
				ctx.NumberFieldF(&activeBox.H, 1.0, 1)
			}
		}
	})
}

func (g *Game) deleteSelectedBox() {
	activeBox := g.getActiveBox()
	if activeBox == nil {
		g.writeLog("tried to remove a nil box wth")
		return
	}

	frameData := g.character.AnimationPlayer.GetActiveFrameData()
	if frameData == nil {
		return
	}

	boxType := g.uiVariables.activeBoxType
	boxIndex := g.uiVariables.activeBoxIndex

	frameData.Boxes[boxType] = append(frameData.Boxes[boxType][:boxIndex], frameData.Boxes[boxType][boxIndex+1:]...)
	g.uiVariables.activeBoxIndex = -1 // Clear selection after deletion
	g.writeLog(fmt.Sprintf("Deleted %s box at index %d", boxType.String(), boxIndex))
}

func (g *Game) addBox() {
	frameData := g.character.AnimationPlayer.GetActiveFrameData()
	if frameData == nil {
		g.writeLog("No active frame data to add box to")
		return
	}
	boxType := collision.BoxType(g.uiVariables.boxDropdownIndex)
	newRect := types.Rect{X: 0, Y: 0, W: 50, H: 50}

	frameData.Boxes[boxType] = append(frameData.Boxes[boxType], newRect)
	g.uiVariables.activeBoxIndex = len(frameData.Boxes[boxType]) - 1
	g.writeLog(fmt.Sprintf("Added %s box", boxType.String()))

	g.uiVariables.activeBoxIndex = len(frameData.Boxes[boxType]) - 1
}
