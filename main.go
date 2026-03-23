package main

import (
	"fgengine/config"
	"fgengine/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	config.InitGameConfig()
	if err := ebiten.RunGame(scene.NewSceneManager()); err != nil {
		panic(err)
	}
}
