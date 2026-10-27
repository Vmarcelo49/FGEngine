package scene

import (
	"fgengine/character"
	"fgengine/config"
	"fgengine/constants"
	"fgengine/gameplay"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/language"
	"fgengine/stage"
	"fgengine/types"
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func MakeGameplayScene(playerOneName, playerTwoName string, cpuOpponent bool) Scene {
	if playerOneName == "" {
		playerOneName = "PlaceHolder"
	}
	if playerTwoName == "" {
		playerTwoName = "PlaceHolder"
	}
	playerOne, err := character.LoadCharacter(playerOneName, 1)
	if err != nil {
		panic(err)
	}

	playerTwo, err := character.LoadCharacter(playerTwoName, 2)
	if err != nil {
		panic(err)
	}

	camera := graphics.NewCamera()
	camera.WorldBoundsLock = true

	// Locale loaded once per scene; missing files fall back to raw keys,
	// never a crash (see tr).
	var gameText map[string]string
	if locale, err := language.LoadLang(language.Lang(config.ActiveConfig.Language)); err != nil {
		fmt.Println("hud locale load failed:", err)
	} else {
		gameText = locale.GameText
	}

	return &GameplayScene{
		camera:      camera,
		stage:       stage.NewSolidColorStage(constants.StageColor),
		gamestate:   gameplay.NewGameState(playerOne, playerTwo, uint64(time.Now().UnixNano())),
		texts:       gameText,
		cpuOpponent: cpuOpponent,
	}
}

type GameplayScene struct {
	camera       *graphics.Camera
	stage        *stage.Stage
	gamestate    gameplay.GameState
	debugui      debugui.DebugUI
	texts        map[string]string
	cpuOpponent  bool
	debugVisible bool
}

func (g *GameplayScene) Update(inputs [2]input.GameInput) SceneStatus {
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.debugVisible = !g.debugVisible
	}
	if g.cpuOpponent {
		inputs[1] = g.cpuInput()
	}
	g.gamestate.Update(inputs)
	if g.gamestate.MatchOver() {
		return SceneMainMenu
	}
	g.updateCamera()
	if g.debugVisible {
		g.updateDebugUI()
	}
	return SceneDontChange
}

func (g *GameplayScene) Draw(screen *ebiten.Image) {
	if g.stage != nil {
		g.stage.Draw(screen, g.camera)
	}
	g.drawArena(screen)

	for _, char := range g.gamestate.Characters {
		if char == nil {
			continue
		}
		graphics.DrawCharacter(screen, char, g.camera)
	}

	g.drawHUD(screen)
	if g.cpuOpponent {
		ebitenutil.DebugPrintAt(screen, "MOVE WASD  /  ATTACK U I O J  /  F1 DEBUG", 12, 344)
	} else {
		ebitenutil.DebugPrintAt(screen, "P1 WASD + U/I/O/J     P2 ARROWS + N/M/,/.     F1 DEBUG", 12, 344)
	}
	if g.debugVisible {
		for _, char := range g.gamestate.Characters {
			if char != nil {
				graphics.DrawCharacterBoxes(screen, char, g.camera)
			}
		}
		g.drawDebugGuides(screen)
		g.debugui.Draw(screen)
		ebitenutil.DebugPrintAt(screen, "DEBUG  /  F1 TO HIDE", 12, 46)
	}
}

func (g *GameplayScene) drawArena(screen *ebiten.Image) {
	// The engine's test stage is a solid fill. These restrained shapes give
	// the playable slice a clear horizon, floor, and ring without affecting
	// the deterministic simulation or camera math.
	for i := 0; i < 5; i++ {
		y := float32(72 + i*37)
		vector.StrokeLine(screen, 0, y, 640, y, 1, color.RGBA{R: 27, G: 38, B: 59, A: 255}, false)
	}
	for i := 0; i < 9; i++ {
		x := float32(36 + i*72)
		vector.DrawFilledRect(screen, x, 68, 1, 210, color.RGBA{R: 25, G: 40, B: 59, A: 255}, false)
		vector.DrawFilledRect(screen, x-3, 96, 7, 2, color.RGBA{R: 244, G: 113, B: 70, A: 200}, false)
	}
	vector.DrawFilledRect(screen, 0, 311, 640, 49, color.RGBA{R: 10, G: 16, B: 29, A: 255}, false)
	vector.DrawFilledRect(screen, 0, 308, 640, 2, color.RGBA{R: 238, G: 111, B: 68, A: 255}, false)
	vector.DrawFilledRect(screen, 0, 311, 640, 1, color.RGBA{R: 246, G: 190, B: 104, A: 255}, false)
	for i := 0; i < 15; i++ {
		x := float32(i * 48)
		vector.StrokeLine(screen, 320, 277, x, 359, 1, color.RGBA{R: 22, G: 34, B: 52, A: 255}, false)
	}
	for i := 0; i < 4; i++ {
		y := float32(315 + i*10)
		vector.StrokeLine(screen, 0, y, 640, y, 1, color.RGBA{R: 20, G: 31, B: 47, A: 255}, false)
	}
}

