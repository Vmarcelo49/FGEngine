package scene

import (
	"fgengine/constants"
	"fgengine/device"
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
	SceneMainMenu
	SceneControllerSelect
	SceneSelectCPU
	SceneSelectVersus
	SceneStartMatch
	SceneQuit
)

type MatchMode byte

const (
	MatchVersusCPU MatchMode = iota
	MatchLocalVersus
)

type SceneManager struct {
	currentScene Scene
	// After a scene swap, discard held input until both players are neutral.
	waitNeutral bool
}

func (sm *SceneManager) Update() error {
	polledInputs := device.UpdateGamepads()
	activeInputs := polledInputs

	if sm.waitNeutral {
		if polledInputs[0] == input.NoInput && polledInputs[1] == input.NoInput {
			sm.waitNeutral = false
		} else {
			activeInputs = [2]input.GameInput{input.NoInput, input.NoInput}
		}
	}

	sceneSignal := sm.currentScene.Update(activeInputs)
	switch sceneSignal {
	case SceneMainMenu:
		sm.currentScene = MakeMainMenuScene()
		sm.waitNeutral = true
	case SceneControllerSelect:
		sm.currentScene = MakeControllerScene()
		sm.waitNeutral = true
	case SceneSelectCPU:
		sm.currentScene = MakeCharacterSelectScene(MatchVersusCPU)
		sm.waitNeutral = true
	case SceneSelectVersus:
		sm.currentScene = MakeCharacterSelectScene(MatchLocalVersus)
		sm.waitNeutral = true
	case SceneStartMatch:
		if selection, ok := sm.currentScene.(*CharacterSelectScene); ok {
			sm.currentScene = MakeGameplayScene(selection.playerOneName(), selection.playerTwoName(), selection.mode == MatchVersusCPU)
			sm.waitNeutral = true
		} else {
			sm.currentScene = MakeMainMenuScene()
			sm.waitNeutral = true
		}
	case SceneQuit:
		return ebiten.Termination
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
