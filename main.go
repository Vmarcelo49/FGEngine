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

	renderQueue *graphics.RenderQueue
	debugui     debugui.DebugUI
}

func (g *Game) Update() error {
	logic.UpdateByInputs([]input.GameInput{g.players[0].Input.GetLocalInputs()}, g.players)
	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Game Debug Info", image.Rect(0, 0, 256, 144), func(layout debugui.ContainerLayout) {
			ctx.Text("Camera Info:")
			ctx.Text(fmt.Sprintf("Position: (%.2f, %.2f)", g.camera.Viewport.X, g.camera.Viewport.Y))
			ctx.Text(fmt.Sprintf("Size: (%.2f, %.2f)", g.camera.Viewport.W, g.camera.Viewport.H))
			ctx.Text("Character Info:")
			for i, p := range g.players {
				ctx.Text(fmt.Sprintf("Player %d:", i+1))
				ctx.Text(fmt.Sprintf("Position: (%.2f, %.2f)", p.Character.GetPosition().X, p.Character.GetPosition().Y))
				ctx.Text(fmt.Sprintf("State: %s", p.Character.StateMachine.ActiveState.String()))
			}
		})
		return nil
	}); err != nil {
		return err
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camera.Viewport.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camera.Viewport.X += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camera.Viewport.Y -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camera.Viewport.Y += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.camera.Scaling *= 1.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.camera.Scaling *= 0.99
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.camera.Viewport.AlignCenter(constants.World)
		g.camera.Scaling = 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	animation.DrawStaticImageStage(g.stageImg, screen, g.camera.WorldToScreen(types.Vector2{X: 0, Y: 0}), g.camera.Scaling)
	g.renderQueue.Draw(screen, g.camera)
	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.LayoutSizeW, config.LayoutSizeH
}

func main() {
	config.InitGameConfig()
	ebiten.SetWindowTitle("Fighting Game")
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)

	player1 := player.NewDebugPlayer()
	// Center player horizontally in the world, position vertically on ground
	playerRect := player1.Character.GetSprite().Rect
	playerRect.AlignCenter(constants.World)
	player1.Character.StateMachine.Position = types.Vector2{X: playerRect.X, Y: constants.GroundLevelY - float64(player1.Character.GetSprite().Rect.H)}

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
