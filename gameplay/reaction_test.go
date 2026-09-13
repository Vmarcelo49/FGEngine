package gameplay

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/input"
	"fgengine/types"
	"path/filepath"
	"runtime"
	"testing"
)

// poseState builds a defender with the reaction set available. The animName
// argument sets the starting animation (idle crouch, ko...).
func poseState(x, y float64, animName string) *animation.StateMachine {
	hurtBox := map[types.BoxType][]types.Rect{
		types.Hurt: {{X: -32, Y: -40, W: 64, H: 40}},
	}
	lyingBox := map[types.BoxType][]types.Rect{
		types.Hurt: {{X: -40, Y: -24, W: 80, H: 24}},
	}
	idleFD := []animation.FrameData{{Duration: 100, Boxes: hurtBox}}
	// Reaction frames carry hurtboxes: stunned defenders can be re-hit
	// (hitstun chaining) and re-blocked (re-guard).
	hold2 := []animation.FrameData{
		{Duration: 3, Boxes: hurtBox},
		{Duration: 3, Boxes: hurtBox},
	}
	loop01 := &animation.LoopFrame{Start: 0, End: 1}
	ap := &animation.AnimationPlayer{Animations: map[string]*animation.Animation{
		"idle":        {FrameData: idleFD},
		"2":           {FrameData: idleFD},
		"fall":        {FrameData: idleFD},
		"hurt":        {FrameData: hold2, LoopFrames: loop01},
		"crouchHurt":  {FrameData: hold2, LoopFrames: loop01},
		"airHurt":     {FrameData: hold2, LoopFrames: loop01},
		"blockHit":    {FrameData: hold2, LoopFrames: loop01},
		"crouchBlock": {FrameData: hold2, LoopFrames: loop01},
		"airBlock":    {FrameData: hold2, LoopFrames: loop01},
		"knockdown":   {FrameData: []animation.FrameData{{Duration: 6, Boxes: lyingBox}, {Duration: 6, Boxes: lyingBox}, {Duration: 6, Boxes: lyingBox}}},
		"getup":       {FrameData: []animation.FrameData{{Duration: 4}, {Duration: 4}}},
		"ko":          {FrameData: []animation.FrameData{{Duration: 8}, {Duration: 8}}},
		"win":         {FrameData: []animation.FrameData{{Duration: 8}, {Duration: 8}}},
	}}
	ap.SetAnimation(animName)
	return &animation.StateMachine{
		HP:         10000,
		MaxHP:      10000,
		Position:   types.Vector2{X: x, Y: y},
		AnimPlayer: ap,
	}
}

func TestHitEntersHurtWithStun(t *testing.T) {
	atk := attackState(100, 100, 2, 0, 0)
	atk.AnimPlayer.Animations["A"].FrameData[0].Hitstun = 12
	g := hitGameState(atk, poseState(105, 382, "idle"))
	g.ResolveHits()
	def := g.Characters[1].StateMachine
	if got := def.AnimPlayer.ActiveAnimationName(); got != "hurt" {
		t.Fatalf("state = %q, want hurt", got)
	}
	if def.StunFrames != 12 {
		t.Fatalf("StunFrames = %d, want 12", def.StunFrames)
	}
	if def.Velocity.X != 2 {
		t.Fatalf("vx = %v, want knockback 2", def.Velocity.X)
	}
}

func TestHitSelectsCrouchHurt(t *testing.T) {
	atk := attackState(100, 100, 0, 0, 0)
	atk.AnimPlayer.Animations["A"].FrameData[0].Hitstun = 12
	g := hitGameState(atk, poseState(105, 382, "2"))
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.AnimPlayer.ActiveAnimationName(); got != "crouchHurt" {
		t.Fatalf("state = %q, want crouchHurt", got)
	}
}

func TestHitSelectsAirHurt(t *testing.T) {
	atk := attackState(100, 100, 0, 0, 0)
	atk.AnimPlayer.Animations["A"].FrameData[0].Hitstun = 12
	// Airborne (y < ground) but inside the hit's vertical span.
	g := hitGameState(atk, poseState(105, 360, "idle"))
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.AnimPlayer.ActiveAnimationName(); got != "airHurt" {
		t.Fatalf("state = %q, want airHurt", got)
	}
}

func TestInvincibilityBlocks(t *testing.T) {
	def := poseState(105, 382, "idle")
	def.AnimPlayer.Animations["idle"].FrameData[0].IsInvincible = true
	g := hitGameState(attackState(100, 100, 5, 5, 5), def)
	g.ResolveHits()
	def = g.Characters[1].StateMachine
	if def.HP != 10000 {
		t.Fatalf("invincible defender took damage: HP = %d", def.HP)
	}
	if got := def.AnimPlayer.ActiveAnimationName(); got != "idle" {
		t.Fatalf("invincible defender left idle: %q", got)
	}
	if len(g.Connects) != 0 {
		t.Fatalf("invincible defender recorded %d connects", len(g.Connects))
	}
}

