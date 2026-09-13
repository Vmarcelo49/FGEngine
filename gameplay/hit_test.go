package gameplay

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/types"
	"testing"
)

func attackState(x float64, damage, kb, kup, pushback int) *animation.StateMachine {
	ap := &animation.AnimationPlayer{
		Animations: map[string]*animation.Animation{
			"A": {FrameData: []animation.FrameData{{
				Duration:  100,
				Damage:    damage,
				Knockback: kb,
				Knockup:   kup,
				Pushback:  pushback,
				Boxes: map[types.BoxType][]types.Rect{
					types.Hit:  {{X: 0, Y: -40, W: 60, H: 40}},
					types.Hurt: {{X: -32, Y: -40, W: 64, H: 40}},
				},
			}}},
		},
	}
	ap.SetAnimation("A")
	return &animation.StateMachine{
		HP:         10000,
		MaxHP:      10000,
		Position:   types.Vector2{X: x, Y: 382},
		AnimPlayer: ap,
	}
}

func idleState(x float64) *animation.StateMachine {
	ap := &animation.AnimationPlayer{
		Animations: map[string]*animation.Animation{
			"idle": {FrameData: []animation.FrameData{{
				Duration: 100,
				Boxes: map[types.BoxType][]types.Rect{
					types.Hurt: {{X: -32, Y: -40, W: 64, H: 40}},
				},
			}}},
		},
	}
	ap.SetAnimation("idle")
	return &animation.StateMachine{
		HP:         10000,
		MaxHP:      10000,
		Position:   types.Vector2{X: x, Y: 382},
		AnimPlayer: ap,
	}
}

func hitGameState(p1, p2 *animation.StateMachine) GameState {
	return NewGameState(
		&character.Character{StateMachine: p1},
		&character.Character{StateMachine: p2},
		7,
	)
}

func TestResolveAppliesDamage(t *testing.T) {
	g := hitGameState(attackState(100, 100, 0, 0, 0), idleState(105))
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.HP; got != 9900 {
		t.Fatalf("defender HP = %d, want 9900", got)
	}
	if got := g.Characters[0].StateMachine.HP; got != 10000 {
		t.Fatalf("attacker HP changed: %d", got)
	}
}

func TestResolveKnockbackDirection(t *testing.T) {
	// Attacker left of defender: knocked right.
	g := hitGameState(attackState(100, 0, 3, 0, 0), idleState(105))
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.Velocity.X; got != 3 {
		t.Fatalf("knockback vx = %v, want +3", got)
	}
	// Mirrored: attacker right of defender, facing it, knocks left.
	mirrored := attackState(180, 0, 3, 0, 0)
	mirrored.IsFacingLeft = animation.Left
	g = hitGameState(idleState(100), mirrored)
	g.ResolveHits()
	if got := g.Characters[0].StateMachine.Velocity.X; got != -3 {
		t.Fatalf("mirrored knockback vx = %v, want -3", got)
	}
}

func TestResolveKnockup(t *testing.T) {
	g := hitGameState(attackState(100, 0, 0, 4, 0), idleState(105))
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.Velocity.Y; got != -4 {
		t.Fatalf("knockup vy = %v, want -4", got)
	}
}

func TestResolvePushbackSplit(t *testing.T) {
	g := hitGameState(attackState(100, 0, 0, 0, 10), idleState(105))
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.Position.X; got != 115 {
		t.Fatalf("defender x = %v, want 115", got)
	}
	if got := g.Characters[0].StateMachine.Position.X; got != 95 {
		t.Fatalf("attacker x = %v, want 95 (pushback/2 recoil)", got)
	}
}

func TestOneHitPerFrameEntry(t *testing.T) {
	g := hitGameState(attackState(100, 100, 0, 0, 0), idleState(105))
	g.ResolveHits()
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.HP; got != 9900 {
		t.Fatalf("second resolve on same frame re-hit: HP = %d, want 9900", got)
	}
}

func TestRehitOnFrameAdvance(t *testing.T) {
	atk := attackState(100, 100, 0, 0, 0)
	// Second active entry with its own hitbox.
	atk.AnimPlayer.Animations["A"].FrameData = append(
		atk.AnimPlayer.Animations["A"].FrameData,
		animation.FrameData{
			Duration: 100,
			Damage:   100,
			Boxes: map[types.BoxType][]types.Rect{
				types.Hit: {{X: 0, Y: -40, W: 60, H: 40}},
			},
		},
	)
	g := hitGameState(atk, idleState(105))
	g.ResolveHits()
	atk.AnimPlayer.FrameIndex = 1
	g.pruneConnects()
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.HP; got != 9800 {
		t.Fatalf("advanced frame should re-hit: HP = %d, want 9800", got)
	}
}

func TestTradeBothHitsLand(t *testing.T) {
	p1 := attackState(100, 100, 0, 0, 0)
	p2 := attackState(140, 100, 0, 0, 0)
	p2.IsFacingLeft = animation.Left
	g := hitGameState(p1, p2)
	g.ResolveHits()
	if got := g.Characters[0].StateMachine.HP; got != 9900 {
		t.Fatalf("P1 HP = %d, want 9900 (P2's trade hit)", got)
	}
	if got := g.Characters[1].StateMachine.HP; got != 9900 {
		t.Fatalf("P2 HP = %d, want 9900 (P1's trade hit)", got)
	}
}

func TestMissChangesNothing(t *testing.T) {
	g := hitGameState(attackState(100, 100, 5, 5, 10), idleState(500))
	g.ResolveHits()
	def := g.Characters[1].StateMachine
	if def.HP != 10000 || def.Velocity.X != 0 || def.Velocity.Y != 0 || def.Position.X != 500 {
		t.Fatalf("miss altered defender: %+v", def)
	}
	if len(g.Connects) != 0 {
		t.Fatalf("miss recorded %d ledger entries", len(g.Connects))
	}
}

// A zero-effect hitbox still connects (ledger records) and, per the
// zero-first rule, stops pre-existing momentum.
func TestZeroTouchZeroesVelocity(t *testing.T) {
	def := idleState(105)
	def.Velocity = types.Vector2{X: 5, Y: -3}
	g := hitGameState(attackState(100, 0, 0, 0, 0), def)
	g.ResolveHits()
	if got := g.Characters[1].StateMachine.Velocity; got.X != 0 || got.Y != 0 {
		t.Fatalf("zero touch should zero velocity, got %v", got)
	}
	if got := g.Characters[1].StateMachine.HP; got != 10000 {
		t.Fatalf("zero touch changed HP: %d", got)
	}
	if !g.HasConnected(0, 1, "A", 0) {
		t.Fatal("zero touch should still record the connect")
	}
}
