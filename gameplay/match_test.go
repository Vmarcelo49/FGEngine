package gameplay

import (
	"fgengine/animation"
	"fgengine/constants"
	"fgengine/input"
	"testing"
)

var neutralPair = [2]input.GameInput{input.NoInput, input.NoInput}

// KO (P2 at 10 HP takes 100) ends the round: award, states, freeze.
func TestKOEndsRound(t *testing.T) {
	atk := attackState(100, 100, 0, 0, 0)
	def := poseState(105, 382, "idle")
	def.HP = 10
	g := hitGameState(atk, def)
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.HP; got != 0 {
		t.Fatalf("HP = %d, want clamped 0", got)
	}
	g.Update(neutralPair)
	if g.Wins != [2]int{1, 0} {
		t.Fatalf("wins = %v, want P1 awarded", g.Wins)
	}
	if g.Phase != PhaseRoundEnd {
		t.Fatalf("phase = %v, want RoundEnd", g.Phase)
	}
	if g.FreezeFrames != constants.RoundEndFreezeFrames {
		t.Fatalf("freeze = %d, want %d", g.FreezeFrames, constants.RoundEndFreezeFrames)
	}
	if got := g.Characters[0].StateMachine.AnimPlayer.ActiveAnimationName(); got != "win" {
		t.Fatalf("winner = %q, want win", got)
	}
	if got := g.Characters[1].StateMachine.AnimPlayer.ActiveAnimationName(); got != "ko" {
		t.Fatalf("loser = %q, want ko", got)
	}
}

// Timer expiry awards the higher HP side.
func TestTimerExpiryWin(t *testing.T) {
	g := hitGameState(poseState(300, 382, "idle"), poseState(400, 382, "idle"))
	g.Characters[0].StateMachine.HP = 5000
	g.Characters[1].StateMachine.HP = 3000
	g.TimerFrames = 3
	for i := 0; i < 3; i++ {
		g.Update(neutralPair)
	}
	if g.Wins != [2]int{1, 0} {
		t.Fatalf("wins = %v, want P1 (higher HP)", g.Wins)
	}
	if g.Phase != PhaseRoundEnd {
		t.Fatalf("phase = %v, want RoundEnd", g.Phase)
	}
	if got := g.Characters[0].StateMachine.AnimPlayer.ActiveAnimationName(); got != "win" {
		t.Fatalf("winner = %q, want win", got)
	}
	if got := g.Characters[1].StateMachine.AnimPlayer.ActiveAnimationName(); got != "ko" {
		t.Fatalf("loser = %q, want ko", got)
	}
}

// Timer tie: no award, full reset, round number unchanged.
func TestTimerTieReplays(t *testing.T) {
	g := hitGameState(poseState(300, 382, "idle"), poseState(400, 382, "idle"))
	g.Characters[0].StateMachine.HP = 4000
	g.Characters[1].StateMachine.HP = 4000
	g.Characters[0].StateMachine.Position.X = 300
	g.Characters[1].StateMachine.Position.X = 360
	g.TimerFrames = 2
	for i := 0; i < 2; i++ {
		g.Update(neutralPair)
	}
	if g.Wins != [2]int{0, 0} {
		t.Fatalf("tie must award nothing: %v", g.Wins)
	}
	if g.Phase != PhaseRoundEnd {
		t.Fatalf("phase = %v, want RoundEnd freeze first", g.Phase)
	}
	for i := 0; i < constants.RoundEndFreezeFrames; i++ {
		g.Update(neutralPair)
	}
	assertRoundReset(t, &g, 1)
}

// Double KO (trade, both drop): tie, no award, replay.
func TestDoubleKONoAward(t *testing.T) {
	p1 := attackState(100, 100, 0, 0, 0)
	p2 := attackState(140, 100, 0, 0, 0)
	p2.IsFacingLeft = animation.Left
	g := hitGameState(p1, p2)
	g.Characters[0].StateMachine.HP = 50
	g.Characters[1].StateMachine.HP = 50
	g.ResolveHits()
	if g.Characters[0].StateMachine.HP != 0 || g.Characters[1].StateMachine.HP != 0 {
		t.Fatalf("trade should double-KO: %d / %d",
			g.Characters[0].StateMachine.HP, g.Characters[1].StateMachine.HP)
	}
	g.Update(neutralPair)
	if g.Wins != [2]int{0, 0} {
		t.Fatalf("double KO must award nothing: %v", g.Wins)
	}
	if g.Phase != PhaseRoundEnd {
		t.Fatalf("phase = %v, want RoundEnd", g.Phase)
	}
	for i := 0; i < constants.RoundEndFreezeFrames; i++ {
		g.Update(neutralPair)
	}
	assertRoundReset(t, &g, 1)
}

// Freeze ignores everything except the countdown.
func TestFreezeIgnoresInputs(t *testing.T) {
	atk := attackState(100, 100, 0, 0, 0)
	def := poseState(105, 382, "idle")
	def.HP = 10
	g := hitGameState(atk, def)
	g.ResolveHits()
	g.Update(neutralPair) // detection -> freeze
	px0 := g.Characters[0].StateMachine.Position.X
	px1 := g.Characters[1].StateMachine.Position.X
	attack := [2]input.GameInput{input.A | input.Right, input.A | input.Left}
	for i := 0; i < 30; i++ {
		g.Update(attack)
	}
	if g.Phase != PhaseRoundEnd {
		t.Fatalf("phase = %v, want continued freeze", g.Phase)
	}
	if g.FreezeFrames != constants.RoundEndFreezeFrames-30 {
		t.Fatalf("freeze = %d, want countdown", g.FreezeFrames)
	}
	if g.Characters[0].StateMachine.Position.X != px0 || g.Characters[1].StateMachine.Position.X != px1 {
		t.Fatal("positions moved during freeze")
	}
	if g.Characters[0].StateMachine.HP != 10000 || g.Characters[1].StateMachine.HP != 0 {
		t.Fatal("HP changed during freeze")
	}
}

