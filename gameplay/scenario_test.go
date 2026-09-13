package gameplay

import (
	"fgengine/character"
	"fgengine/input"
	"testing"
)

// Golden match hash (SPEC §3.6, F7): the full 2-round scripted scenario.
// Same regenerate-discipline as the replay golden.
// History: F7 genesis; 236A became a launcher (knockup, long stun,
// knockdown flag, big hitbox) so the scenario trajectory changed.
const goldenMatchHash uint64 = 7203099118063159584

const (
	scenP1X = 300.0
	scenP2X = 360.0
)

var scenNeutral = [2]input.GameInput{input.NoInput, input.NoInput}

func loadScenarioChars(t *testing.T) (*character.Character, *character.Character) {
	t.Helper()
	p1, err := character.LoadCharacter("PlaceHolder", 1)
	if err != nil {
		t.Fatalf("load P1: %v", err)
	}
	p2, err := character.LoadCharacter("PlaceHolder", 2)
	if err != nil {
		t.Fatalf("load P2: %v", err)
	}
	return p1, p2
}

// anchor pins both fighters' X (Y stays grounded). Deterministic
// drift control between script phases.
func anchor(g *GameState, x1, x2 float64) {
	g.Characters[0].StateMachine.Position.X = x1
	g.Characters[1].StateMachine.Position.X = x2
}

func stepN(g *GameState, in [2]input.GameInput, n int) {
	for i := 0; i < n; i++ {
		g.Update(in)
	}
}

func roundOne(t *testing.T, g *GameState) {
	p1sm := g.Characters[0].StateMachine
	p2sm := g.Characters[1].StateMachine

	// Approach: both walk forward 50 frames from spawn.
	stepN(g, [2]input.GameInput{input.Right, input.Left}, 50)
	gap := p2sm.Position.X - p1sm.Position.X
	if gap >= 384 || gap <= 0 {
		t.Fatalf("approach failed to converge: gap=%v", gap)
	}

	anchor(g, scenP1X, scenP2X)

	// Jump: P1 up, must leave the ground quickly, then land and recover.
	g.Update([2]input.GameInput{input.Up, input.NoInput})
	airborne := false
	for i := 0; i < 10; i++ {
		g.Update(scenNeutral)
		if p1sm.Position.Y < 382 {
			airborne = true
			break
		}
	}
	if !airborne {
		t.Fatal("P1 never left the ground")
	}
	landed := false
	for i := 0; i < 40; i++ {
		g.Update(scenNeutral)
		if p1sm.Position.Y == 382 && p1sm.AnimPlayer.ActiveAnimationName() == "idle" {
			landed = true
			break
		}
	}
	if !landed {
		t.Fatalf("P1 never landed+recovered: y=%v anim=%q",
			p1sm.Position.Y, p1sm.AnimPlayer.ActiveAnimationName())
	}

	// Normals exchange with exact HP accounting (A observed mid-flight).
	want := 10000
	anchor(g, scenP1X, scenP2X)
	g.Update([2]input.GameInput{input.A, input.NoInput})
	stepN(g, scenNeutral, 11)
	if got := p2sm.AnimPlayer.ActiveAnimationName(); got != "hurt" {
		t.Fatalf("mid-hit state = %q, want hurt", got)
	}
	if p2sm.StunFrames <= 0 {
		t.Fatalf("mid-hit stun = %d, want > 0", p2sm.StunFrames)
	}
	want -= 100
	if p2sm.HP != want {
		t.Fatalf("after A: HP=%d, want %d", p2sm.HP, want)
	}
	stepN(g, scenNeutral, 18)
	for _, n := range []struct {
		btn input.GameInput
		dmg int
	}{{input.B, 150}, {input.C, 200}, {input.D, 120}} {
		anchor(g, scenP1X, scenP2X)
		g.Update([2]input.GameInput{n.btn, input.NoInput})
		stepN(g, scenNeutral, 29)
		want -= n.dmg
		if p2sm.HP != want {
			t.Fatalf("after button: HP=%d, want %d", p2sm.HP, want)
		}
	}

	// Special: 236A motion fires and connects.
	anchor(g, scenP1X, scenP2X)
	for _, m := range []input.GameInput{input.Down, input.Down | input.Right, input.Right, input.A} {
		g.Update([2]input.GameInput{m, input.NoInput})
	}
	stepN(g, scenNeutral, 21)
	want -= 300
	if p2sm.HP != want {
		t.Fatalf("after special: HP=%d, want %d", p2sm.HP, want)
	}
	if got := p1sm.AnimPlayer.ActiveAnimationName(); got != "idle" {
		t.Fatalf("P1 after special = %q, want idle", got)
	}

	// Block: P2 holds back through 3 attacks, HP frozen.
	for i := 0; i < 3; i++ {
		anchor(g, scenP1X, scenP2X)
		hold := [2]input.GameInput{input.NoInput, input.Right}
		g.Update([2]input.GameInput{input.A, input.Right})
		stepN(g, hold, 29)
		if p2sm.HP != want {
			t.Fatalf("blocked hit damaged: HP=%d, want %d", p2sm.HP, want)
		}
	}
	anchor(g, scenP1X, scenP2X)
	g.Update([2]input.GameInput{input.A, input.Right})
	stepN(g, [2]input.GameInput{input.NoInput, input.Right}, 7)
	if got := p2sm.AnimPlayer.ActiveAnimationName(); got != "blockHit" {
		t.Fatalf("blocked mid-state = %q, want blockHit", got)
	}
	if p2sm.StunFrames <= 0 {
		t.Fatalf("blocked mid-stun = %d, want > 0", p2sm.StunFrames)
	}
	stepN(g, [2]input.GameInput{input.NoInput, input.Right}, 22)

	// Trade twice: mutual A, both drop 100 each (P1 has taken nothing
	// until now; P2 sits at `want`).
	for i := 0; i < 2; i++ {
		anchor(g, scenP1X, scenP2X+40)
		g.Update([2]input.GameInput{input.A, input.A})
		stepN(g, scenNeutral, 24)
		want -= 100
		if g.Characters[1].StateMachine.HP != want {
			t.Fatalf("after trade %d: P2 HP=%d, want %d", i, g.Characters[1].StateMachine.HP, want)
		}
	}
	if got := g.Characters[0].StateMachine.HP; got != 9800 {
		t.Fatalf("P1 took trade damage: HP=%d, want 9800", got)
	}

	whittleToKO(t, g, [2]int{1, 0})
}

