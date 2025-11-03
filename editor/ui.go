package editor

import (
	"fgengine/config"
	"fgengine/state"
	"fmt"
	"image"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	panelWidth    = 320
	HeaderHeight  = 32
	toolbarHeight = 52

	leftPanelX       = 0
	rightPanelX      = -panelWidth
	bottomPanely     = -200
	framePanelHeight = 200

	toolBarButtonWidth = 100
)

func (g *Game) updateDebugUI() error {
	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		g.uiToolbar(ctx)
		g.uiProjectPanel(ctx)

		if g.editorManager.activeAnimation != nil {
			g.uiFrameProperties(ctx)
			g.guiTimeline(ctx)
		}

		g.logWindow(ctx)

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (g *Game) writeLog(text string) {
	if len(g.editorManager.logBuf) > 0 {
		g.editorManager.logBuf += "\n"
	}
	g.editorManager.logBuf += text
	g.editorManager.logUpdated = true
}

func (g *Game) resetCharacterState() {
	g.writeLog("There's a character loaded, clearing current state")
	g.activeCharacter = nil
	g.editorManager.activeAnimation = nil
	g.editorManager.previousAnimationName = ""
	g.writeLog("Cleared current state")
}

func (g *Game) uiToolbar(ctx *debugui.Context) {
	ctx.Window("Toolbar", image.Rect(0, 0, config.WindowWidth, toolbarHeight-1), func(layout debugui.ContainerLayout) {
		ctx.SetGridLayout([]int{toolBarButtonWidth, toolBarButtonWidth, toolBarButtonWidth, toolBarButtonWidth, -1, 35, 70}, nil)
		ctx.Button("New Character").On(func() {
			g.createCharacter()
		})
		ctx.Button("Load Character").On(func() {
			g.loadCharacter()
		})
		if g.activeCharacter != nil {
			ctx.Button("Save Character").On(func() {
				g.saveCharacter()
			})
		}
	})
}

func (g *Game) logWindow(ctx *debugui.Context) {
	ctx.Window("Log Window", image.Rect(panelWidth+1, toolbarHeight, 650, 290), func(layout debugui.ContainerLayout) {
		ctx.SetGridLayout([]int{-1}, []int{-1, 0})
		ctx.Panel(func(layout debugui.ContainerLayout) {
			ctx.SetGridLayout([]int{-1}, []int{-1})
			ctx.Text(g.editorManager.logBuf)
			if g.editorManager.logUpdated {
				ctx.SetScroll(image.Pt(layout.ScrollOffset.X, layout.ContentSize.Y))
				g.editorManager.logUpdated = false
			}
		})
		ctx.GridCell(func(bounds image.Rectangle) {
			submit := func() {
				if g.editorManager.logSubmitBuf == "" {
					return
				}
				g.writeLog(g.editorManager.logSubmitBuf)
				g.editorManager.logSubmitBuf = ""
			}
			ctx.SetGridLayout([]int{-3, -1}, nil)
			ctx.TextField(&g.editorManager.logSubmitBuf).On(func() {
				if ebiten.IsKeyPressed(ebiten.KeyEnter) {
					submit()
					ctx.SetTextFieldValue(g.editorManager.logSubmitBuf)
				}
			})
			ctx.Button("Submit").On(func() {
				submit()
			})
		})
	})
}

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
			if frameIndex >= 0 && frameIndex < len(g.editorManager.activeAnimation.FrameData) {
				properties := &g.editorManager.activeAnimation.FrameData[frameIndex]

				// Add SpriteIndex editor
				ctx.Text("Sprite Index:")
				spriteIndex := properties.SpriteIndex
				maxSpriteIndex := len(g.editorManager.activeAnimation.Sprites) - 1
				if maxSpriteIndex < 0 {
					maxSpriteIndex = 0
				}
				ctx.Slider(&spriteIndex, 0, maxSpriteIndex, 1).On(func() {
					if spriteIndex >= 0 && spriteIndex < len(g.editorManager.activeAnimation.Sprites) {
						properties.SpriteIndex = spriteIndex
						// AnimationPlayer automatically handles sprite selection
					}
				})
				ctx.Text(fmt.Sprintf("Points to sprite: %d / %d", spriteIndex+1, len(g.editorManager.activeAnimation.Sprites)))

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
				if len(g.editorManager.activeAnimation.FrameData) != len(g.editorManager.activeAnimation.Sprites) {
					ctx.Text(fmt.Sprintf("error? currently there are %d props, but the number of frames is %d", len(g.editorManager.activeAnimation.FrameData), len(g.editorManager.activeAnimation.Sprites)))
				}
				if len(g.editorManager.activeAnimation.FrameData) == 0 && len(g.editorManager.activeAnimation.Sprites) == 0 {
					ctx.Text("no sprites or properties found")
				}
			}
		})
	})
}