// Reset contents are exhaustive: a dirty state must come back pristine.
func TestResetContents(t *testing.T) {
	atk := attackState(100, 100, 0, 0, 0)
	def := poseState(105, 382, "idle")
	def.HP = 10
	g := hitGameState(atk, def)
	g.ResolveHits()
	g.Update(neutralPair) // detection -> freeze, Wins[0] = 1
	// Dirty everything the reset must clean.
	d := g.Characters[1].StateMachine
	d.StunFrames = 5
	d.KnockdownPending = true
	d.WallBouncePending = true
	d.Velocity.X = 9
	g.inputHist[0] = append(g.inputHist[0], input.A)
	g.inputHist[1] = append(g.inputHist[1], input.Left)
	g.RecordConnect(0, 1, "A", 0, 0)
	for i := 0; i < constants.RoundEndFreezeFrames; i++ {
		g.Update(neutralPair)
	}
	assertRoundReset(t, &g, 2)
	if g.Wins != [2]int{1, 0} {
		t.Fatalf("wins = %v, want preserved [1 0]", g.Wins)
	}
	if len(g.Connects) != 0 {
		t.Fatalf("ledger survived reset: %d entries", len(g.Connects))
	}
}

func assertRoundReset(t *testing.T, g *GameState, wantRound int) {
	t.Helper()
	if g.Phase != PhaseFight {
		t.Fatalf("phase = %v, want Fight", g.Phase)
	}
	if g.Round != wantRound {
		t.Fatalf("round = %d, want %d", g.Round, wantRound)
	}
	if g.TimerFrames != constants.RoundTimerFrames {
		t.Fatalf("timer = %d, want refilled %d", g.TimerFrames, constants.RoundTimerFrames)
	}
	p1 := g.Characters[0].StateMachine
	p2 := g.Characters[1].StateMachine
	if p1.Position.X != constants.WorldWidth/4 || p2.Position.X != 3*constants.WorldWidth/4 {
		t.Fatalf("positions not reset: %v %v", p1.Position, p2.Position)
	}
	if p1.Position.Y != constants.GroundLevelY || p2.Position.Y != constants.GroundLevelY {
		t.Fatalf("not grounded: %v %v", p1.Position, p2.Position)
	}
	if p1.HP != p1.MaxHP || p2.HP != p2.MaxHP {
		t.Fatalf("HP not refilled: %d/%d", p1.HP, p2.HP)
	}
	if p1.StunFrames != 0 || p2.StunFrames != 0 {
		t.Fatal("stun survived reset")
	}
	if p1.KnockdownPending || p2.KnockdownPending || p1.WallBouncePending || p2.WallBouncePending ||
		p1.GroundBounceArmed || p2.GroundBounceArmed || p1.GroundBounceUsed || p2.GroundBounceUsed {
		t.Fatal("launch flags survived reset")
	}
	if p1.Velocity.X != 0 || p1.Velocity.Y != 0 || p2.Velocity.X != 0 || p2.Velocity.Y != 0 {
		t.Fatal("velocity survived reset")
	}
	if p1.AnimPlayer.ActiveAnimationName() != "idle" || p2.AnimPlayer.ActiveAnimationName() != "idle" {
		t.Fatal("states not reset to idle")
	}
	for i := range g.inputHist {
		if len(g.inputHist[i]) != 0 {
			t.Fatalf("player %d history survived reset", i)
		}
	}
}

// Second round win ends the match.
func TestMatchEnd(t *testing.T) {
	if hitGameState(poseState(0, 0, "idle"), poseState(0, 0, "idle")).MatchOver() {
		t.Fatal("fresh match must not be over")
	}
	atk := attackState(100, 100, 0, 0, 0)
	def := poseState(105, 382, "idle")
	def.HP = 10
	g := hitGameState(atk, def)
	g.Wins = [2]int{1, 0}
	g.ResolveHits()
	g.Update(neutralPair)
	if g.Wins != [2]int{2, 0} {
		t.Fatalf("wins = %v, want [2 0]", g.Wins)
	}
	for i := 0; i < constants.RoundEndFreezeFrames; i++ {
		g.Update(neutralPair)
	}
	if g.Phase != PhaseMatchEnd {
		t.Fatalf("phase = %v, want MatchEnd", g.Phase)
	}
	if !g.MatchOver() {
		t.Fatal("MatchOver must report true")
	}
}

// Winner in win stays there (rule-5 terminal).
func TestWinStaysTerminal(t *testing.T) {
	sm := poseState(100, 382, "win")
	g := hitGameState(idleState(500), sm)
	for i := 0; i < 30; i++ {
		g.Update(neutralPair)
	}
	if got := g.Characters[1].StateMachine.AnimPlayer.ActiveAnimationName(); got != "win" {
		t.Fatalf("win fell through to %q", got)
	}
}
