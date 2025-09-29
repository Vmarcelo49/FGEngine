package main

import (
	"fgengine/animation"
	"fgengine/config"
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/logic"
	"fgengine/player"
	"fgengine/types"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	players []*player.Player
	camera  *graphics.Camera
}

func (g *Game) Update() error {
	logic.UpdateByInputs([]input.GameInput{g.players[0].InputManager.GetLocalInputs()}, g.players)
	g.camera.UpdatePosition(types.Vector2{X: g.players[0].Character.GetPosition().X, Y: g.players[0].Character.GetPosition().Y})
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	animation.DrawGridStage(20, color.RGBA{R: 255, G: 0, B: 0, A: 255}, constants.StageColor, screen)
	for _, p := range g.players {
		graphics.Draw(p.Character, screen, g.camera)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.LayoutSizeW, config.LayoutSizeH
}

func main() {
	config.InitDefaultConfig()
	ebiten.SetWindowTitle("Fighting Game")
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)

	player1 := player.NewDebugPlayer()
	player1.Character.StateMachine.Position = types.Vector2{X: constants.WorldWidth / 2, Y: constants.WorldHeight - 200}

	game := &Game{
		players: []*player.Player{player1},
		camera:  graphics.NewCamera(),
	}
	game.camera.LockWorldBounds = true

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
