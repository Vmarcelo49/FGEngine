package editor

import (
	"fgengine/animation"
	"fgengine/config"
	"fgengine/constants"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	Zoom = float64(config.WindowWidth) / constants.WorldWidth // 2.5 since WorldWidth is 640
)

type Game struct {
	debugui         debugui.DebugUI
	activeCharacter *animation.Character
	editorManager   *EditorManager
}

func (g *Game) Update() error {
	if err := g.updateDebugUI(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.editorManager.activeAnimation != nil { // this will be a fun rewrite
		//g.renderCurrentFrame(screen)
		//g.drawBoxes(screen)

	}
	g.debugui.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.WindowWidth, config.WindowHeight
}

func Run() {
	config.InitDefaultConfig()

	ebiten.SetWindowTitle("Animation Editor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)

	if err := ebiten.RunGame(&Game{
		editorManager: &EditorManager{
			logBuf: "Animation Editor Started",
		},
	}); err != nil {
		panic(err)
	}
}
