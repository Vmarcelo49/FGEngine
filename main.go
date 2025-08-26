package main

import (
	"fgengine/camera"
	"fgengine/config"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/player"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	players      []*player.Player
	inputManager *input.InputManager
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		graphics.DrawRenderable(p.Character, screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return camera.GetDimensions()
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
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
