package scene

import (
	"fgengine/constants"
	"fgengine/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update([2]input.Input) SceneStatus
	Draw(*ebiten.Image)
}

type SceneStatus byte

const (
	SceneDontChange SceneStatus = iota
	Scene1
	Scene2
	SceneController
)

type SceneManager struct {
	currentScene Scene
	//inputs
}

func NewSceneManager() *SceneManager {
	return &SceneManager{
		currentScene: MakeControllerScene(),
	}
}

func (sm *SceneManager) Update() error {
	input.UpdateGamepads()
	sceneSignal := sm.currentScene.Update(input.CurrentInputs())
	switch sceneSignal {
	case Scene1:
		//sm.currentScene = MakeScene1()
	case Scene2:
		//sm.currentScene = MakeScene2()
	case SceneController:
		sm.currentScene = MakeControllerScene()
	}
	return nil
}

func (sm *SceneManager) Draw(screen *ebiten.Image) {
	sm.currentScene.Draw(screen)
}

func (sm *SceneManager) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(constants.CameraWidth), int(constants.CameraHeight)
}
