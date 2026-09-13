package gameplay

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/input"
	"testing"
)

// guardGame builds P1-attacker vs P2-defender with P2 facing left and the
// given physical input as the history tail (Right = holding back for P2).
func guardGame(atk *animation.StateMachine, def *animation.StateMachine, hold input.GameInput) GameState {
	def.IsFacingLeft = animation.Left
	g := hitGameState(atk, def)
	g.inputHist[1] = []input.GameInput{hold}
	return g
}

func blockAttack() *animation.StateMachine {
	atk := attackState(100, 100, 5, 0, 10)
	atk.AnimPlayer.Animations["A"].FrameData[0].Hitstun = 12
	atk.AnimPlayer.Animations["A"].FrameData[0].Blockstun = 8
	return atk
}

func TestGroundBlockHolds(t *testing.T) {
	g := guardGame(blockAttack(), poseState(105, 382, "idle"), input.Right)
	g.ResolveHits()
	def := g.Characters[1].StateMachine
	if def.HP != 10000 {
		t.Fatalf("blocked HP = %d, want no damage (chip 0)", def.HP)
	}
	if got := def.AnimPlayer.ActiveAnimationName(); got != "blockHit" {
		t.Fatalf("state = %q, want blockHit", got)
	}
	if def.StunFrames != 8 {
		t.Fatalf("StunFrames = %d, want attacker blockstun 8", def.StunFrames)
	}
	if def.Velocity.X != 0 || def.Velocity.Y != 0 {
		t.Fatalf("blocked defender took velocity: %v", def.Velocity)
	}
	if def.Position.X != 115 || g.Characters[0].StateMachine.Position.X != 95 {
		t.Fatalf("pushback split wrong: def=%v atk=%v",
			def.Position, g.Characters[0].StateMachine.Position)
	}
}

func TestCrouchBlock(t *testing.T) {
	g := guardGame(blockAttack(), poseState(105, 382, "2"), input.Right)
	g.ResolveHits()
	def := g.Characters[1].StateMachine
	if got := def.AnimPlayer.ActiveAnimationName(); got != "crouchBlock" {
		t.Fatalf("state = %q, want crouchBlock", got)
	}
	if def.HP != 10000 || def.StunFrames != 8 {
		t.Fatalf("crouch block wrong: HP=%d stun=%d", def.HP, def.StunFrames)
	}
}

func TestAirBlock(t *testing.T) {
	g := guardGame(blockAttack(), poseState(105, 360, "fall"), input.Right)
	g.ResolveHits()
	def := g.Characters[1].StateMachine
	if got := def.AnimPlayer.ActiveAnimationName(); got != "airBlock" {
		t.Fatalf("state = %q, want airBlock", got)
	}
	if def.HP != 10000 || def.StunFrames != 8 {
		t.Fatalf("air block wrong: HP=%d stun=%d", def.HP, def.StunFrames)
	}
	if def.Position.Y != 360 {
		t.Fatalf("direct resolve must not move vertically: y=%v", def.Position.Y)
	}
}

func TestNoGuardWithoutHold(t *testing.T) {
	g := guardGame(blockAttack(), poseState(105, 382, "idle"), input.NoInput)
	g.ResolveHits()
	def := g.Characters[1].StateMachine
	if got := def.AnimPlayer.ActiveAnimationName(); got != "hurt" {
		t.Fatalf("state = %q, want hurt (no guard without back held)", got)
	}
	if def.HP != 9900 {
		t.Fatalf("HP = %d, want 9900", def.HP)
	}
}

func TestStunDoesNotGuard(t *testing.T) {
	def := poseState(105, 382, "hurt")
	def.StunFrames = 5
	g := guardGame(blockAttack(), def, input.Right)
	g.ResolveHits()
	def = g.Characters[1].StateMachine
	if def.HP != 9900 {
		t.Fatalf("hitstun must not guard: HP = %d", def.HP)
	}
	if got := def.AnimPlayer.ActiveAnimationName(); got != "hurt" {
		t.Fatalf("state = %q, want hurt re-entry", got)
	}
}

func TestAttackDoesNotGuard(t *testing.T) {
	def := attackState(105, 0, 0, 0, 0)
	g := guardGame(blockAttack(), def, input.Right)
	g.ResolveHits()
	// Defender is Characters[1]; its zero-damage attack also trades back,
	// which is fine — the point is it got hit (no guard in attack state).
	if got := g.Characters[1].StateMachine.HP; got != 9900 {
		t.Fatalf("attack state must not guard: HP = %d", got)
	}
}