// whittleToKO has P1 press C until P2 drops, re-anchoring every 2 hits,
// then asserts the KO end states and round wins.
func whittleToKO(t *testing.T, g *GameState, wantWins [2]int) {
	t.Helper()
	p1sm := g.Characters[0].StateMachine
	p2sm := g.Characters[1].StateMachine
	cycles := 0
	for p2sm.HP > 0 && cycles < 100 {
		if cycles%2 == 0 {
			anchor(g, scenP1X, scenP2X)
		}
		g.Update([2]input.GameInput{input.C, input.NoInput})
		stepN(g, scenNeutral, 24)
		cycles++
	}
	if p2sm.HP != 0 {
		t.Fatalf("whittle failed to KO in %d cycles (HP=%d)", cycles, p2sm.HP)
	}
	if got := p2sm.AnimPlayer.ActiveAnimationName(); got != "ko" {
		t.Fatalf("loser = %q, want ko", got)
	}
	if got := p1sm.AnimPlayer.ActiveAnimationName(); got != "win" {
		t.Fatalf("winner = %q, want win", got)
	}
	if g.Wins != wantWins {
		t.Fatalf("wins = %v, want %v", g.Wins, wantWins)
	}
}

func TestFullMatchScenario(t *testing.T) {
	t.Chdir(gameplayRepoRoot(t))
	play := func() GameState {
		p1, p2 := loadScenarioChars(t)
		g := NewGameState(p1, p2, 777)
		roundOne(t, &g)
		stepN(&g, scenNeutral, 70)
		if g.Round != 2 {
			t.Fatalf("round = %d, want 2 after freeze", g.Round)
		}
		if x := g.Characters[0].StateMachine.Position.X; x != 192 {
			t.Fatalf("P1 not reset: x=%v", x)
		}
		if hp := g.Characters[1].StateMachine.HP; hp != 10000 {
			t.Fatalf("P2 HP not refilled: %d", hp)
		}
		// Round 2: approach and whittle only.
		stepN(&g, [2]input.GameInput{input.Right, input.Left}, 80)
		anchor(&g, scenP1X, scenP2X)
		whittleToKO(t, &g, [2]int{2, 0})
		stepN(&g, scenNeutral, 70)
		if !g.MatchOver() {
			t.Fatal("match should be over after round 2")
		}
		return g
	}

	g1 := play()
	g2 := play()
	if h1, h2 := g1.Hash(), g2.Hash(); h1 != h2 {
		t.Fatalf("scenario not repeatable: %d != %d", h1, h2)
	}
	if h1 := g1.Hash(); h1 != goldenMatchHash {
		t.Fatalf("hash %d does not match golden %d", h1, goldenMatchHash)
	}
}

// TestLauncherKnockdownCycle proves the 236A launcher interaction chain on
// real data: big hitbox connects → launch with long stun → still stunned
// mid-flight → knockdown on touchdown → getup → idle.
func TestLauncherKnockdownCycle(t *testing.T) {
	t.Chdir(gameplayRepoRoot(t))
	p1, p2 := loadScenarioChars(t)
	p1.StateMachine.Position.X = 300
	p2.StateMachine.Position.X = 360
	g := NewGameState(p1, p2, 4242)
	d := g.Characters[1].StateMachine

	// 236A motion: Down, Down-Forward, Forward, A (P1 faces right).
	for _, m := range []input.GameInput{input.Down, input.Down | input.Right, input.Right, input.A} {
		g.Update([2]input.GameInput{m, input.NoInput})
	}

	launched, downed, risen, done := false, false, false, false
	for i := 0; i < 90 && !done; i++ {
		g.Update(scenNeutral)
		name := d.AnimPlayer.ActiveAnimationName()
		switch {
		case !launched && d.Position.Y < 382 && d.StunFrames > 0:
			launched = true
			if name != "hurt" {
				t.Fatalf("launched in %q, want hurt (grounded at hit)", name)
			}
		case launched && !downed && name == "knockdown":
			downed = true
			if d.StunFrames <= 0 {
				t.Fatal("touchdown must happen while still stunned")
			}
		case downed && !risen && name == "getup":
			risen = true
		case risen && name == "idle":
			done = true
		}
	}
	if !launched {
		t.Fatal("victim never launched stunned")
	}
	if !downed {
		t.Fatal("never converted to knockdown")
	}
	if !risen {
		t.Fatal("never got up")
	}
	if !done {
		t.Fatal("never recovered to idle")
	}
	if d.HP != 9700 {
		t.Fatalf("HP=%d, want exactly the 300 launcher damage", d.HP)
	}
	if d.StunFrames != 0 {
		t.Fatalf("StunFrames=%d, want spent", d.StunFrames)
	}
	if d.KnockdownPending {
		t.Fatal("KnockdownPending must clear on conversion")
	}
}