func TestArmorAbsorbs(t *testing.T) {
	def := poseState(105, 382, "idle")
	def.AnimPlayer.Animations["idle"].FrameData[0].HasArmor = true
	g := hitGameState(attackState(100, 100, 5, 0, 10), def)
	g.ResolveHits()
	def = g.Characters[1].StateMachine
	if def.HP != 9900 {
		t.Fatalf("armored defender HP = %d, want 9900 (damage applies)", def.HP)
	}
	if got := def.AnimPlayer.ActiveAnimationName(); got != "idle" {
		t.Fatalf("armored defender reacted: %q", got)
	}
	if def.Velocity.X != 0 || def.Position.X != 105 || g.Characters[0].StateMachine.Position.X != 100 {
		t.Fatalf("armored defender moved or attacker pushed: def=%v atk=%v",
			def.Position, g.Characters[0].StateMachine.Position)
	}
	if len(g.Connects) != 1 {
		t.Fatalf("armored hit must still record the connect, got %d", len(g.Connects))
	}
}

func TestKOHPlumbs(t *testing.T) {
	def := poseState(105, 382, "idle")
	def.HP = 50
	g := hitGameState(attackState(100, 100, 5, 0, 0), def)
	g.ResolveHits()
	def = g.Characters[1].StateMachine
	if def.HP != 0 {
		t.Fatalf("HP = %d, want clamped 0 (not -50)", def.HP)
	}
	if got := def.AnimPlayer.ActiveAnimationName(); got != "ko" {
		t.Fatalf("state = %q, want ko", got)
	}
	if def.StunFrames != 0 {
		t.Fatalf("StunFrames = %d on KO, want 0", def.StunFrames)
	}
}

func TestKoTakesNone(t *testing.T) {
	def := poseState(105, 382, "ko")
	def.HP = 5000
	g := hitGameState(attackState(100, 100, 5, 0, 0), def)
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.HP; got != 5000 {
		t.Fatalf("ko defender HP changed: %d", got)
	}
	if len(g.Connects) != 0 {
		t.Fatalf("ko defender recorded %d connects", len(g.Connects))
	}
}

func TestKoStaysThroughUpdate(t *testing.T) {
	atk := poseState(100, 382, "idle")
	def := poseState(105, 382, "ko")
	g := hitGameState(atk, def)
	neutral := [2]input.GameInput{input.NoInput, input.NoInput}
	for i := 0; i < 30; i++ {
		g.Update(neutral)
	}
	if got := g.Characters[1].StateMachine.AnimPlayer.ActiveAnimationName(); got != "ko" {
		t.Fatalf("ko fell through to %q (rule 5 violation)", got)
	}
}

func gameplayRepoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "..")
}

// TestCombatCycleStableAndRepeatable is the F2 acceptance: a scripted
// A-attack on the real placeholder runs hit → HP drop → hurt → idle, and
// re-running the script reproduces the exact state.
func TestCombatCycleStableAndRepeatable(t *testing.T) {
	runCycle := func() GameState {
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
		script := make([][2]input.GameInput, 40)
		script[5] = [2]input.GameInput{input.A, input.NoInput}
		for _, in := range script {
			g.Update(in)
		}
		return g
	}

	// Mid-cycle: hit landed, defender in hurt with stun remaining.
	mid := runCyclePartial(t, 12)
	if got := mid.Characters[1].StateMachine.AnimPlayer.ActiveAnimationName(); got != "hurt" {
		t.Fatalf("mid-cycle state = %q, want hurt", got)
	}
	if stun := mid.Characters[1].StateMachine.StunFrames; stun <= 0 {
		t.Fatalf("mid-cycle StunFrames = %d, want > 0", stun)
	}
	if hp := mid.Characters[1].StateMachine.HP; hp != 9900 {
		t.Fatalf("mid-cycle HP = %d, want 9900", hp)
	}

	// End of cycle: back to idle, stun spent.
	end1 := runCycle()
	end2 := runCycle()
	p2 := end1.Characters[1].StateMachine
	if got := p2.AnimPlayer.ActiveAnimationName(); got != "idle" {
		t.Fatalf("end-cycle state = %q, want idle", got)
	}
	if p2.StunFrames != 0 {
		t.Fatalf("end-cycle StunFrames = %d, want 0", p2.StunFrames)
	}
	if p2.HP != 9900 {
		t.Fatalf("end-cycle HP = %d, want exactly one hit (9900)", p2.HP)
	}
	if h1, h2 := end1.Hash(), end2.Hash(); h1 != h2 {
		t.Fatalf("cycle not repeatable: %d != %d", h1, h2)
	}
}

func runCyclePartial(t *testing.T, frames int) GameState {
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
	if frames > 5 {
		script[5] = [2]input.GameInput{input.A, input.NoInput}
	}
	for _, in := range script {
		g.Update(in)
	}
	return g
}
