package scene

import (
	"fgengine/character"
	"fgengine/constants"
	"fgengine/gameplay"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/stage"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

func MakeGameplayScene() Scene {
	playerOne, err := character.LoadCharacter("PlaceHolder", 1)
	if err != nil {
		panic(err)
	}

	playerTwo, err := character.LoadCharacter("PlaceHolder", 2)
	if err != nil {
		panic(err)
	}

	camera := graphics.NewCamera()
	camera.LockWorldBounds = true

	return &GameplayScene{
		camera: camera,
		stage:  stage.NewSolidColorStage(constants.StageColor),
		gamestate: gameplay.GameState{
			Characters: [2]*character.Character{
				playerOne,
				playerTwo,
			}}}
}

type GameplayScene struct {
	camera    *graphics.Camera
	stage     *stage.Stage
	gamestate gameplay.GameState
}

func (g *GameplayScene) Update(inputs [2]input.GameInput) SceneStatus {
	g.gamestate.Update(inputs)
	g.updateCamera()
	return SceneDontChange
}

func (g *GameplayScene) Draw(screen *ebiten.Image) {
	if g.stage != nil {
		g.stage.Draw(screen, g.camera)
	}

	for _, char := range g.gamestate.Characters {
		if char == nil {
			continue
		}
		char.Draw(screen, g.camera)
	}
}

func (g *GameplayScene) updateCamera() {
	if g.camera == nil {
		return
	}

	left := g.gamestate.Characters[0]
	right := g.gamestate.Characters[1]
	if left == nil && right == nil {
		return
	}
	if left == nil {
		g.camera.UpdatePosition(right.Position())
		return
	}
	if right == nil {
		g.camera.UpdatePosition(left.Position())
		return
	}

	midpoint := types.Vector2{
		X: (left.Position().X + right.Position().X) / 2,
		Y: constants.WorldHeight / 2,
	}
	g.camera.UpdatePosition(midpoint)
}