func (g *GameplayScene) cpuInput() input.GameInput {
	if g.gamestate.Phase != gameplay.PhaseFight {
		return input.NoInput
	}
	p1 := g.gamestate.Characters[0].StateMachine
	p2 := g.gamestate.Characters[1].StateMachine
	if p1 == nil || p2 == nil || p2.StunFrames > 0 {
		return input.NoInput
	}

	distance := math.Abs(p1.Position.X - p2.Position.X)
	toward, away := input.Left, input.Right
	if p2.Position.X < p1.Position.X {
		toward, away = input.Right, input.Left
	}
	frame := constants.RoundTimerFrames - g.gamestate.TimerFrames

	if seq := cpuSpecialPattern(g.gamestate.Characters[1].Properties.SpecialInput); len(seq) > 0 {
		if slot := frame % 240; slot >= 0 && slot < len(seq) {
			return faceRelativeInput(seq[slot], bool(p2.IsFacingLeft))
		}
	}

	if distance > 190 {
		return toward
	}
	if distance > 135 && frame%84 < 12 {
		return toward
	}
	if distance < 154 && frame%132 >= 62 && frame%132 < 71 {
		return away
	}
	if distance <= 138 && frame%47 == 7 {
		return input.B
	}
	if distance <= 158 && frame%31 == 4 {
		return input.A
	}
	return input.NoInput
}

func cpuSpecialPattern(name string) []input.GameInput {
	if len(name) < 3 {
		return nil
	}
	var button input.GameInput
	switch name[len(name)-1:] {
	case "A":
		button = input.A
	case "B":
		button = input.B
	case "C":
		button = input.C
	case "D":
		button = input.D
	default:
		return nil
	}
	switch strings.TrimSuffix(name, name[len(name)-1:]) {
	case "236":
		return []input.GameInput{input.Down, input.Down | input.Right, input.Right, button}
	case "214":
		return []input.GameInput{input.Down, input.Down | input.Left, input.Left, button}
	case "623":
		return []input.GameInput{input.Right, input.Down, input.Down | input.Right, button}
	case "423":
		return []input.GameInput{input.Left, input.Down, input.Down | input.Left, button}
	case "22":
		return []input.GameInput{input.Down, input.Down, button}
	case "246":
		return []input.GameInput{input.Down, input.Left, input.Right, button}
	case "642":
		return []input.GameInput{input.Right, input.Left, input.Down, button}
	default:
		return nil
	}
}

func faceRelativeInput(value input.GameInput, facingLeft bool) input.GameInput {
	if !facingLeft {
		return value
	}
	if value&input.Left != 0 {
		value = (value &^ input.Left) | input.Right
	} else if value&input.Right != 0 {
		value = (value &^ input.Right) | input.Left
	}
	return value
}

