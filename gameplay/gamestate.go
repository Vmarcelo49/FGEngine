package gameplay

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/constants"
	"fgengine/input"
	"slices"
)

type GameState struct {
	Characters [2]*character.Character
	inputHist  [2][]input.GameInput
}

type playerFrameContext struct {
	stateMachine    *animation.StateMachine
	intentAnimation string
	wasAirborne     bool
}

func (g *GameState) Update(inputs [2]input.GameInput) {
	p1 := g.Characters[0].StateMachine
	p2 := g.Characters[1].StateMachine

	g.resolveFacing(p1, p2)

	frame := [2]playerFrameContext{}
	for i, sm := range []*animation.StateMachine{p1, p2} {
		g.pushInputToHistory(i, inputs[i])

		frame[i] = playerFrameContext{
			stateMachine:    sm,
			intentAnimation: input.CheckInputIntent(correctInputByFacing(g.inputHist[i], sm.IsFacingLeft)),
			wasAirborne:     sm.IsAirborne(),
		}
	}

	for _, ctx := range frame {
		// Apply velocity from the framedata
		ctx.stateMachine.ApplyVelocity()

		// Check for animation cancels before applying physics, as some cancels may modify velocity.
		g.checkCancelAnim(ctx)

		// Apply physics (gravity, friction)
		ctx.stateMachine.ApplyPhysics()
	}

	// Resolve player pushbox overlap once after both players have integrated physics.
	CheckHits(p1, p2)
	ResolveBodyCollision(p1, p2)

	for _, ctx := range frame {
		// Checking for landing/falling/idle animations after physics has been applied, as the animation may depend on whether the character is airborne or not.
		g.applyAnimationPostPhysics(ctx)

		// some animations may need info on the input to check some logic
		ctx.stateMachine.AnimPlayer.Update(ctx.intentAnimation)
	}
}

func (g *GameState) resolveFacing(p1, p2 *animation.StateMachine) {
	if p1.Position.X > p2.Position.X {
		if !p1.IsAirborne() {
			p1.IsFacingLeft = animation.Left
		}
		if !p2.IsAirborne() {
			p2.IsFacingLeft = animation.Right
		}
		return
	}

	if !p1.IsAirborne() {
		p1.IsFacingLeft = animation.Right
	}
	if !p2.IsAirborne() {
		p2.IsFacingLeft = animation.Left
	}
}

// pushInputToHistory adds a new input to the player's input history, and ensures the history doesn't exceed the maximum length defined in constants.
func (g *GameState) pushInputToHistory(playerIndex int, in input.GameInput) {
	history := append(g.inputHist[playerIndex], in)
	if len(history) > constants.MaxInputHistory {
		history = history[1:]
	}
	g.inputHist[playerIndex] = history
}

func correctInputByFacing(history []input.GameInput, facing animation.Orientation) []input.GameInput {
	if facing != animation.Left {
		return history
	}

	corrected := make([]input.GameInput, 0, len(history))
	for _, gInput := range history {
		if gInput&input.Left != 0 {
			gInput = (gInput &^ input.Left) | input.Right
		} else if gInput&input.Right != 0 {
			gInput = (gInput &^ input.Right) | input.Left
		}
		corrected = append(corrected, gInput)
	}

	return corrected
}

func (g *GameState) applyAnimationPostPhysics(ctx playerFrameContext) {
	sm := ctx.stateMachine
	if sm == nil || sm.AnimPlayer == nil {
		return
	}

	isAirborne := sm.IsAirborne()
	landedThisFrame := ctx.wasAirborne && !isAirborne

	if landedThisFrame {
		if _, hasLanding := sm.AnimPlayer.Animations["landing"]; hasLanding && sm.AnimPlayer.ActiveAnimationName() != "landing" {
			sm.AnimPlayer.SetAnimation("landing")
		}
	}

	if !sm.AnimPlayer.IsFinished() {
		return
	}

	currentAnim := sm.AnimPlayer.ActiveAnimationName()
	if isAirborne {
		if _, hasFall := sm.AnimPlayer.Animations["fall"]; hasFall && currentAnim != "fall" {
			sm.AnimPlayer.SetAnimation("fall")
		}
		return
	}

	if currentAnim == "landing" || currentAnim == "fall" {
		if ctx.intentAnimation != "" && currentAnim != ctx.intentAnimation {
			sm.AnimPlayer.SetAnimation(ctx.intentAnimation)
		} else if currentAnim != "idle" {
			sm.AnimPlayer.SetAnimation("idle")
		}
		return
	}

	if ctx.intentAnimation != "" {
		if currentAnim != ctx.intentAnimation {
			sm.AnimPlayer.SetAnimation(ctx.intentAnimation)
		}
		return
	}

	if currentAnim == "idle" {
		sm.AnimPlayer.SetAnimation("idle")
		return
	}

	if currentAnim != "idle" {
		sm.AnimPlayer.SetAnimation("idle")
	}
}

func (g *GameState) checkCancelAnim(ctx playerFrameContext) {
	sm := ctx.stateMachine
	if ctx.intentAnimation == "" {
		return
	}

	frameData := sm.AnimPlayer.ActiveFrameData()
	if frameData == nil {
		return
	}

	if !canCancelTo(frameData, sm, ctx.intentAnimation) {
		return
	}

	sm.AnimPlayer.SetAnimation(ctx.intentAnimation)
}

func canCancelTo(frameData *animation.FrameData, sm *animation.StateMachine, intentAnimation string) bool {
	if intentAnimation == "" {
		return false
	}

	if sm.AnimPlayer.ActiveAnimationName() == intentAnimation {
		return false
	}

	if len(frameData.CancelTypes) == 0 {
		return false
	}

	// Prevent jump-start animations while already airborne.
	if (intentAnimation == "7" || intentAnimation == "8" || intentAnimation == "9") && sm.IsAirborne() {
		return false
	}

	if frameData.CancelTypes[0] == "any" {
		return true
	}

	return slices.Contains(frameData.CancelTypes, intentAnimation)
}

/*
char.update()
  func update()
    1. check input / state machine
    2. check physics (gravity, friction, velocity)
    3. check collisions
    4. check animation
*/
