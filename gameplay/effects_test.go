package gameplay

import (
	"fgengine/animation"
	"fgengine/input"
	"math"
	"testing"
)

func TestRecoveryBlocksCancel(t *testing.T) {
	sm := cancelTestSM("A")
	if !canCancelTo(&animation.FrameData{CancelTypes: []string{"any"}}, sm, "B") {
		t.Fatal("control: any should cancel without recovery")
	}
	if canCancelTo(&animation.FrameData{CancelTypes: []string{"any"}, IsRecovery: true}, sm, "B") {
		t.Fatal("recovery must block cancel despite any")
	}
}

func TestHoldableHoldsWhileHeld(t *testing.T) {
	ap := &animation.AnimationPlayer{Animations: map[string]*animation.Animation{
		"walk": {FrameData: []animation.FrameData{{Duration: 3, IsHoldable: true}, {Duration: 3}}},
	}}
	ap.SetAnimation("walk")
	for i := 0; i < 10; i++ {
		ap.Update("walk", 0)
		if ap.FrameIndex != 0 {
			t.Fatalf("held frame advanced: idx=%d", ap.FrameIndex)
		}
	}
	for i := 0; i < 5; i++ {
		ap.Update("", 0)
	}
	if ap.FrameIndex != 1 {
		t.Fatalf("released holdable should advance past, idx=%d", ap.FrameIndex)
	}
}

func TestAnimationSwitch(t *testing.T) {
	ap := &animation.AnimationPlayer{Animations: map[string]*animation.Animation{
		"A": {FrameData: []animation.FrameData{{Duration: 2, AnimationSwitch: "B"}, {Duration: 5}}},
		"B": {FrameData: []animation.FrameData{{Duration: 4}}},
	}}
	ap.SetAnimation("A")
	ap.Update("", 0)
	ap.Update("", 0)
	if got := ap.ActiveAnimationName(); got != "B" {
		t.Fatalf("switch did not fire: active=%q", got)
	}
	if ap.FrameIndex != 0 {
		t.Fatalf("switch target should start at frame 0, idx=%d", ap.FrameIndex)
	}
}

func TestAnimationSwitchOnLastFrame(t *testing.T) {
	ap := &animation.AnimationPlayer{Animations: map[string]*animation.Animation{
		"C": {FrameData: []animation.FrameData{{Duration: 1}, {Duration: 1, AnimationSwitch: "B"}}},
		"B": {FrameData: []animation.FrameData{{Duration: 4}}},
	}}
	ap.SetAnimation("C")
	ap.Update("", 0)
	ap.Update("", 0)
	if got := ap.ActiveAnimationName(); got != "B" {
		t.Fatalf("last-frame switch did not fire: active=%q", got)
	}
}

func waitAnim(t *testing.T, g *GameState, want string, cap int) {
	t.Helper()
	neutral := [2]input.GameInput{input.NoInput, input.NoInput}
	for i := 0; i < cap; i++ {
		if g.Characters[1].StateMachine.AnimPlayer.ActiveAnimationName() == want {
			return
		}
		g.Update(neutral)
	}
	t.Fatalf("never reached %q", want)
}

func TestKnockdownEndToEnd(t *testing.T) {
	atk := attackState(100, 60, 0, 8, 0)
	atk.AnimPlayer.Animations["A"].FrameData[0].Hitstun = 30
	atk.AnimPlayer.Animations["A"].FrameData[0].CanHardKnockdown = true
	def := poseState(105, 382, "idle")
	g := hitGameState(atk, def)
	g.ResolveHits()
	d := g.Characters[1].StateMachine
	if !d.KnockdownPending {
		t.Fatal("knockdown hit must set KnockdownPending")
	}
	if d.Velocity.Y != -8 {
		t.Fatalf("launch vy = %v, want -8", d.Velocity.Y)
	}
	// Attacker retreats so the long active frame can't interfere.
	g.Characters[0].StateMachine.Position.X = 500

	waitAnim(t, &g, "knockdown", 40)
	if d.KnockdownPending {
		t.Fatal("flag must clear on knockdown conversion")
	}
	waitAnim(t, &g, "getup", 40)
	waitAnim(t, &g, "idle", 30)
}

func TestStaleFlagCleared(t *testing.T) {
	def := poseState(105, 382, "idle")
	def.KnockdownPending = true
	def.WallBouncePending = true
	g := hitGameState(idleState(500), def)
	neutral := [2]input.GameInput{input.NoInput, input.NoInput}
	for i := 0; i < 3; i++ {
		g.Update(neutral)
	}
	d := g.Characters[1].StateMachine
	if d.KnockdownPending || d.WallBouncePending {
		t.Fatal("grounded flags without landing must clear as stale")
	}
	if got := d.AnimPlayer.ActiveAnimationName(); got != "idle" {
		t.Fatalf("stale flag caused spurious transition: %q", got)
	}
}

