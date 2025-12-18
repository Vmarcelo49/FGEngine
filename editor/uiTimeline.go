package editor

import (
	"fgengine/config"
	"fgengine/filepicker"
	"fmt"
	"image"

	"github.com/ebitengine/debugui"
)

func (g *Game) guiTimeline(ctx *debugui.Context) {
	topY := config.WindowHeight - framePanelHeight
	rightX := config.WindowWidth - panelWidth - 1

	ctx.Window("Timeline", image.Rect(panelWidth, topY, rightX, config.WindowHeight), func(layout debugui.ContainerLayout) {
		sprite := g.character.AnimationPlayer.ActiveSprite()

		if sprite == nil {
			ctx.Text("No frame selected")
			return
		}
		ctx.SetGridLayout([]int{100, -1, 60}, nil)

		ctx.Text("Active Frame:")
		frameCount := len(g.ActiveAnimation().FrameData)
		if frameCount > 0 {
			ctx.Slider(&g.character.AnimationPlayer.FrameIndex, 0, frameCount-1, 1).On(func() {
				// Reset frame time when manually changing frame
				g.character.AnimationPlayer.FrameTimeLeft = g.ActiveAnimation().FrameData[g.character.AnimationPlayer.FrameIndex].Duration
			})
		}
		framedataIndex := g.character.AnimationPlayer.FrameIndex
		ctx.Text(fmt.Sprintf("%d / %d", framedataIndex+1, frameCount))

		ctx.Text("Framedata:")
		framedataLen := len(g.ActiveAnimation().FrameData)
		if framedataLen > 1 {
			ctx.Slider(&g.uiVariables.frameDataIndex, 0, framedataLen-1, 1).On(func() {
				g.character.AnimationPlayer.FrameIndex = g.uiVariables.frameDataIndex
				g.character.AnimationPlayer.FrameTimeLeft = g.ActiveAnimation().FrameData[g.uiVariables.frameDataIndex].Duration
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
			g.character.AnimationPlayer.FrameIndex = 0
			g.character.AnimationPlayer.FrameTimeLeft = g.ActiveAnimation().FrameData[0].Duration
		})

	})
}

func (g *Game) AddImageToAnimation() {
	picker := filepicker.GetFilePicker()
	filter := filepicker.FileFilter{
		Description: "Image files",
		Extensions:  []string{"png"},
	}

	path, err := picker.LoadFile(filter)
	if err != nil {
		g.writeLog(fmt.Sprintf("failed to load image: %s", err))
		return
	}

	if err := g.addSpriteByFile(path); err != nil {
		g.writeLog(fmt.Sprintf("failed to add image to frame: %s", err))
		return
	}
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
	lastIndex := len(g.ActiveAnimation().FrameData) - 1
	frameData := g.ActiveAnimation().FrameData
	g.ActiveAnimation().FrameData = append(frameData[:g.uiVariables.frameDataIndex], frameData[g.uiVariables.frameDataIndex+1:]...)

	// Adjust frameIndex after removal
	if g.uiVariables.frameDataIndex > 0 {
		g.uiVariables.frameDataIndex = lastIndex - 1
	}

	if len(g.ActiveAnimation().FrameData) == 0 {
		g.uiVariables.frameDataIndex = 0
	}
}
