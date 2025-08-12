package main

import (
	"FGEngine/animation"
	"FGEngine/config"
	"FGEngine/graphics"
	"FGEngine/player"
	"FGEngine/types"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	players          []*player.Player
	animationManager *animation.ComponentManager
}

func (g *Game) Update() error {
	// Update all animation components each frame
	g.animationManager.UpdateAll()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		graphics.DrawRenderable(p, screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.CameraWidth, config.CameraHeight
}

func main() {
	initializeSystems()
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle("Fighting Game")

	// Create the animation manager
	animManager := animation.NewComponentManager()

	// Create player with the animation manager
	player1 := player.CreateDebugPlayer(animManager)
	player1.SetAnimation("idle")
	player1.Position = types.Vector2{X: config.WorldWidth / 2, Y: config.WorldHeight / 2}

	game := &Game{
		players:          []*player.Player{player1},
		animationManager: animManager,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
