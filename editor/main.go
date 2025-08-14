package editor

import (
	"fgengine/animation"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WindowWidth  = 1600
	WindowHeight = 900
)

type Game struct {
	debugui debugui.DebugUI

	activeCharacter *animation.Character

	editorManager *EditorManager
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

func main() {
	ebiten.SetWindowTitle("Animation Editor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(WindowWidth, WindowHeight)

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
