package gameplay

import (
	"fgengine/animation"
	"fgengine/types"
	"testing"
)

func boxSM(x, y float64, boxes ...types.Rect) *animation.StateMachine {
	ap := &animation.AnimationPlayer{
		Animations: map[string]*animation.Animation{
			"idle": {FrameData: []animation.FrameData{{
				Duration: 10,
				Boxes:    map[types.BoxType][]types.Rect{types.Collision: boxes},
			}}},
		},
	}
	ap.SetAnimation("idle")
	return &animation.StateMachine{
		Position:   types.Vector2{X: x, Y: y},
		AnimPlayer: ap,
	}
}

func worldOverlapping(sm1, sm2 *animation.StateMachine) bool {
	for _, b1 := range collisionBoxesInWorld(sm1) {
		for _, b2 := range collisionBoxesInWorld(sm2) {
			if b1.IsOverlapping(b2) {
				return true
			}
		}
	}
	return false
}

func TestResolveBodyCollisionSeparates(t *testing.T) {
	p1 := boxSM(100, 382, types.Rect{X: -10, Y: -10, W: 20, H: 20})
	p2 := boxSM(105, 382, types.Rect{X: -10, Y: -10, W: 20, H: 20})
	if !worldOverlapping(p1, p2) {
		t.Fatal("fixture must start overlapped")
	}
	ResolveBodyCollision(p1, p2)
	if worldOverlapping(p1, p2) {
		t.Fatalf("still overlapped after resolve: %v %v", p1.Position, p2.Position)
	}
}

func TestResolveBodyCollisionConsidersAllBoxes(t *testing.T) {
	// First box is far away; only the second overlaps. The old
	// first-only code missed this case entirely.
	far := types.Rect{X: 500, Y: 500, W: 10, H: 10}
	near := types.Rect{X: -10, Y: -10, W: 20, H: 20}
	p1 := boxSM(100, 382, far, near)
	p2 := boxSM(105, 382, near)
	if !worldOverlapping(p1, p2) {
		t.Fatal("fixture must start overlapped")
	}
	ResolveBodyCollision(p1, p2)
	if worldOverlapping(p1, p2) {
		t.Fatalf("all-boxes resolve missed the second-box overlap: %v %v", p1.Position, p2.Position)
	}
}

func TestResolveBodyCollisionNoOp(t *testing.T) {
	p1 := boxSM(100, 382, types.Rect{X: -10, Y: -10, W: 20, H: 20})
	p2 := boxSM(500, 382, types.Rect{X: -10, Y: -10, W: 20, H: 20})
	ResolveBodyCollision(p1, p2)
	if p1.Position.X != 100 || p2.Position.X != 500 {
		t.Fatalf("separated players moved: %v %v", p1.Position, p2.Position)
	}
}