func TestReguardRefreshes(t *testing.T) {
	def := poseState(105, 382, "blockHit")
	def.StunFrames = 3
	g := guardGame(blockAttack(), def, input.Right)
	g.ResolveHits()
	def = g.Characters[1].StateMachine
	if got := def.AnimPlayer.ActiveAnimationName(); got != "blockHit" {
		t.Fatalf("state = %q, want continued blockHit", got)
	}
	if def.StunFrames != 8 {
		t.Fatalf("StunFrames = %d, want refreshed 8", def.StunFrames)
	}
	if def.HP != 10000 {
		t.Fatalf("re-guard took damage: HP = %d", def.HP)
	}
}

func TestLandingConvertsAirBlock(t *testing.T) {
	atk := idleState(500)
	def := poseState(300, 370, "airBlock")
	def.StunFrames = 20
	g := hitGameState(atk, def)
	neutral := [2]input.GameInput{input.NoInput, input.NoInput}
	for i := 0; i < 10; i++ {
		g.Update(neutral)
	}
	def = g.Characters[1].StateMachine
	if got := def.AnimPlayer.ActiveAnimationName(); got != "blockHit" {
		t.Fatalf("landed airBlock became %q, want blockHit", got)
	}
	if def.StunFrames <= 0 {
		t.Fatalf("counter lost on conversion: StunFrames = %d", def.StunFrames)
	}
	if def.Position.Y != 382 {
		t.Fatalf("not grounded after 10 frames: y=%v", def.Position.Y)
	}
}

func TestBlockstunExpiryMidAirFalls(t *testing.T) {
	atk := idleState(500)
	def := poseState(300, 200, "airBlock")
	def.StunFrames = 4
	g := hitGameState(atk, def)
	neutral := [2]input.GameInput{input.NoInput, input.NoInput}
	for i := 0; i < 15; i++ {
		g.Update(neutral)
	}
	def = g.Characters[1].StateMachine
	if got := def.AnimPlayer.ActiveAnimationName(); got != "fall" {
		t.Fatalf("expired mid-air block became %q, want fall", got)
	}
	if def.Position.Y >= 382 {
		t.Fatalf("should still be airborne: y=%v", def.Position.Y)
	}
}

func guardCycle(t *testing.T, holdBack bool, frames int) GameState {
	t.Helper()
	t.Chdir(gameplayRepoRoot(t))
	p1, err := character.LoadCharacter("PlaceHolder", 1)
	if err != nil {
		t.Fatalf("load P1: %v", err)
	}
	p2, err := character.LoadCharacter("PlaceHolder", 2)
	if err != nil {
		t.Fatalf("load P2: %v", err)
	}
	p1.StateMachine.Position.X = 300
	p2.StateMachine.Position.X = 360
	g := NewGameState(p1, p2, 42)
	script := make([][2]input.GameInput, frames)
	for f := range script {
		// P2 holds physical back (Right, facing left) while the flag says so.
		if holdBack && f <= 20 {
			script[f][1] = input.Right
		}
	}
	if frames > 5 {
		script[5][0] = input.A
	}
	for _, in := range script {
		g.Update(in)
	}
	return g
}

// TestGuardAcceptance is the F3 acceptance on the real placeholder:
// holding back turns the scripted A-attack into blockstun with no damage;
// without the hold it hits (contrast).
func TestGuardAcceptance(t *testing.T) {
	mid := guardCycle(t, true, 12)
	p2 := mid.Characters[1].StateMachine
	if got := p2.AnimPlayer.ActiveAnimationName(); got != "blockHit" {
		t.Fatalf("guarded mid-cycle state = %q, want blockHit", got)
	}
	if p2.HP != 10000 {
		t.Fatalf("guarded mid-cycle HP = %d, want no damage", p2.HP)
	}
	if p2.StunFrames <= 0 {
		t.Fatalf("guarded mid-cycle StunFrames = %d, want > 0", p2.StunFrames)
	}

	end := guardCycle(t, true, 40)
	p2 = end.Characters[1].StateMachine
	if got := p2.AnimPlayer.ActiveAnimationName(); got != "idle" {
		t.Fatalf("guarded end-cycle state = %q, want idle", got)
	}
	if p2.HP != 10000 || p2.StunFrames != 0 {
		t.Fatalf("guarded end-cycle wrong: HP=%d stun=%d", p2.HP, p2.StunFrames)
	}

	noGuard := guardCycle(t, false, 12)
	p2 = noGuard.Characters[1].StateMachine
	if got := p2.AnimPlayer.ActiveAnimationName(); got != "hurt" {
		t.Fatalf("unguarded state = %q, want hurt", got)
	}
	if p2.HP != 9900 {
		t.Fatalf("unguarded HP = %d, want 9900", p2.HP)
	}
}
