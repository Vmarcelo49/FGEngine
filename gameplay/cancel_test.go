package gameplay

import (
	"fgengine/animation"
	"testing"
)

func cancelTestSM(active string) *animation.StateMachine {
	ap := &animation.AnimationPlayer{
		Animations: map[string]*animation.Animation{
			"A": {FrameData: []animation.FrameData{{Duration: 1}}},
			"B": {FrameData: []animation.FrameData{{Duration: 1}}},
			"8": {FrameData: []animation.FrameData{{Duration: 1}}},
		},
	}
	ap.SetAnimation(active)
	return &animation.StateMachine{AnimPlayer: ap}
}

// "any" anywhere in the list matches everything (SPEC §6.5) — including
// non-first position, where the old [0]-only check failed (audit B4).
func TestCanCancelAnyInNonFirstPosition(t *testing.T) {
	sm := cancelTestSM("A")
	fd := &animation.FrameData{CancelTypes: []string{"A", "any"}}
	if !canCancelTo(fd, sm, "B") {
		t.Fatal("any in non-first position must match intent B")
	}
}

func TestCanCancelRules(t *testing.T) {
	sm := cancelTestSM("A")

	if !canCancelTo(&animation.FrameData{CancelTypes: []string{"any"}}, sm, "B") {
		t.Error("any must match")
	}
	if !canCancelTo(&animation.FrameData{CancelTypes: []string{"B"}}, sm, "B") {
		t.Error("exact listed intent must match")
	}
	if canCancelTo(&animation.FrameData{CancelTypes: []string{"B"}}, sm, "C") {
		t.Error("unlisted intent must not match")
	}
	if canCancelTo(&animation.FrameData{}, sm, "B") {
		t.Error("empty cancelTypes must not match")
	}
	if canCancelTo(&animation.FrameData{CancelTypes: []string{"any"}}, sm, "A") {
		t.Error("cancel into the active animation must not match")
	}
}
