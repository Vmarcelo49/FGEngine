package gameplay

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/types"
	"testing"
)

func testStateMachine() *animation.StateMachine {
	ap := &animation.AnimationPlayer{
		Animations: map[string]*animation.Animation{
			"idle": {
				FrameData: []animation.FrameData{
					{Duration: 2},
					{Duration: 2},
				},
			},
		},
	}
	ap.SetAnimation("idle")
	return &animation.StateMachine{
		HP:         10000,
		Position:   types.Vector2{X: 100, Y: 382},
		AnimPlayer: ap,
	}
}

func testGameState(seed uint64) GameState {
	return NewGameState(
		&character.Character{StateMachine: testStateMachine()},
		&character.Character{StateMachine: testStateMachine()},
		seed,
	)
}

// Two identically constructed states must hash identically.
func TestHashDeterministic(t *testing.T) {
	a := testGameState(999)
	b := testGameState(999)
	if a.Hash() != b.Hash() {
		t.Fatalf("identical states hash differently: %d != %d", a.Hash(), b.Hash())
	}
}

// Mutating any snapshotted field must change the hash.
func TestHashSensitive(t *testing.T) {
	base := testGameState(999).Hash()
	mutations := map[string]func(*GameState){
		"position":  func(g *GameState) { g.Characters[0].StateMachine.Position.X++ },
		"velocity":  func(g *GameState) { g.Characters[1].StateMachine.Velocity.Y = 3 },
		"hp":        func(g *GameState) { g.Characters[0].StateMachine.HP-- },
		"stun":      func(g *GameState) { g.Characters[0].StateMachine.StunFrames = 7 },
		"anim":      func(g *GameState) { g.Characters[0].StateMachine.AnimPlayer.FrameIndex = 1 },
		"facing":    func(g *GameState) { g.Characters[0].StateMachine.IsFacingLeft = animation.Left },
		"timer":     func(g *GameState) { g.TimerFrames-- },
		"rng":       func(g *GameState) { g.RNG.Next() },
		"ledger":    func(g *GameState) { g.RecordConnect(0, 1, "idle", 0) },
		"history":   func(g *GameState) { g.inputHist[0] = append(g.inputHist[0], 6) },
		"roundWins": func(g *GameState) { g.Wins[1]++ },
	}
	for name, mutate := range mutations {
		g := testGameState(999)
		mutate(&g)
		if g.Hash() == base {
			t.Errorf("mutating %s did not change the hash", name)
		}
	}
}

// Ledger entries survive while the attacker stays on the recorded frame
// and clear once it advances (SPEC §7.1).
func TestConnectLedgerPrune(t *testing.T) {
	g := testGameState(1)

	if g.HasConnected(0, 1, "idle", 0) {
		t.Fatal("empty ledger must report no connection")
	}
	g.RecordConnect(0, 1, "idle", 0)
	if !g.HasConnected(0, 1, "idle", 0) {
		t.Fatal("recorded connection not found")
	}
	if g.HasConnected(0, 1, "idle", 1) {
		t.Fatal("different frame must not match")
	}
	if g.HasConnected(1, 0, "idle", 0) {
		t.Fatal("swapped players must not match")
	}

	// Same frame: prune keeps the entry.
	g.pruneConnects()
	if !g.HasConnected(0, 1, "idle", 0) {
		t.Fatal("prune dropped an entry for the still-active frame")
	}

	// Attacker advances: prune clears it.
	g.Characters[0].StateMachine.AnimPlayer.FrameIndex = 1
	g.pruneConnects()
	if g.HasConnected(0, 1, "idle", 0) {
		t.Fatal("prune kept an entry for a departed frame")
	}
}
