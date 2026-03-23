package scene

import (
	"fgengine/character"
	"fgengine/constants"
	"fgengine/gameplay"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/stage"
	"fgengine/types"
	"fmt"
	"image"
	"image/color"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func MakeGameplayScene() Scene {
	playerOne, err := character.LoadCharacter("PlaceHolder", 1)
	if err != nil {
		panic(err)
	}

	playerTwo, err := character.LoadCharacter("roa", 2)
	if err != nil {
		panic(err)
	}

	camera := graphics.NewCamera()
	camera.LockWorldBounds = true

	return &GameplayScene{
		camera: camera,
		stage:  stage.NewSolidColorStage(constants.StageColor),
		gamestate: gameplay.GameState{
			Characters: [2]*character.Character{
				playerOne,
				playerTwo,
			}}}
}

type GameplayScene struct {
	camera    *graphics.Camera
	stage     *stage.Stage
	gamestate gameplay.GameState
	debugui   debugui.DebugUI
}

func (g *GameplayScene) Update(inputs [2]input.GameInput) SceneStatus {
	g.gamestate.Update(inputs)
	g.updateCamera()
	g.updateDebugUI()
	return SceneDontChange
}

func (g *GameplayScene) Draw(screen *ebiten.Image) {
	if g.stage != nil {
		g.stage.Draw(screen, g.camera)
	}

	for _, char := range g.gamestate.Characters {
		if char == nil {
			continue
		}
		char.Draw(screen, g.camera)
	}

	g.drawDebugGuides(screen)

	g.debugui.Draw(screen)
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
		Y: constants.WorldHeight / 2,
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
				if sm.ActiveAnim != nil && sm.ActiveAnim.ActiveAnimation != nil {
					animName = sm.ActiveAnim.ActiveAnimation.Name
					frameIndex = sm.ActiveAnim.FrameIndex
					frameTimeLeft = sm.ActiveAnim.FrameTimeLeft
				}

				ctx.Text(fmt.Sprintf("P%d pos=(%.2f, %.2f)", i+1, sm.Position.X, sm.Position.Y))
				ctx.Text(fmt.Sprintf("P%d vel=(%.2f, %.2f)", i+1, sm.Velocity.X, sm.Velocity.Y))
				ctx.Text(fmt.Sprintf("P%d facing=%v", i+1, sm.Facing))
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
				ctx.Text(fmt.Sprintf("scaling=%.2f lockWorld=%v", g.camera.Scaling, g.camera.LockWorldBounds))
			}
		})
		return nil
	})
	if err != nil {
		fmt.Println("debugui update error:", err)
	}
}
