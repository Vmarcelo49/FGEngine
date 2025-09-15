package editor

// Editor can edit characters and animations
// every animation has a set of properties and boxes that can be adjusted
// run go run .\cmd\editor\

import (
	"fgengine/character"
	"fgengine/config"
	"fgengine/constants"
	"fgengine/graphics"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	Zoom float64
)

type Game struct {
	debugui         debugui.DebugUI
	activeCharacter *character.Character
	editorManager   *EditorManager
}

func (g *Game) Update() error {
	// g.handleMouseInput()
	if err := g.updateDebugUI(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.editorManager.activeAnimation != nil && g.activeCharacter != nil {
		graphics.DrawRenderableWithScale(g.activeCharacter, screen, Zoom, Zoom)
		graphics.DrawBoxesOf(g.activeCharacter, screen)
	}
	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.WindowWidth, config.WindowHeight
}

func Run() {
	config.InitDefaultConfig()
	Zoom = float64(config.WindowWidth) / constants.WorldWidth // 2.5 since WorldWidth is 640
	ebiten.SetWindowTitle("Animation Editor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)

	game := &Game{
		editorManager: &EditorManager{
			logBuf: "Animation Editor Started",
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
