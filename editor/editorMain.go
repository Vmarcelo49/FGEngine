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
	debugui         debugui.DebugUI
	activeCharacter *character.Character
	editorManager   *EditorManager
	inputManager    *input.InputManager
	lastMouseX      int
	lastMouseY      int
	isDragging      bool
	camera          *graphics.Camera
}

func (g *Game) Update() error {
	g.handleCameraInput()
	g.handleBoxMouseEdit()
	g.updateAnimationPlayback() // Add animation playback logic
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
	return config.LayoutSizeW, config.LayoutSizeH
}

func Run() {
	config.SetEditorConfig()
	ebiten.SetWindowTitle("Animation Editor")
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)

	game := &Game{
		editorManager: &EditorManager{
			logBuf: "Move Camera with Right click drag\n",
		},
		inputManager: input.NewInputManager(),
		camera:       graphics.NewCamera(),
	}
	game.camera.Scaling = float64(config.LayoutSizeW) / constants.Camera.W
	game.camera.SetPosition(types.Vector2{X: (-constants.World.W / 2), Y: (-constants.World.H / 2)})

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

// updateAnimationPlayback handles automatic frame advancement when playingAnim is true
func (g *Game) updateAnimationPlayback() {
	if !g.editorManager.playingAnim || g.editorManager.activeAnimation == nil {
		return
	}

	g.editorManager.frameCounter++

	// Get current frame duration
	if g.editorManager.frameIndex < len(g.editorManager.activeAnimation.FrameData) {
		currentFrameDuration := g.editorManager.activeAnimation.FrameData[g.editorManager.frameIndex].Duration

		// Check if current frame duration has elapsed (Duration is in game frames, not seconds)
		if g.editorManager.frameCounter >= currentFrameDuration {
			g.editorManager.frameCounter = 0
			g.editorManager.frameIndex++

			// Handle end of animation - loop back to start
			if g.editorManager.frameIndex >= len(g.editorManager.activeAnimation.FrameData) {
				g.editorManager.frameIndex = 0
			}

			// AnimationPlayer automatically handles sprite selection
		}
	}
}
