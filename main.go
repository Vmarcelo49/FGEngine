package main

import (
	"fgengine/constants"
	"fgengine/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(int(constants.CameraWidth), int(constants.CameraHeight))
	ebiten.SetWindowTitle("fgengine")
	if err := ebiten.RunGame(scene.NewSceneManager()); err != nil {
		panic(err)
	}
}
