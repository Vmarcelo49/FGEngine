//go:build !js && !wasm
// +build !js,!wasm

package editorimgui

import (
	"fmt"

	"github.com/sqweek/dialog"
)

func (g *Game) AddImageToAnimation() {
	path, err := dialog.File().Filter("PNG Image", "png").Load()
	if err != nil {
		g.writeLog(fmt.Sprintf("failed to load image: %s", err))
		return
	}

	if err := g.addSpritesFromFiles([]string{path}); err != nil {
		g.writeLog(fmt.Sprintf("failed to add image to animation: %s", err))
		return
	}

	g.writeLog("Added frame to animation")
}

func (g *Game) duplicateLastFrameData() {
	lastFrameIndex := len(g.ActiveAnimation().FrameData) - 1
	if lastFrameIndex < 0 {
		return
	}
	lastFrameData := g.ActiveAnimation().FrameData[lastFrameIndex]
	g.ActiveAnimation().FrameData = append(g.ActiveAnimation().FrameData, lastFrameData)
}

func (g *Game) removeFrame() {
	if g.ActiveAnimation() == nil || len(g.ActiveAnimation().FrameData) == 0 {
		return
	}
	player := g.activeAnimPlayer()
	if player == nil {
		return
	}

	frameCount := len(g.ActiveAnimation().FrameData)
	if g.uiVariables.frameDataIndex < 0 || g.uiVariables.frameDataIndex >= frameCount {
		g.uiVariables.frameDataIndex = player.FrameIndex
	}
	if g.uiVariables.frameDataIndex < 0 || g.uiVariables.frameDataIndex >= frameCount {
		g.uiVariables.frameDataIndex = frameCount - 1
	}

	frameData := g.ActiveAnimation().FrameData
	g.ActiveAnimation().FrameData = append(frameData[:g.uiVariables.frameDataIndex], frameData[g.uiVariables.frameDataIndex+1:]...)

	if len(g.ActiveAnimation().FrameData) == 0 {
		player.FrameIndex = 0
		player.FrameTimeLeft = 0
		g.uiVariables.frameDataIndex = 0
		return
	}

	if g.uiVariables.frameDataIndex >= len(g.ActiveAnimation().FrameData) {
		g.uiVariables.frameDataIndex = len(g.ActiveAnimation().FrameData) - 1
	}
	player.FrameIndex = g.uiVariables.frameDataIndex
	player.FrameTimeLeft = g.ActiveAnimation().FrameData[player.FrameIndex].Duration
}
