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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	animation.DrawGridStage(10, color.RGBA{R: 255, G: 0, B: 0, A: 255}, constants.StageColor, screen)
	for _, p := range g.players {
		graphics.Draw(p.Character, screen, g.camera)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.camera.Viewport.W), int(g.camera.Viewport.H)
}

func main() {
	initializeSystems()
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle("Fighting Game")

	player1 := player.CreateDebugPlayer()
	player1.Character.Position = types.Vector2{X: constants.WorldWidth / 2, Y: constants.WorldHeight - 200}

	game := &Game{
		players: []*player.Player{player1},
		camera:  graphics.NewDefaultCamera(),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
