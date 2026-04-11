package scene

import (
	"fgengine/constants"
	"fgengine/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update([2]input.GameInput) SceneStatus
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
	waitNeutral  bool
}

func (sm *SceneManager) Update() error {
	polledInputs := input.UpdateGamepads()
	activeInputs := polledInputs

	// Prevent button carry-over between scenes by waiting for full release.
	if sm.waitNeutral {
		if polledInputs[0] == input.NoInput && polledInputs[1] == input.NoInput {
			sm.waitNeutral = false
		} else {
			activeInputs = [2]input.GameInput{input.NoInput, input.NoInput}
		}
	}

	sceneSignal := sm.currentScene.Update(activeInputs)
	switch sceneSignal {
	case Scene1:
		sm.currentScene = MakeMainMenuScene()
		sm.waitNeutral = true
	case Scene2:
		sm.currentScene = MakeGameplayScene()
		sm.waitNeutral = true
	case SceneController:
		sm.currentScene = MakeControllerScene()
		sm.waitNeutral = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

func (sm *SceneManager) Draw(screen *ebiten.Image) {
	sm.currentScene.Draw(screen)
}

func (sm *SceneManager) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(constants.CameraWidth), int(constants.CameraHeight)
}

func NewSceneManager() *SceneManager {
	return &SceneManager{
		currentScene: MakeControllerScene(),
		waitNeutral:  true,
	}
}
