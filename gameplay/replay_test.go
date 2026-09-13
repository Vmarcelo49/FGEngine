package gameplay

import (
	"fgengine/character"
	"fgengine/input"
	"testing"
)

// Golden replay hash (SPEC §3.6). Regenerate deliberately: change the
// constant only together with the simulation change that moved it, after
// confirming the new behavior is intended.
// History: F0 genesis; F2 MaxHP joined the snapshot; F4 launch flags
// (KnockdownPending, WallBouncePending, GroundBounceArmed/Used) joined
// the snapshot (no behavior change, script never connects/launches).
const goldenReplayHash uint64 = 8515230550504361799

// replayScript is a fixed 300-frame input stream exercising walks, a jump,
// a normal, a dash, and a special motion on both sides.
func replayScript() [][2]input.GameInput {
	const n = 300
	script := make([][2]input.GameInput, n)
	at := func(frame int, p1, p2 input.GameInput) {
		if frame >= 0 && frame < n {
			script[frame] = [2]input.GameInput{p1, p2}
		}
	}
	for f := 10; f < 70; f++ {
		at(f, input.Right, input.Left)
	}
	at(70, input.Up, input.NoInput)
	at(100, input.A, input.Left)
	// P1 dash: Right, button, Right.
	at(150, input.Right, input.NoInput)
	at(151, input.B, input.NoInput)
	at(152, input.Right, input.NoInput)
	// P2 fireball motion, one step per frame.
	at(200, input.NoInput, input.Down)
	at(201, input.NoInput, input.Down|input.Right)
	at(202, input.NoInput, input.Right)
	at(203, input.NoInput, input.A)
	return script
}

func runReplay(seed uint64) GameState {
	p1 := testStateMachine()
	p1.Position.X = 192
	p2 := testStateMachine()
	p2.Position.X = 576
	g := NewGameState(
		&character.Character{StateMachine: p1},
		&character.Character{StateMachine: p2},
		seed,
	)
	script := replayScript()
	for _, in := range script {
		g.Update(in)
	}
	return g
}

// TestReplayDeterministic is the SPEC §3.6 gate: the same seed + input
// stream must reproduce the exact same end state, matching gold.
func TestReplayDeterministic(t *testing.T) {
	const seed = 20260912
	g1 := runReplay(seed)
	g2 := runReplay(seed)
	h1, h2 := g1.Hash(), g2.Hash()
	t.Logf("replay hash: %d", h1)
	if h1 != h2 {
		t.Fatalf("replay diverged: %d != %d", h1, h2)
	}
	if h1 != goldenReplayHash {
		t.Fatalf("hash %d does not match golden %d", h1, goldenReplayHash)
	}
	if g1.TimerFrames != 5940-len(replayScript()) {
		t.Fatalf("timer = %d, want %d", g1.TimerFrames, 5940-len(replayScript()))
	}
	for i, h := range g1.inputHist {
		if len(h) > 30 {
			t.Fatalf("player %d history exceeded cap: %d", i, len(h))
		}
	}
}
