package tests

import (
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/player"
	"fgengine/types"
	"fmt"
	"image"
	"reflect"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	players []*player.Player
	camera  *graphics.Camera

	tickCount int
	lastInput input.GameInput
	debugui   debugui.DebugUI
}

func (g *Game) Update() error {
	p1Inputs := g.players[0].Input.GetLocalInputs()

	sm := g.players[0].Character.StateMachine

	sm.InputHistory = append(sm.InputHistory, p1Inputs)
	if len(sm.InputHistory) > constants.MaxInputHistory {
		sm.InputHistory = sm.InputHistory[len(sm.InputHistory)-constants.MaxInputHistory:]
	}
	detectedInputSTR := ""
	for key, seq := range input.InputSequences { // cooldown here probably would be good
		if input.DetectInputSequence(seq, sm.InputHistory) {
			if reflect.DeepEqual(seq, input.InputSequences[key]) {
				detectedInputSTR = key
			}
		}
	}

	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Input Info", image.Rect(0, 0, 512, 288), func(layout debugui.ContainerLayout) {
			ctx.Text(fmt.Sprintf("Current Tick: %2f, Current Fps: %2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
			ctx.Text("Player 1 Inputs:")
			ctx.Text(p1Inputs.String())
			ctx.Text("Input History:")
			if detectedInputSTR != "" {
				ctx.Text("Detected Input: " + detectedInputSTR)

			} else {
				ctx.Text("No special input detected")
			}
			historyLength := len(g.players[0].Character.StateMachine.InputHistory)
			if historyLength > 0 {
				ctx.Loop(historyLength, func(index int) {
					historyInput := g.players[0].Character.StateMachine.InputHistory[index]
					if historyInput == g.lastInput {
						g.tickCount++
					} else {
						g.lastInput = historyInput
						g.tickCount = 1
					}
					ctx.Text(historyInput.String() + fmt.Sprintf(" (%d)", g.tickCount))
				})
			} else {
				ctx.Text("No input history yet")
			}
		})
		return nil
	}); err != nil {
		return err
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.debugui.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func Run() {
	player1 := player.NewDebugPlayer()
	game := &Game{
		players: []*player.Player{player1},
		camera: &graphics.Camera{
			Viewport: types.Rect{
				X: 0, Y: 0,
				W: 1600, H: 900,
			},
		},
	}
	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowTitle("FGEngine - Special Inputs Test")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
