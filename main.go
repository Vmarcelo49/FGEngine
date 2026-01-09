package main

import (
	"fgengine/character"
	"fgengine/config"
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/input"
	"fgengine/logic"
	"fgengine/player"
	"fgengine/scene"
	"fgengine/stage"
	"fgengine/state"
	"fgengine/types"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	players  []*player.Player
	camera   *graphics.Camera
	stageImg *ebiten.Image

	sceneManager *scene.SceneManager
	renderQueue  *graphics.RenderQueue
	debugui      debugui.DebugUI
}

func (g *Game) updateCamera() {
	if len(g.players) == 0 {
		return
	}

	if len(g.players) == 1 {
		g.camera.UpdatePosition(g.players[0].Character.Position())
		return
	}

	p1 := g.players[0].Character.Position()
	p2 := g.players[1].Character.Position()
	mid := types.Vector2{
		X: (p1.X + p2.X) / 2,
		Y: (p1.Y + p2.Y) / 2,
	}
	g.camera.UpdatePosition(mid)
}

func (g *Game) resolveHits() {
	if len(g.players) < 2 {
		return
	}
	g.applyHit(g.players[0], g.players[1])
	g.applyHit(g.players[1], g.players[0])
}

func (g *Game) applyHit(attacker, defender *player.Player) {
	attackerChar := attacker.Character
	frame := attackerChar.AnimationPlayer.ActiveFrameData()
	if frame == nil {
		return
	}
	if frame.Phase != state.Active {
		return
	}
	if !attackerChar.StateMachine.IsAttacking() || attackerChar.AttackHasHit {
		return
	}

	if !attackerChar.BoundingBox().IsOverlapping(defender.Character.BoundingBox()) {
		return
	}

	damage := frame.Damage
	if damage == 0 {
		damage = 8
	}

	defenderSM := defender.Character.StateMachine
	defenderSM.HP -= damage
	if defenderSM.HP < 0 {
		defenderSM.HP = 0
	}
	hitstun := frame.Hitstun
	if hitstun == 0 {
		hitstun = 30
	}
	defenderSM.HitstunFrames = hitstun

	dir := 1.0
	if attackerChar.Position().X > defender.Character.Position().X {
		dir = -1
	}
	defenderSM.Velocity.X = dir * 4
	defenderSM.Velocity.Y = -2
	defenderSM.ClearState(state.StateAttack | state.StateSpecialAttack | state.StateSuperAttack | state.StateA | state.StateB | state.StateC)
	defenderSM.AddState(state.StateOnHitsun | state.StateAirborne | state.StateFalling)
	defenderSM.ClearState(state.StateGrounded | state.StateNeutral)
	attackerChar.AttackHasHit = true
}

func (g *Game) Update() error {
	if g.sceneManager != nil {
		input.Update(&g.sceneManager.ActiveScene)
	}

	logic.UpdateFacings()

	inputs := make([]input.GameInput, len(g.players))
	for i, p := range g.players {
		inputs[i] = p.Input.Poll()
	}

	logic.UpdateByInputs(inputs, g.players)
	g.resolveHits()
	g.updateCamera()

	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Game Debug Info", image.Rect(0, 0, 320, 200), func(layout debugui.ContainerLayout) {
			ctx.Text("Camera Info:")
			ctx.Text(fmt.Sprintf("Position: (%.2f, %.2f)", g.camera.Viewport.X, g.camera.Viewport.Y))
			ctx.Text(fmt.Sprintf("Size: (%.2f, %.2f)", g.camera.Viewport.W, g.camera.Viewport.H))
			ctx.Text("Character Info:")
			for i, p := range g.players {
				ctx.Text(fmt.Sprintf("Player %d:", i+1))
				ctx.Text(fmt.Sprintf("Position: (%.2f, %.2f)", p.Character.Position().X, p.Character.Position().Y))
				ctx.Text(fmt.Sprintf("State: %s", p.Character.StateMachine.ActiveState.String()))
				ctx.Text(fmt.Sprintf("HP: %d", p.Character.StateMachine.HP))
			}
		})
		return nil
	}); err != nil {
		return err
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

	p1Bindings := map[input.GameInput][]ebiten.Key{
		input.Up:    {ebiten.KeyW},
		input.Down:  {ebiten.KeyS},
		input.Left:  {ebiten.KeyA},
		input.Right: {ebiten.KeyD},
		input.A:     {ebiten.KeyU},
		input.B:     {ebiten.KeyI},
		input.C:     {ebiten.KeyO},
		input.D:     {ebiten.KeyJ},
	}
	p2Bindings := map[input.GameInput][]ebiten.Key{
		input.Up:    {ebiten.KeyUp},
		input.Down:  {ebiten.KeyDown},
		input.Left:  {ebiten.KeyLeft},
		input.Right: {ebiten.KeyRight},
		input.A:     {ebiten.KeyN},
		input.B:     {ebiten.KeyM},
		input.C:     {ebiten.KeyComma},
	}

	p1Input := input.LoadKeyboardBinding(p1Bindings)
	p2Input := input.LoadKeyboardBinding(p2Bindings)

	player1, err := player.NewPlayerWithInput(p1Input)
	if err != nil {
		log.Fatal(err)
	}
	player2, err := player.NewPlayerWithInput(p2Input)
	if err != nil {
		log.Fatal(err)
	}

	centerRect := player1.Character.Sprite().Rect
	centerRect.AlignCenter(constants.World)
	player1.Character.StateMachine.Position = types.Vector2{X: centerRect.X - 120, Y: constants.GroundLevelY - centerRect.H}
	player2Rect := player2.Character.Sprite().Rect
	player2Rect.AlignCenter(constants.World)
	player2.Character.StateMachine.Position = types.Vector2{X: centerRect.X + 120, Y: constants.GroundLevelY - player2Rect.H}

	sm := scene.NewSceneManager()
	game := &Game{
		players:      []*player.Player{player1, player2},
		camera:       graphics.NewCamera(),
		renderQueue:  &graphics.RenderQueue{},
		sceneManager: &sm,
	}
	game.camera.LockWorldBounds = true
	game.updateFacings()

	game.stageImg, _, _ = ebitenutil.NewImageFromFile("assets/stages/PlaceMarkers.png")
	stageDrawable := stage.NewGridStage(32, color.RGBA{40, 40, 40, 255}, color.RGBA{12, 12, 12, 255})
	if game.stageImg != nil {
		stageDrawable = stage.NewImageStage(game.stageImg)
	}

	game.renderQueue.Add(stageDrawable, constants.LayerBG)
	game.renderQueue.Add(player1.Character, constants.LayerPlayer)
	game.renderQueue.Add(player2.Character, constants.LayerPlayer)
	game.renderQueue.Add(&character.BoxDrawable{Character: player1.Character}, constants.LayerEffects)
	game.renderQueue.Add(&character.BoxDrawable{Character: player2.Character}, constants.LayerEffects)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
