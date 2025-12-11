package editor

// Editor can edit characters and animations
// every animation has a set of properties and boxes that can be adjusted
// run go run .\cmd\editor\

import (
	"fgengine/character"
	"fgengine/config"
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/types"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	character    *character.Character
	uiVariables  *uiVariables
	camera       *graphics.Camera
	inputManager *input.InputManager
	mouse        *MouseInput
	debugui      debugui.DebugUI
}

type MouseInput struct {
	lastMouseX int
	lastMouseY int
	isDragging bool
}

func (g *Game) Update() error {
	g.handleCameraInput()
	g.handleBoxMouseEdit()
	g.updateAnimationFrame()
	if err := g.updateDebugUI(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.character != nil && g.character.AnimationPlayer.ActiveAnimation != nil {
		//graphics.Draw(g.character, screen, g.camera)
		g.character.Draw(screen, g.camera)
		graphics.DrawBoxes(g.character.AnimationPlayer.GetActiveFrameData(), screen, g.camera, g.character)
	}

	g.drawMouseCrosshair(screen)
	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.LayoutSizeW, config.LayoutSizeH
}

func MakeEditorGame() *Game {
	game := &Game{
		uiVariables: &uiVariables{
			logBuf:           "Move Camera with Right click drag\n",
			enableMouseInput: new(bool),
		},
		inputManager: input.NewInputManager(),
		camera:       graphics.NewCamera(),
	}
	game.camera.Scaling = float64(config.LayoutSizeW) / constants.Camera.W
	game.camera.SetPosition(types.Vector2{X: (-constants.World.W / 2), Y: (-constants.World.H / 2)})
	game.mouse = &MouseInput{}

	return game
}

func Run() {
	config.SetEditorConfig()
	ebiten.SetWindowTitle("Animation Editor")
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)

	game := MakeEditorGame()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
