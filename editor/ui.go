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
	framePanelHeight = 200

	toolBarButtonWidth = 100
)

func (g *Game) updateDebugUI() error {
	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		g.uiToolbar(ctx)
		g.uiProjectPanel(ctx)
		if g.character != nil {
			if g.character.AnimationPlayer.ActiveAnimation != nil {
				g.uiFrameProperties(ctx)
				g.guiTimeline(ctx)
			}
		}

		g.logWindow(ctx)

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (g *Game) writeLog(text string) {
	if len(g.uiVariables.logBuf) > 0 {
		g.uiVariables.logBuf += "\n"
	}
	g.uiVariables.logBuf += text
	g.uiVariables.logUpdated = true
}

func (g *Game) resetCharacterState() {
	g.character = nil
	g.writeLog("There was a Character loaded, cleared current state")
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
		if g.character != nil {
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
			ctx.Text(g.uiVariables.logBuf)
			if g.uiVariables.logUpdated {
				ctx.SetScroll(image.Pt(layout.ScrollOffset.X, layout.ContentSize.Y))
				g.uiVariables.logUpdated = false
			}
		})
		ctx.GridCell(func(bounds image.Rectangle) {
			submit := func() {
				if g.uiVariables.logSubmitBuf == "" {
					return
				}
				g.writeLog(g.uiVariables.logSubmitBuf)
				g.uiVariables.logSubmitBuf = ""
			}
			ctx.SetGridLayout([]int{-3, -1}, nil)
			ctx.TextField(&g.uiVariables.logSubmitBuf).On(func() {
				if ebiten.IsKeyPressed(ebiten.KeyEnter) {
					submit()
					ctx.SetTextFieldValue(g.uiVariables.logSubmitBuf)
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
			frameData := g.character.AnimationPlayer.ActiveFrameData()
			if frameData == nil {
				ctx.Text("No frame selected")
				return
			}

			ctx.Text("Sprite Index:")

			lastSpriteIndex := len(g.character.AnimationPlayer.ActiveAnimation.Sprites) - 1
			if lastSpriteIndex < 0 {
				lastSpriteIndex = 0
			}

			ctx.Slider(&frameData.SpriteIndex, 0, lastSpriteIndex, 1) // Here is where we set the sprite index

			ctx.Text(fmt.Sprintf("Points to sprite: %d / %d", frameData.SpriteIndex+1, len(g.ActiveAnimation().Sprites)))

			ctx.SetGridLayout([]int{-1, -1, -1}, nil)
			ctx.Loop(len(state.OrderedStates), func(i int) { // Set checkboxes for each state
				stateFlag := state.OrderedStates[i]
				isActive := (frameData.State & stateFlag) != 0
				ctx.Checkbox(&isActive, stateFlag.String()).On(func() {
					if isActive {
						frameData.State |= stateFlag
					} else {
						frameData.State &= ^stateFlag
					}
				})
			})
			ctx.SetGridLayout(nil, nil)
			ctx.Text("Current State: " + frameData.State.String())
		})
	})
}
