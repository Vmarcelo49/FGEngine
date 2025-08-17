package editor

import (
	"fgengine/config"
	"image"

	"github.com/ebitengine/debugui"
)

const (
	panelWidth    = 320
	HeaderHeight  = 32
	toolbarHeight = 52

	leftPanelX   = 0
	rightPanelX  = -panelWidth
	bottomPanely = -200

	toolBarButtonWidth = 100
)

func (g *Game) updateDebugUI() error {
	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		g.guiMainToolbar(ctx)
		g.guiProjectPanel(ctx)

		if g.editorManager.activeAnimation != nil {
			g.guiFrameProperties(ctx)
			g.guiFramePanel(ctx)
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
	// g.editorManager.boxRenderer = nil // TODO, we have to check how the current box renderer interacts with the editor
	g.editorManager.uiPrevAnimationName = ""
	g.writeLog("Cleared current state")
}

func (g *Game) guiMainToolbar(ctx *debugui.Context) {
	ctx.Window("Toolbar", image.Rect(0, 0, config.WindowWidth, toolbarHeight), func(layout debugui.ContainerLayout) {
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

func (g *Game) guiProjectPanel(ctx *debugui.Context) {
}

func (g *Game) guiFrameProperties(ctx *debugui.Context) {
}

func (g *Game) guiFramePanel(ctx *debugui.Context) {
}

func (g *Game) logWindow(ctx *debugui.Context) {
}
