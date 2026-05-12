package main

import (
	"fgengine/editor"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("FGEngine Character Editor")
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(editor.NewCharacterEditor()); err != nil {
		panic(err)
	}
}
