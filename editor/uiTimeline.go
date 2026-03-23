package editor

import (
	"fgengine/config"
	"fmt"
	"image"

	"github.com/ebitengine/debugui"
	"github.com/sqweek/dialog"
)

func (g *Game) guiTimeline(ctx *debugui.Context) {
	topY := config.WindowHeight - framePanelHeight
	rightX := config.WindowWidth - panelWidth - 1

	ctx.Window("Timeline", image.Rect(panelWidth, topY, rightX, config.WindowHeight), func(layout debugui.ContainerLayout) {
		player := g.activeAnimPlayer()
		if player == nil {
			ctx.Text("No animation player available")
			return
		}

		sprite := player.ActiveSprite()

		if sprite == nil {
			ctx.Text("No frame selected")
			return
		}
		ctx.SetGridLayout([]int{100, -1, 60}, nil)

		ctx.Text("Active Frame:")
		frameCount := len(g.ActiveAnimation().FrameData)
		if frameCount == 0 {
			ctx.Text("No frame data")
			return
		}
		if player.FrameIndex < 0 || player.FrameIndex >= frameCount {
			player.FrameIndex = 0
		}
		if g.uiVariables.frameDataIndex < 0 || g.uiVariables.frameDataIndex >= frameCount {
			g.uiVariables.frameDataIndex = player.FrameIndex
		}

		if frameCount > 0 {
			ctx.Slider(&player.FrameIndex, 0, frameCount-1, 1).On(func() {
				// Reset frame time when manually changing frame
				player.FrameTimeLeft = g.ActiveAnimation().FrameData[player.FrameIndex].Duration
				g.uiVariables.frameDataIndex = player.FrameIndex
			})
		}
		framedataIndex := player.FrameIndex
		ctx.Text(fmt.Sprintf("%d / %d", framedataIndex+1, frameCount))

		ctx.Text("Framedata:")
		framedataLen := len(g.ActiveAnimation().FrameData)
		if framedataLen > 1 {
			ctx.Slider(&g.uiVariables.frameDataIndex, 0, framedataLen-1, 1).On(func() {
				player.FrameIndex = g.uiVariables.frameDataIndex
				player.FrameTimeLeft = g.ActiveAnimation().FrameData[g.uiVariables.frameDataIndex].Duration
			})
		}
		ctx.Text(fmt.Sprintf("%d / %d", g.uiVariables.frameDataIndex+1, framedataLen))

		ctx.SetGridLayout([]int{-1, 0, -1, -1, -1, -1}, nil)

		ctx.Text("Frame Duration:")
		duration := g.ActiveAnimation().FrameData[framedataIndex].Duration
		ctx.NumberField(&duration, 1).On(func() {
			if duration < 1 {
				duration = 1
			}
			g.ActiveAnimation().FrameData[framedataIndex].Duration = duration
		})

		ctx.Button("Add Frame by Image").On(func() {
			g.AddImageToAnimation()
		})
		ctx.Button("Duplicate Frame").On(func() {
			g.duplicateLastFrameData()
		})
		ctx.Button("Remove Frame").On(func() {
			g.removeFrame()
		})
		playPauseToggleText := "Play"
		if g.uiVariables.playingAnim {
			playPauseToggleText = "Stop"
		}
		ctx.Button(playPauseToggleText).On(func() {
			if playPauseToggleText == "Play" {
				g.uiVariables.playingAnim = true
			} else {
				g.uiVariables.playingAnim = false
			}
			player.FrameIndex = 0
			player.FrameTimeLeft = g.ActiveAnimation().FrameData[0].Duration
		})

	})
}

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
