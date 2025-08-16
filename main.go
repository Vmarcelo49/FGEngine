package main

import (
	"fgengine/animation"
	"fgengine/camera"
	"fgengine/config"
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/player"
	"fgengine/types"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	players          []*player.Player
	animationManager *animation.AnimationRegistry
	inputManager     *input.InputManager // Needed to make logging the inputs work properly.
}

func (g *Game) Update() error {
	g.inputManager.UpdateGamepadList() // logging user input
	g.inputManager.GetLocalInputs()	   //
	g.animationManager.UpdateAll()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		graphics.DrawRenderable(p, screen)
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
	player1.AnimationComponent.SetAnimation("idle")
	player1.Position = types.Vector2{X: constants.WorldWidth / 2, Y: constants.WorldHeight / 2}

	game := &Game{
		players:          []*player.Player{player1},
		animationManager: animManager,
		inputManager:     inputManager,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
