package editor

import (
	"fgengine/animation"
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
		sprite := g.editorManager.getCurrentSprite()
		spriteIndex := g.getActiveSpriteIndex()

		if sprite == nil {
			ctx.Text("No frame selected")
			return
		}
		ctx.SetGridLayout([]int{100, -1, 60}, nil)
		ctx.Text("Navigate:")
		frameCount := len(g.editorManager.activeAnimation.Sprites)
		frameIndex := int(g.editorManager.frameIndex)
		ctx.Slider(&frameIndex, 0, frameCount-1, 1).On(func() {
			g.editorManager.frameIndex = frameIndex
			// Update character's active sprite to match the current frame
			if g.activeCharacter != nil && g.editorManager.activeAnimation != nil {
				currentSprite := g.editorManager.getCurrentSprite()
				if currentSprite != nil {
					g.activeCharacter.ActiveSprite = currentSprite
				}
			}
			g.refreshBoxEditor() // Refresh box editor when frame changes
		})
		ctx.Text(fmt.Sprintf("%d / %d", g.editorManager.frameIndex+1, frameCount))

		ctx.SetGridLayout([]int{-1, 0, -1, -1, -1, -1}, nil)

		ctx.Text("Frame Duration:")
		duration := int(g.editorManager.activeAnimation.Prop[spriteIndex].Duration)
		ctx.NumberField(&duration, 1)
		if duration < 1 {
			duration = 1
		}
		g.editorManager.activeAnimation.Prop[spriteIndex].Duration = uint(duration)

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
			} else {
				g.editorManager.playingAnim = false
			}
		})

	})
}

func (g *Game) getActiveSpriteIndex() int {
	if g.editorManager.activeAnimation == nil || g.editorManager.frameIndex < 0 || g.editorManager.frameIndex >= len(g.editorManager.activeAnimation.Sprites) {
		return -1
	}
	animation := g.editorManager.activeAnimation
	for i, sprite := range animation.Sprites {
		if sprite == g.editorManager.getCurrentSprite() {
			return i
		}
	}
	return -1
}

func (g *Game) AddImageToFrame() {
	if g.editorManager.activeAnimation == nil {
		return
	}

	path, err := dialog.File().Filter("Image files", "png").Load()
	if err != nil {
		g.writeLog(fmt.Sprintf("failed to load image: %s", err))
		return
	}

	if err := g.editorManager.addSpriteByFile(path); err != nil {
		g.writeLog(fmt.Sprintf("failed to add image to frame: %s", err))
		return
	}

	g.editorManager.frameCount = len(g.editorManager.activeAnimation.Sprites)
	g.editorManager.frameIndex = g.editorManager.frameCount - 1
	g.editorManager.activeAnimation.Prop = append(g.editorManager.activeAnimation.Prop, animation.FrameProperties{})
}

func (g *Game) copyLastFrame() {
	if g.editorManager.activeAnimation == nil {
		return
	}
	lastFrameIndex := len(g.editorManager.activeAnimation.Sprites) - 1
	if lastFrameIndex < 0 {
		return
	}

	lastFrame := g.editorManager.activeAnimation.Sprites[lastFrameIndex]
	newFrame := deepCopySprite(lastFrame)

	g.editorManager.activeAnimation.Sprites = append(g.editorManager.activeAnimation.Sprites, newFrame)

	if lastFrameIndex < len(g.editorManager.activeAnimation.Prop) {
		lastProp := g.editorManager.activeAnimation.Prop[lastFrameIndex]
		g.editorManager.activeAnimation.Prop = append(g.editorManager.activeAnimation.Prop, lastProp)
	} else {
		// If no properties exist for the last frame, add an empty one
		g.editorManager.activeAnimation.Prop = append(g.editorManager.activeAnimation.Prop, animation.FrameProperties{})
	}
}

func (g *Game) removeFrame() {
	if g.editorManager.activeAnimation == nil || len(g.editorManager.activeAnimation.Sprites) == 0 {
		return
	}

	frames := g.editorManager.activeAnimation.Sprites
	if g.editorManager.frameIndex >= 0 && g.editorManager.frameIndex < len(frames) {
		g.editorManager.activeAnimation.Sprites = append(frames[:g.editorManager.frameIndex], frames[g.editorManager.frameIndex+1:]...)

		// Remove the corresponding property at the same index
		if g.editorManager.frameIndex < len(g.editorManager.activeAnimation.Prop) {
			props := g.editorManager.activeAnimation.Prop
			g.editorManager.activeAnimation.Prop = append(props[:g.editorManager.frameIndex], props[g.editorManager.frameIndex+1:]...)
		}
	}

	if g.editorManager.frameIndex > 0 {
		g.editorManager.frameIndex--
	}

	if len(g.editorManager.activeAnimation.Sprites) == 0 {
		g.editorManager.frameIndex = 0
	}
}
