package gameplay

// SplitMix64 is the single deterministic PRNG for simulation randomness
// (SPEC §3.4). Owned by GameState and seeded per match; its state is part
// of the state snapshot (§3.5).
type SplitMix64 struct {
	State uint64
}

// Next returns the next pseudo-random uint64.
func (r *SplitMix64) Next() uint64 {
	r.State += 0x9E3779B97F4A7C15
	z := r.State
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}
