package gameplay

import (
	"testing"
)

// Same seed must produce the identical stream on independent instances.
func TestSplitMix64Deterministic(t *testing.T) {
	a := SplitMix64{State: 12345}
	b := SplitMix64{State: 12345}
	for i := 0; i < 8; i++ {
		if got, want := a.Next(), b.Next(); got != want {
			t.Fatalf("step %d: same seed diverged: %d != %d", i, got, want)
		}
	}
}

// Different seeds must (overwhelmingly certainly) diverge immediately.
func TestSplitMix64SeedSensitive(t *testing.T) {
	a := SplitMix64{State: 1}
	b := SplitMix64{State: 2}
	for i := 0; i < 4; i++ {
		if got, want := a.Next(), b.Next(); got == want {
			t.Fatalf("step %d: different seeds produced identical output %d", i, got)
		}
	}
}
