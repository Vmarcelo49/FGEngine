package animation

import (
	"testing"
)

func stunTestPlayer(start, end int) *AnimationPlayer {
	ap := &AnimationPlayer{
		Animations: map[string]*Animation{
			"hurt": {
				FrameData: []FrameData{
					{Duration: 2},
					{Duration: 2},
					{Duration: 2},
					{Duration: 2},
				},
				LoopFrames: &LoopFrame{Start: start, End: end},
			},
		},
	}
	ap.SetAnimation("hurt")
	return ap
}

// While stunned, the player must hold within the loopFrames range and
// never advance past it (SPEC §7.6).
func TestStunHoldLoopsRange(t *testing.T) {
	ap := stunTestPlayer(1, 2)
	// Prime: burn the entry frame's duration (index 0 shows briefly).
	ap.Update("", 5)
	ap.Update("", 5)
	for i := 0; i < 12; i++ {
		ap.Update("", 5)
		if ap.FrameIndex < 1 || ap.FrameIndex > 2 {
			t.Fatalf("step %d: stunned frame escaped hold range: %d", i, ap.FrameIndex)
		}
	}
}

// start == end holds a single frame (SPEC §6.4, §7.6).
func TestStunHoldSingleFrame(t *testing.T) {
	ap := stunTestPlayer(2, 2)
	// Prime: burn the entry frame's duration (index 0 shows briefly).
	ap.Update("", 5)
	ap.Update("", 5)
	for i := 0; i < 8; i++ {
		ap.Update("", 5)
		if ap.FrameIndex != 2 {
			t.Fatalf("step %d: single-frame hold left frame 2: %d", i, ap.FrameIndex)
		}
	}
}

// Stunned without loopFrames: no hold available, normal update applies
// (must not panic, must still advance).
func TestStunWithoutLoopFramesFallsThrough(t *testing.T) {
	ap := &AnimationPlayer{
		Animations: map[string]*Animation{
			"hurt": {FrameData: []FrameData{{Duration: 1}, {Duration: 1}}},
		},
	}
	ap.SetAnimation("hurt")
	for i := 0; i < 6; i++ {
		ap.Update("", 5)
	}
	if ap.FrameIndex != 1 {
		t.Fatalf("expected advance to last frame, got %d", ap.FrameIndex)
	}
}

// Unstunned update keeps the pre-existing advance behavior.
func TestUnstunnedUpdateAdvances(t *testing.T) {
	ap := stunTestPlayer(1, 2)
	for i := 0; i < 20; i++ {
		ap.Update("", 0)
	}
	if ap.FrameIndex != 3 {
		t.Fatalf("expected unstunned release to run past the loop range to frame 3, got %d", ap.FrameIndex)
	}
}
