package scene

import (
	"fgengine/character"
	"fgengine/gameplay"
	"fgengine/input"

	"github.com/hajimehoshi/ebiten/v2"
)

func MakeGameplayScene() Scene {
	return &GameplayScene{
		gamestate: gameplay.GameState{
			Characters: [2]*character.Character{
				character.MakeTestCharacter(),
				character.MakeTestCharacter(),
			}}}
}

type GameplayScene struct {
	gamestate gameplay.GameState
}

func (g *GameplayScene) Update(inputs [2]input.GameInput) SceneStatus {
	g.gamestate.Update(inputs)
	return SceneDontChange
}

func (g *GameplayScene) Draw(screen *ebiten.Image) {
	// Placeholder for drawing the gameplay scene
}
