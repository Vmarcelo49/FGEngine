package main

import (
	"fgengine/animation"
	"fgengine/camera"
	"fgengine/config"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/player"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	players          []*player.Player
	animationManager *animation.AnimationRegistry
	inputManager     *input.InputManager
}

func (g *Game) Update() error {
	g.animationManager.UpdateAll()
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

	animManager := animation.NewAnimationRegistry()
	inputManager := input.NewInputManager()

	player1 := player.CreateDebugPlayer(animManager)

	game := &Game{
		players:          []*player.Player{player1},
		animationManager: animManager,
		inputManager:     inputManager,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
