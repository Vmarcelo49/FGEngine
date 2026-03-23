//go:build !js && !wasm
// +build !js,!wasm

package editorimgui

import (
	"fgengine/collision"
	"fgengine/types"
	"fmt"
)

func (g *Game) activeBox() *types.Rect {
	player := g.activeAnimPlayer()
	if player == nil {
		return nil
	}

	frameData := player.ActiveFrameData()
	if frameData == nil {
		return nil
	}
	if frameData.Boxes == nil {
		return nil
	}

	boxType := g.uiVariables.activeBoxType
	boxIndex := g.uiVariables.activeBoxIndex
	boxes, exists := frameData.Boxes[boxType]
	if !exists || boxIndex < 0 || boxIndex >= len(boxes) {
		return nil
	}

	return &frameData.Boxes[boxType][boxIndex]
}

func (g *Game) deleteSelectedBox() {
	activeBox := g.activeBox()
	if activeBox == nil {
		g.writeLog("cannot delete box: no active box selected")
		return
	}

	player := g.activeAnimPlayer()
	if player == nil {
		return
	}

	frameData := player.ActiveFrameData()
	if frameData == nil {
		return
	}

	boxType := g.uiVariables.activeBoxType
	boxIndex := g.uiVariables.activeBoxIndex

	frameData.Boxes[boxType] = append(frameData.Boxes[boxType][:boxIndex], frameData.Boxes[boxType][boxIndex+1:]...)
	g.uiVariables.activeBoxIndex = -1
	g.writeLog(fmt.Sprintf("Deleted %s box at index %d", boxType.String(), boxIndex))
}

func (g *Game) addBox() {
	player := g.activeAnimPlayer()
	if player == nil {
		g.writeLog("No animation player available")
		return
	}

	frameData := player.ActiveFrameData()
	if frameData == nil {
		g.writeLog("No active frame data to add box to")
		return
	}
	if frameData.Boxes == nil {
		frameData.Boxes = make(map[collision.BoxType][]types.Rect)
	}
	boxType := collision.BoxType(g.uiVariables.boxDropdownTypeIndex)
	newRect := types.Rect{X: 0, Y: 0, W: 50, H: 50}

	frameData.Boxes[boxType] = append(frameData.Boxes[boxType], newRect)
	g.uiVariables.activeBoxIndex = len(frameData.Boxes[boxType]) - 1
	g.writeLog(fmt.Sprintf("Added %s box", boxType.String()))
}
