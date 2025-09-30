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
	"fmt"
	"image"
	"log"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	players  []*player.Player
	camera   *graphics.Camera
	stageImg *ebiten.Image

	debugui debugui.DebugUI
}

func (g *Game) Update() error {
	logic.UpdateByInputs([]input.GameInput{g.players[0].InputManager.GetLocalInputs()}, g.players)
	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Game Debug Info", image.Rect(0, 0, 256, 144), func(layout debugui.ContainerLayout) {
			ctx.Text("Camera Info:")
			ctx.Text(fmt.Sprintf(" Position: (%.2f, %.2f)", g.camera.Viewport.X, g.camera.Viewport.Y))
			ctx.Text(fmt.Sprintf(" Size: (%.2f, %.2f)", g.camera.Viewport.W, g.camera.Viewport.H))
			ctx.Text("Character Info:")
			for i, p := range g.players {
				ctx.Text(fmt.Sprintf(" Player %d:", i+1))
				ctx.Text(fmt.Sprintf("  Position: (%.2f, %.2f)", p.Character.GetPosition().X, p.Character.GetPosition().Y))
				ctx.Text(fmt.Sprintf("  State: %s", p.Character.StateMachine.ActiveState.String()))
			}
		})
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	animation.DrawStaticImageStage(g.stageImg, screen, g.camera.WorldToScreen(types.Vector2{X: 0, Y: 0}))
	for _, p := range g.players {
		graphics.Draw(p.Character, screen, g.camera)
	}
	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.LayoutSizeW, config.LayoutSizeH
}

func main() {
	config.InitDefaultConfig()
	ebiten.SetWindowTitle("Fighting Game")
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)

	player1 := player.NewDebugPlayer()
	player1.Character.StateMachine.Position = types.Vector2{X: (constants.WorldWidth - player1.Character.ActiveSprite.Rect.W) / 2, Y: constants.GroundLevelY - float64(player1.Character.ActiveSprite.Rect.H)} // start in the middle of the world

	game := &Game{
		players: []*player.Player{player1},
		camera:  graphics.NewCamera(),
	}
	game.camera.LockWorldBounds = true

	game.stageImg, _, _ = ebitenutil.NewImageFromFile("assets/stages/PlaceMarkers.png")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
