package gameplay

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/constants"
	"fgengine/input"
	"fgengine/types"
	"slices"
)

type GameState struct {
	Characters [2]*character.Character
	inputHist  [2][]input.GameInput
	// Rollback/determinism substrate (SPEC §3.4, §3.5, §7.1, §7.7).
	RNG          SplitMix64
	Connects     []ConnectKey
	TimerFrames  int
	Round        int
	Wins         [2]int
	Phase        MatchPhase
	FreezeFrames int
}

// NewGameState builds a match-ready GameState. The seed fixes all future
// simulation randomness (SPEC §3.4): versus play passes a time-based seed
// at match creation (outside the sim path), tests pass fixed seeds.
func NewGameState(p1, p2 *character.Character, seed uint64) GameState {
	return GameState{
		Characters:  [2]*character.Character{p1, p2},
		RNG:         SplitMix64{State: seed},
		TimerFrames: constants.RoundTimerFrames,
		Round:       1,
		Phase:       PhaseFight,
	}
}

type playerFrameContext struct {
	stateMachine    *animation.StateMachine
	intentAnimation string
	wasAirborne     bool
}

func (g *GameState) Update(inputs [2]input.GameInput) {
	if g.Phase != PhaseFight {
		// Round over: full freeze, countdown only (SPEC §7.7).
		g.updateMatchFlow()
		return
	}

	g.pruneConnects()

	// Round timer ticks during the fight phase.
	if g.TimerFrames > 0 {
		g.TimerFrames--
	}

	p1 := g.Characters[0].StateMachine
	p2 := g.Characters[1].StateMachine

	g.resolveFacing(p1, p2)

	frame := [2]playerFrameContext{}
	for i, sm := range []*animation.StateMachine{p1, p2} {
		// Stun countdowns tick before anything else can set them this
		// frame (SPEC §4.1 step 2, §7.6).
		if sm.StunFrames > 0 {
			sm.StunFrames--
		}
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

	// Resolve hits (F1 subset of §7.2), then player pushbox overlap, once
	// after both players have integrated physics.
	g.ResolveHits()
	ResolveBodyCollision(p1, p2)

	for _, ctx := range frame {
		// Checking for landing/falling/idle animations after physics has been applied, as the animation may depend on whether the character is airborne or not.
		g.applyAnimationPostPhysics(ctx)

		// some animations may need info on the input to check some logic
		ctx.stateMachine.AnimPlayer.Update(ctx.intentAnimation, ctx.stateMachine.StunFrames)
	}

	// Match flow last: round-end detection (SPEC §4.1 step 7, §7.7).
	g.updateMatchFlow()
}

// updateMatchFlow runs round-end detection, the freeze countdown, and
// round reset / match end (SPEC §7.7).
func (g *GameState) updateMatchFlow() {
	switch g.Phase {
	case PhaseMatchEnd:
		return
	case PhaseRoundEnd:
		if g.FreezeFrames > 0 {
			g.FreezeFrames--
		}
		if g.FreezeFrames > 0 {
			return
		}
		if g.Wins[0] >= constants.RoundsToWin || g.Wins[1] >= constants.RoundsToWin {
			g.Phase = PhaseMatchEnd
			return
		}
		g.resetRound()
		return
	default: // PhaseFight
		p1Dead := g.Characters[0].StateMachine.HP <= 0
		p2Dead := g.Characters[1].StateMachine.HP <= 0
		switch {
		case p1Dead && p2Dead:
			// Double KO: tie, no award (SPEC §7.7).
			g.beginRoundEnd(-1)
		case p1Dead:
			g.beginRoundEnd(1)
		case p2Dead:
			g.beginRoundEnd(0)
		case g.TimerFrames <= 0:
			hp0 := g.Characters[0].StateMachine.HP
			hp1 := g.Characters[1].StateMachine.HP
			switch {
			case hp0 > hp1:
				g.beginRoundEnd(0)
			case hp1 > hp0:
				g.beginRoundEnd(1)
			default:
				g.beginRoundEnd(-1)
			}
		}
	}
}

// beginRoundEnd awards the round (unless tie), forces end states, and
// starts the freeze. Winner -1 = tie/double-KO (SPEC §7.7).
func (g *GameState) beginRoundEnd(winner int) {
	if winner >= 0 {
		g.Wins[winner]++
		// Awarded rounds advance the round number; ties replay it.
		g.Round++
		loser := 1 - winner
		g.Characters[loser].StateMachine.AnimPlayer.SetAnimation("ko")
		g.Characters[loser].StateMachine.StunFrames = 0
		g.Characters[winner].StateMachine.AnimPlayer.SetAnimation("win")
		g.Characters[winner].StateMachine.StunFrames = 0
	}
	g.Phase = PhaseRoundEnd
	g.FreezeFrames = constants.RoundEndFreezeFrames
}

// resetRound starts the next round: positions, HP, timer, and states reset;
// histories and ledger clear so nothing carries over (SPEC §7.7).
func (g *GameState) resetRound() {
	starts := [2]float64{constants.WorldWidth / 4, 3 * constants.WorldWidth / 4}
	for i, sm := range []*animation.StateMachine{g.Characters[0].StateMachine, g.Characters[1].StateMachine} {
		sm.Position = types.Vector2{X: starts[i], Y: constants.GroundLevelY}
		sm.Velocity = types.Vector2{}
		sm.HP = sm.MaxHP
		sm.StunFrames = 0
		sm.KnockdownPending = false
		sm.WallBouncePending = false
		sm.GroundBounceArmed = false
		sm.GroundBounceUsed = false
		sm.IgnoreGravityFrames = 0
		sm.IsFacingLeft = animation.Right
		if i == 1 {
			sm.IsFacingLeft = animation.Left
		}
		sm.AnimPlayer.SetAnimation("idle")
		g.inputHist[i] = nil
	}
	g.Connects = nil
	g.TimerFrames = constants.RoundTimerFrames
	g.Phase = PhaseFight
	g.FreezeFrames = 0
}

// MatchOver reports a decided match (SPEC §7.7). The scene exits on it.
func (g GameState) MatchOver() bool {
	return g.Phase == PhaseMatchEnd
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

	// Explicit exits (SPEC §6.7 rule 5): these states never fall through
	// to rule 4.
	if sm.AnimPlayer.ActiveAnimation != nil {
		switch sm.AnimPlayer.ActiveAnimation.Name {
		case "ko", "win":
			return
		case "knockdown":
			if sm.AnimPlayer.IsFinished() {
				sm.AnimPlayer.SetAnimation("getup")
			}
			return
		case "getup":
			if sm.AnimPlayer.IsFinished() {
				sm.AnimPlayer.SetAnimation("idle")
			}
			return
		}
	}

	isAirborne := sm.IsAirborne()
	landedThisFrame := ctx.wasAirborne && !isAirborne

	if landedThisFrame {
		// Ground contact consumes launch memory (§7.5).
		knockdown := sm.KnockdownPending
		sm.KnockdownPending = false
		sm.WallBouncePending = false
		sm.GroundBounceArmed = false
		sm.GroundBounceUsed = false
		current := sm.AnimPlayer.ActiveAnimationName()
		switch {
		case current == "airBlock" && sm.StunFrames > 0:
			// Landing converts remaining air blockstun to blockHit,
			// counter preserved (§7.4).
			sm.AnimPlayer.SetAnimation("blockHit")
		case knockdown:
			// Knockdown launch touches down → knockdown state (§7.5).
			sm.AnimPlayer.SetAnimation("knockdown")
		default:
			if _, hasLanding := sm.AnimPlayer.Animations["landing"]; hasLanding && current != "landing" {
				sm.AnimPlayer.SetAnimation("landing")
			}
		}
	} else if !isAirborne && sm.Velocity.Y >= 0 {
		// Grounded without a landing event and not launching this frame:
		// any launch flags are stale. (A hit victim keeps its flags here
		// because its just-set upward velocity reads vy < 0.)
		sm.KnockdownPending = false
		sm.WallBouncePending = false
		sm.GroundBounceArmed = false
		sm.GroundBounceUsed = false
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
	// Recovery frames never cancel, even with cancelTypes listed (SPEC §6.5).
	if frameData.IsRecovery {
		return false
	}

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

	// "any" anywhere in the list matches everything (SPEC §6.5).
	if slices.Contains(frameData.CancelTypes, "any") {
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
