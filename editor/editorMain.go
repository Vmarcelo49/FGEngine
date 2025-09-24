package editor

// Editor can edit characters and animations
// every animation has a set of properties and boxes that can be adjusted
// run go run .\cmd\editor\

import (
	"fgengine/character"
	"fgengine/config"
	"fgengine/graphics"
	"fgengine/input"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	debugui         debugui.DebugUI
	activeCharacter *character.Character
	editorManager   *EditorManager
	inputManager    *input.InputManager
	lastMouseX      int
	lastMouseY      int
	isDragging      bool
	camera          *graphics.Camera
	zoom            float64
}

func (g *Game) Update() error {
	g.handleCameraInput()
	g.handleBoxMouseEdit()
	if err := g.updateDebugUI(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.editorManager.activeAnimation != nil && g.activeCharacter != nil {
		graphics.Draw(g.activeCharacter, screen, g.camera)
		graphics.DrawBoxes(g.activeCharacter, screen, g.camera)
	}

	g.drawMouseCrosshair(screen)

	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.WindowWidth, config.WindowHeight // in main its return int(g.camera.Viewport.W), int(g.camera.Viewport.H)
}

func Run() {
	config.InitDefaultConfig()
	camera := graphics.NewDefaultCamera() // with this, the character should not be visible at start, since it's at 0,0 and camera is at center of screen
	ebiten.SetWindowTitle("Animation Editor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)

	game := &Game{
		editorManager: &EditorManager{
			logBuf: "Move Camera with Right click drag\n",
		},
		inputManager: input.NewInputManager(),
		camera:       camera,
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
