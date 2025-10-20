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
		sprite := g.editorManager.getCurrentSprite()

		if sprite == nil {
			ctx.Text("No frame selected")
			return
		}
		ctx.SetGridLayout([]int{100, -1, 60}, nil)
		ctx.Text("Navigate:")
		frameCount := len(g.editorManager.activeAnimation.FrameData)
		frameIndex := int(g.editorManager.frameIndex)
		if frameCount > 0 {
			ctx.Slider(&frameIndex, 0, frameCount-1, 1).On(func() {
				g.editorManager.frameIndex = frameIndex
				g.editorManager.frameCounter = 0 // Reset counter when manually changing frame
				// Update character's active sprite to match the current frame
				if g.activeCharacter != nil && g.editorManager.activeAnimation != nil {
					currentSprite := g.editorManager.getCurrentSprite()
					if currentSprite != nil {
						g.activeCharacter.ActiveSprite = currentSprite
					}
				}
				g.refreshBoxEditor() // Refresh box editor when frame changes
			})
		}
		ctx.Text(fmt.Sprintf("%d / %d", g.editorManager.frameIndex+1, frameCount))

		ctx.SetGridLayout([]int{-1, 0, -1, -1, -1, -1}, nil)

		ctx.Text("Frame Duration:")
		if g.editorManager.frameIndex >= 0 && g.editorManager.frameIndex < len(g.editorManager.activeAnimation.FrameData) {
			duration := int(g.editorManager.activeAnimation.FrameData[g.editorManager.frameIndex].Duration)
			ctx.NumberField(&duration, 1)
			if duration < 1 {
				duration = 1
			}
			g.editorManager.activeAnimation.FrameData[g.editorManager.frameIndex].Duration = duration

			ctx.Text("Animation Switch:")
			animationSwitch := g.editorManager.activeAnimation.FrameData[g.editorManager.frameIndex].AnimationSwitch
			ctx.TextField(&animationSwitch).On(func() {
				g.editorManager.activeAnimation.FrameData[g.editorManager.frameIndex].AnimationSwitch = animationSwitch
			})
		} else {
			ctx.Text("Invalid frame index")
		}

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
		if g.editorManager.playingAnim {
			playPauseToggleText = "Stop"
		}
		ctx.Button(playPauseToggleText).On(func() {
			if playPauseToggleText == "Play" {
				g.editorManager.playingAnim = true
				g.editorManager.frameCounter = 0 // Reset counter when starting playback
			} else {
				g.editorManager.playingAnim = false
				g.editorManager.frameCounter = 0 // Reset counter when stopping playback
			}
		})

	})
}

func (g *Game) AddImageToFrame() {
	if g.editorManager.activeAnimation == nil {
		return
	}

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

	if err := g.editorManager.addSpriteByFile(path); err != nil {
		g.writeLog(fmt.Sprintf("failed to add image to frame: %s", err))
		return
	}

	g.editorManager.frameCount = len(g.editorManager.activeAnimation.FrameData)
	g.editorManager.frameIndex = g.editorManager.frameCount - 1
	// Removed the duplicate FrameData append since addSpriteByFile already handles it
}

func (g *Game) copyLastFrame() {
	if g.editorManager.activeAnimation == nil {
		return
	}
	lastFrameIndex := len(g.editorManager.activeAnimation.FrameData) - 1
	if lastFrameIndex < 0 {
		return
	}

	// Get the sprite that the last frame points to
	lastFrameData := g.editorManager.activeAnimation.FrameData[lastFrameIndex]
	if lastFrameData.SpriteIndex >= 0 && lastFrameData.SpriteIndex < len(g.editorManager.activeAnimation.Sprites) {
		lastFrame := g.editorManager.activeAnimation.Sprites[lastFrameData.SpriteIndex]
		newFrame := deepCopySprite(lastFrame)

		g.editorManager.activeAnimation.Sprites = append(g.editorManager.activeAnimation.Sprites, newFrame)

		// Copy the last frame's properties and update SpriteIndex to point to the newly created sprite
		lastProp := lastFrameData
		lastProp.SpriteIndex = len(g.editorManager.activeAnimation.Sprites) - 1
		g.editorManager.activeAnimation.FrameData = append(g.editorManager.activeAnimation.FrameData, lastProp)
	}
}

func (g *Game) removeFrame() {
	if g.editorManager.activeAnimation == nil || len(g.editorManager.activeAnimation.FrameData) == 0 {
		return
	}

	// Remove the FrameData at current frameIndex
	if g.editorManager.frameIndex >= 0 && g.editorManager.frameIndex < len(g.editorManager.activeAnimation.FrameData) {
		props := g.editorManager.activeAnimation.FrameData
		g.editorManager.activeAnimation.FrameData = append(props[:g.editorManager.frameIndex], props[g.editorManager.frameIndex+1:]...)
	}

	// Adjust frameIndex after removal
	if g.editorManager.frameIndex > 0 {
		g.editorManager.frameIndex--
	}

	if len(g.editorManager.activeAnimation.FrameData) == 0 {
		g.editorManager.frameIndex = 0
	}
}
