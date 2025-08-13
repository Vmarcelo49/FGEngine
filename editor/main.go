package editor

import (
	"github.com/ebitengine/debugui"
)

const (
	WindowWidth  = 1600
	WindowHeight = 900
)

type Game struct {
	debugui debugui.DebugUI
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

/*
func main() {
	ebiten.SetWindowTitle("Animation Editor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(WindowWidth, WindowHeight)

	game := &Game{
		imageCache:        newImageCache(),
		zoom:              DefaultZoom,
		zoomSelectorIndex: 1,
		anim:              AnimationState{},
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
*/
