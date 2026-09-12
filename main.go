package main

import (
	"fgengine/config"
	"fgengine/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	config.InitGameConfig()
	ebiten.SetWindowSize(config.ActiveConfig.WindowWidth, config.ActiveConfig.WindowHeight)
	ebiten.SetWindowTitle("FG Engine")
	if err := ebiten.RunGame(scene.NewSceneManager()); err != nil {
		panic(err)
	}
}
