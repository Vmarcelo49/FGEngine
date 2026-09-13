package gameplay

import (
	"encoding/binary"
	"fgengine/animation"
	"hash/fnv"
	"math"
)

// MatchPhase is the GameState-owned match flow state (SPEC §7.7).
// Transitions (fight → round-end → reset) land in F5; F0 only ticks the timer.
type MatchPhase byte

const (
	PhaseFight MatchPhase = iota
	PhaseRoundEnd
)

// ConnectKey records one landed hitbox for the one-hit-per-frame rule
// (SPEC §7.1). Entries clear when the attacker's frame advances.
type ConnectKey struct {
	Attacker int
	Anim     string
	Frame    int
	Defender int
}

// Hash returns the FNV-1a 64-bit hash over the canonical snapshot encoding
// (SPEC §3.5): fixed field order, floats by IEEE-754 bits. Slices iterate
// in insertion order, which the fixed pipeline order keeps deterministic.
// Animation names come from the stamped ActiveAnimation.Name field — never
// a map scan (SPEC §3.3).
func (g GameState) Hash() uint64 {
	h := fnv.New64a()
	var buf [8]byte
	putU64 := func(u uint64) {
		binary.LittleEndian.PutUint64(buf[:], u)
		_, _ = h.Write(buf[:])
	}
	putI64 := func(i int) {
		putU64(uint64(int64(i)))
	}
	putF64 := func(f float64) {
		putU64(math.Float64bits(f))
	}
	putString := func(s string) {
		putU64(uint64(len(s)))
		_, _ = h.Write([]byte(s))
	}
	putBool := func(b bool) {
		if b {
			putU64(1)
		} else {
			putU64(0)
		}
	}

	for i := 0; i < 2; i++ {
		sm := g.Characters[i].StateMachine
		putF64(sm.Position.X)
		putF64(sm.Position.Y)
		putF64(sm.Velocity.X)
		putF64(sm.Velocity.Y)
		putI64(sm.HP)
		putI64(sm.MaxHP)
		putI64(sm.StunFrames)
		putBool(sm.KnockdownPending)
		putBool(sm.WallBouncePending)
		putBool(sm.GroundBounceArmed)
		putBool(sm.GroundBounceUsed)
		putI64(sm.IgnoreGravityFrames)
		putBool(sm.IsFacingLeft == animation.Left)
		name, frameIndex, timeLeft := "", 0, 0
		var queue []string
		if sm.AnimPlayer != nil {
			if sm.AnimPlayer.ActiveAnimation != nil {
				name = sm.AnimPlayer.ActiveAnimation.Name
			}
			frameIndex = sm.AnimPlayer.FrameIndex
			timeLeft = sm.AnimPlayer.FrameTimeLeft
			queue = sm.AnimPlayer.AnimationQueue
		}
		putString(name)
		putI64(frameIndex)
		putI64(timeLeft)
		putU64(uint64(len(queue)))
		for _, q := range queue {
			putString(q)
		}
		putU64(uint64(len(g.inputHist[i])))
		for _, in := range g.inputHist[i] {
			putU64(uint64(in))
		}
	}

	putU64(uint64(len(g.Connects)))
	for _, c := range g.Connects {
		putI64(c.Attacker)
		putString(c.Anim)
		putI64(c.Frame)
		putI64(c.Defender)
	}

	putU64(g.RNG.State)
	putI64(g.TimerFrames)
	putI64(g.Round)
	putU64(uint64(g.Wins[0]))
	putU64(uint64(g.Wins[1]))
	putU64(uint64(g.Phase))
	putI64(g.FreezeFrames)
	return h.Sum64()
}

// pruneConnects drops ledger entries whose attacker has left the recorded
// frame (SPEC §7.1). Called at the start of every Update, before hit
// detection (SPEC §4.1 step 4).
func (g *GameState) pruneConnects() {
	kept := g.Connects[:0]
	for _, c := range g.Connects {
		if c.Attacker < 0 || c.Attacker > 1 {
			continue
		}
		name, frame := "", -1
		if sm := g.Characters[c.Attacker].StateMachine; sm != nil && sm.AnimPlayer != nil && sm.AnimPlayer.ActiveAnimation != nil {
			name = sm.AnimPlayer.ActiveAnimation.Name
			frame = sm.AnimPlayer.FrameIndex
		}
		if c.Anim == name && c.Frame == frame {
			kept = append(kept, c)
		}
	}
	for i := len(kept); i < len(g.Connects); i++ {
		g.Connects[i] = ConnectKey{}
	}
	g.Connects = kept
}

// HasConnected reports whether this attacker frame already landed on this
// defender (one-hit-per-frame rule, SPEC §7.1).
func (g *GameState) HasConnected(attacker, defender int, anim string, frame int) bool {
	for _, c := range g.Connects {
		if c.Attacker == attacker && c.Defender == defender && c.Anim == anim && c.Frame == frame {
			return true
		}
	}
	return false
}

// RecordConnect logs a landed hit for the one-hit-per-frame rule.
func (g *GameState) RecordConnect(attacker, defender int, anim string, frame int) {
	g.Connects = append(g.Connects, ConnectKey{Attacker: attacker, Anim: anim, Frame: frame, Defender: defender})
}
