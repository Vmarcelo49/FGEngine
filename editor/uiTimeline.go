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
		sprite := g.character.AnimationPlayer.GetSpriteFromFrameCounter()

		if sprite == nil {
			ctx.Text("No frame selected")
			return
		}
		ctx.SetGridLayout([]int{100, -1, 60}, nil)
		ctx.Text("Navigate:")
		frameCount := g.getActiveAnimation().Duration()
		if frameCount > 0 {
			ctx.Slider(&g.character.AnimationPlayer.FrameCounter, 0, frameCount-1, 1) // imagine if this works correctly from the start
		}
		ctx.Text(fmt.Sprintf("%d / %d", g.character.AnimationPlayer.GetActiveFrameDataIndex()+1, frameCount))

		ctx.SetGridLayout([]int{-1, 0, -1, -1, -1, -1}, nil)

		ctx.Text("Frame Duration:")
		duration := g.getActiveAnimation().FrameData[g.character.AnimationPlayer.GetActiveFrameDataIndex()].Duration
		ctx.NumberField(&duration, 1).On(func() {
			if duration < 1 {
				duration = 1
			}
			g.getActiveAnimation().FrameData[g.character.AnimationPlayer.GetActiveFrameDataIndex()].Duration = duration
		})

		ctx.Button("Add Image").On(func() {
			g.AddImageToFrame()
		})
		ctx.Button("Copy Last Frame").On(func() {
			g.copyLastFrame()
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
			g.character.AnimationPlayer.FrameCounter = 0
		})

	})
}

func (g *Game) AddImageToFrame() { // Rename this to AddImageToAnimation
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

func (g *Game) copyLastFrame() {
	lastFrameIndex := len(g.getActiveAnimation().FrameData) - 1
	if lastFrameIndex < 0 {
		return
	}
	lastFrameData := g.getActiveAnimation().FrameData[lastFrameIndex]
	g.getActiveAnimation().FrameData = append(g.getActiveAnimation().FrameData, lastFrameData)
}

func (g *Game) removeFrame() {
	if g.getActiveAnimation() == nil || len(g.getActiveAnimation().FrameData) == 0 {
		return
	}
	lastIndex := len(g.getActiveAnimation().FrameData) - 1
	frameData := g.getActiveAnimation().FrameData
	g.getActiveAnimation().FrameData = append(frameData[:g.uiVariables.frameDataIndex], frameData[g.uiVariables.frameDataIndex+1:]...)

	// Adjust frameIndex after removal
	if g.uiVariables.frameDataIndex > 0 {
		g.uiVariables.frameDataIndex = lastIndex - 1
	}

	if len(g.getActiveAnimation().FrameData) == 0 {
		g.uiVariables.frameDataIndex = 0
	}
}
