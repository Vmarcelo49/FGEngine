package editor

import (
	"fgengine/config"
	"fgengine/state"
	"fmt"
	"image"

	"github.com/ebitengine/debugui"
)

func (g *Game) uiFrameProperties(ctx *debugui.Context) {
	ctx.Window("Properties", image.Rect(config.WindowWidth-panelWidth, toolbarHeight, config.WindowWidth, config.WindowHeight), func(layout debugui.ContainerLayout) {
		ctx.Header("Frame Info", true, func() {
			currentFrame := g.editorManager.getCurrentSprite()
			if currentFrame == nil {
				ctx.Text("No frame selected")
				return
			}
			// Get current frame properties
			frameIndex := g.editorManager.frameIndex
			if frameIndex >= 0 && frameIndex < len(g.editorManager.activeAnimation.Prop) {
				properties := &g.editorManager.activeAnimation.Prop[frameIndex]
				ctx.SetGridLayout([]int{-1, -1, -1}, nil)
				ctx.Loop(len(state.OrderedStates), func(i int) {
					stateFlag := state.OrderedStates[i]
					isActive := (properties.State & stateFlag) != 0
					ctx.Checkbox(&isActive, stateFlag.String()).On(func() {
						if isActive {
							properties.State |= stateFlag
						} else {
							properties.State &= ^stateFlag
						}
					})
				})
				ctx.SetGridLayout(nil, nil)
				ctx.Text("Current State: " + properties.State.String())

			} else {
				ctx.Text(fmt.Sprintf("currently there are %d props, but the number of frames is %d", len(g.editorManager.activeAnimation.Prop), len(g.editorManager.activeAnimation.Sprites)))
			}
		})
	})
}