func TestWallBounce(t *testing.T) {
	def := poseState(5, 200, "idle")
	def.Velocity.X = -8
	def.WallBouncePending = true
	g := hitGameState(idleState(500), def)
	g.Update([2]input.GameInput{input.NoInput, input.NoInput})
	d := g.Characters[1].StateMachine
	if d.Position.X != 0 {
		t.Fatalf("x = %v, want clamped 0", d.Position.X)
	}
	if math.Abs(d.Velocity.X-4.8) > 1e-9 {
		t.Fatalf("vx = %v, want reflected +4.8", d.Velocity.X)
	}

	plain := poseState(5, 200, "idle")
	plain.Velocity.X = -8
	g2 := hitGameState(idleState(500), plain)
	g2.Update([2]input.GameInput{input.NoInput, input.NoInput})
	if vx := g2.Characters[1].StateMachine.Velocity.X; vx != 0 {
		t.Fatalf("unflagged wall contact must stop: vx = %v", vx)
	}
}

func TestGroundBounceOnce(t *testing.T) {
	def := poseState(300, 378, "idle")
	def.Velocity.Y = 5
	def.GroundBounceArmed = true
	g := hitGameState(idleState(100), def)
	neutral := [2]input.GameInput{input.NoInput, input.NoInput}
	g.Update(neutral)
	d := g.Characters[1].StateMachine
	if d.Position.Y != 381 {
		t.Fatalf("y = %v, want 1px above ground after bounce", d.Position.Y)
	}
	if math.Abs(d.Velocity.Y+2.4) > 1e-9 {
		t.Fatalf("vy = %v, want reflected -2.4", d.Velocity.Y)
	}
	if !d.GroundBounceUsed {
		t.Fatal("bounce must mark Used")
	}
	for i := 0; i < 30; i++ {
		g.Update(neutral)
	}
	d = g.Characters[1].StateMachine
	if d.Position.Y != 382 || d.Velocity.Y != 0 {
		t.Fatalf("second touchdown must land normally: y=%v vy=%v", d.Position.Y, d.Velocity.Y)
	}
}

func otgAttacker(x float64, otg bool) *animation.StateMachine {
	atk := attackState(x, 50, 2, 5, 0)
	fd := &atk.AnimPlayer.Animations["A"].FrameData[0]
	fd.Hitstun = 10
	fd.CanOTG = otg
	return atk
}

func TestOTGConnectsAndPops(t *testing.T) {
	g := hitGameState(otgAttacker(100, true), poseState(105, 382, "knockdown"))
	g.ResolveHits()
	d := g.Characters[1].StateMachine
	if d.HP != 9950 {
		t.Fatalf("OTG HP = %d, want 9950", d.HP)
	}
	if got := d.AnimPlayer.ActiveAnimationName(); got != "airHurt" {
		t.Fatalf("OTG state = %q, want airHurt pop-up", got)
	}
	if d.Velocity.Y != -5 {
		t.Fatalf("OTG vy = %v, want -5", d.Velocity.Y)
	}
}

func TestNonOTGPassesThroughDowned(t *testing.T) {
	g := hitGameState(otgAttacker(100, false), poseState(105, 382, "knockdown"))
	g.ResolveHits()
	d := g.Characters[1].StateMachine
	if d.HP != 10000 {
		t.Fatalf("non-OTG hit the downed defender: HP = %d", d.HP)
	}
	if got := d.AnimPlayer.ActiveAnimationName(); got != "knockdown" {
		t.Fatalf("state = %q, want undisturbed knockdown", got)
	}
	if len(g.Connects) != 0 {
		t.Fatalf("passthrough recorded %d connects", len(g.Connects))
	}
}

func TestPrioritySuppressesLower(t *testing.T) {
	// P1 wins the trade.
	a := attackState(100, 100, 0, 0, 0)
	b := attackState(140, 100, 0, 0, 0)
	b.IsFacingLeft = animation.Left
	a.AnimPlayer.Animations["A"].FrameData[0].Priority = 5
	b.AnimPlayer.Animations["A"].FrameData[0].Priority = 3
	g := hitGameState(a, b)
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.HP; got != 9900 {
		t.Fatalf("lower-priority defender HP = %d, want 9900", got)
	}
	if got := g.Characters[0].StateMachine.HP; got != 10000 {
		t.Fatalf("higher-priority attacker HP = %d, want untouched 10000", got)
	}

	// Reversed: P2 wins.
	a2 := attackState(100, 100, 0, 0, 0)
	b2 := attackState(140, 100, 0, 0, 0)
	b2.IsFacingLeft = animation.Left
	a2.AnimPlayer.Animations["A"].FrameData[0].Priority = 2
	b2.AnimPlayer.Animations["A"].FrameData[0].Priority = 7
	g2 := hitGameState(a2, b2)
	g2.ResolveHits()
	if got := g2.Characters[0].StateMachine.HP; got != 9900 {
		t.Fatalf("reversed: P1 HP = %d, want 9900", got)
	}
	if got := g2.Characters[1].StateMachine.HP; got != 10000 {
		t.Fatalf("reversed: P2 HP = %d, want untouched 10000", got)
	}
}

func TestGetupBlocksWithoutFlags(t *testing.T) {
	// poseState getup carries no isInvincible: the state rule alone (§7.6)
	// must block.
	g := hitGameState(attackState(100, 100, 0, 0, 0), poseState(105, 382, "getup"))
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.HP; got != 10000 {
		t.Fatalf("getup HP = %d, want untouched", got)
	}
	if len(g.Connects) != 0 {
		t.Fatalf("getup recorded %d connects", len(g.Connects))
	}
}
