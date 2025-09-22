package main

import (
	"fgengine/config"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/player"
	"fgengine/types"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	players      []*player.Player
	inputManager *input.InputManager
	camera       *graphics.Camera
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
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

	inputManager := input.NewInputManager()

	player1 := player.CreateDebugPlayer()

	game := &Game{
		players:      []*player.Player{player1},
		inputManager: inputManager,
		camera: &graphics.Camera{
			Viewport: types.Rect{
				X: 0,
				Y: 0,
				W: float64(config.WindowWidth),
				H: float64(config.WindowHeight),
			},
			LockWorldBounds: true,
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
