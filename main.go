package main

import (
	"FGEngine/config"
	"FGEngine/graphics"
	"FGEngine/player"
	"FGEngine/types"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	player []*player.Player
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.player {
		graphics.DrawPlayer(p, screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.CameraWidth, config.CameraHeight
}

func main() {
	initializeSystems()
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle("Fighting Game")
	player1 := player.CreateDebugPlayer()
	player1.SetAnimation("idle")
	player1.State.Position = types.Vector2{X: config.WorldWidth / 2, Y: config.WorldHeight / 2}
	if err := ebiten.RunGame(&Game{player: []*player.Player{player1}}); err != nil {
		log.Fatal(err)
	}
}