func (g *GameplayScene) drawDebugGuides(screen *ebiten.Image) {
	if g.camera == nil {
		return
	}

	groundColor := color.RGBA{R: 245, G: 179, B: 0, A: 255}
	wallColor := color.RGBA{R: 50, G: 205, B: 50, A: 255}
	anchorColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	leftWallScreen := g.camera.WorldToScreen(types.Vector2{X: constants.World.X, Y: constants.World.Y})
	rightWallScreen := g.camera.WorldToScreen(types.Vector2{X: constants.World.Right(), Y: constants.World.Y})
	topScreen := g.camera.WorldToScreen(types.Vector2{X: constants.World.X, Y: constants.World.Y})
	bottomScreen := g.camera.WorldToScreen(types.Vector2{X: constants.World.X, Y: constants.World.Bottom()})
	groundStart := g.camera.WorldToScreen(types.Vector2{X: constants.World.X, Y: constants.GroundLevelY})
	groundEnd := g.camera.WorldToScreen(types.Vector2{X: constants.World.Right(), Y: constants.GroundLevelY})

	vector.StrokeLine(screen,
		float32(groundStart.X), float32(groundStart.Y),
		float32(groundEnd.X), float32(groundEnd.Y),
		2,
		groundColor,
		false,
	)

	vector.StrokeLine(screen,
		float32(leftWallScreen.X), float32(topScreen.Y),
		float32(leftWallScreen.X), float32(bottomScreen.Y),
		2,
		wallColor,
		false,
	)

	vector.StrokeLine(screen,
		float32(rightWallScreen.X), float32(topScreen.Y),
		float32(rightWallScreen.X), float32(bottomScreen.Y),
		2,
		wallColor,
		false,
	)

	const crosshairSize = 10.0
	for _, char := range g.gamestate.Characters {
		if char == nil {
			continue
		}

		anchorScreen := g.camera.WorldToScreen(char.Position())

		vector.StrokeLine(screen,
			float32(anchorScreen.X-crosshairSize), float32(anchorScreen.Y),
			float32(anchorScreen.X+crosshairSize), float32(anchorScreen.Y),
			2,
			anchorColor,
			false,
		)

		vector.StrokeLine(screen,
			float32(anchorScreen.X), float32(anchorScreen.Y-crosshairSize),
			float32(anchorScreen.X), float32(anchorScreen.Y+crosshairSize),
			2,
			anchorColor,
			false,
		)
	}
}

func (g *GameplayScene) updateCamera() {
	if g.camera == nil {
		return
	}

	left := g.gamestate.Characters[0]
	if g.gamestate.Characters[0].Position().X < g.gamestate.Characters[1].Position().X {
		left = g.gamestate.Characters[0]
	}
	right := g.gamestate.Characters[1]
	if g.gamestate.Characters[1].Position().X < g.gamestate.Characters[0].Position().X {
		right = g.gamestate.Characters[1]
	}

	if left == nil && right == nil {
		return
	}
	if left == nil {
		g.camera.UpdatePosition(right.Position())
		return
	}
	if right == nil {
		g.camera.UpdatePosition(left.Position())
		return
	}

	midpoint := types.Vector2{
		X: (left.Position().X + right.Position().X) / 2,
		Y: constants.WorldHeight - constants.CameraHeight/2,
	}
	g.camera.UpdatePosition(midpoint)
}

func (g *GameplayScene) updateDebugUI() {
	_, err := g.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Gameplay Debug", image.Rect(0, 0, 320, 340), func(layout debugui.ContainerLayout) {
			ctx.Text("Players")

			for i, char := range g.gamestate.Characters {
				if char == nil || char.StateMachine == nil {
					ctx.Text(fmt.Sprintf("P%d: nil", i+1))
					continue
				}

				sm := char.StateMachine
				animName := "none"
				frameIndex := -1
				frameTimeLeft := 0
				if sm.AnimPlayer != nil {
					animName = sm.AnimPlayer.ActiveAnimationName()
					frameIndex = sm.AnimPlayer.FrameIndex
					frameTimeLeft = sm.AnimPlayer.FrameTimeLeft
				}

				ctx.Text(fmt.Sprintf("P%d pos=(%.2f, %.2f)", i+1, sm.Position.X, sm.Position.Y))
				ctx.Text(fmt.Sprintf("P%d vel=(%.2f, %.2f)", i+1, sm.Velocity.X, sm.Velocity.Y))
				ctx.Text(fmt.Sprintf("P%d facing=%v", i+1, sm.IsFacingLeft))
				ctx.Text(fmt.Sprintf("P%d anim=%s frame=%d t=%d", i+1, animName, frameIndex, frameTimeLeft))
			}

			if g.camera != nil {
				ctx.Text("Camera")
				ctx.Text(fmt.Sprintf("viewport=(x=%.2f y=%.2f w=%.2f h=%.2f)",
					g.camera.Viewport.X,
					g.camera.Viewport.Y,
					g.camera.Viewport.W,
					g.camera.Viewport.H,
				))
				ctx.Text(fmt.Sprintf("scaling=%.2f lockWorld=%v", g.camera.Scaling, g.camera.WorldBoundsLock))
			}
		})
		return nil
	})
	if err != nil {
		fmt.Println("debugui update error:", err)
	}
}
