package gameplay

import (
	"fgengine/animation"
	"fgengine/types"
	"math"
)

// ResolveBodyCollision checks for collisions on the collision boxes, then moves the players out of each other if they overlap.
// All collision boxes of both frames participate; the first overlapping
// pair in slice order (P1 outer, P2 inner) wins (SPEC §6.6, §8.4).
func ResolveBodyCollision(p1, p2 *animation.StateMachine) {
	p1Boxes := collisionBoxesInWorld(p1)
	p2Boxes := collisionBoxesInWorld(p2)
	for _, p1CollisionBox := range p1Boxes {
		for _, p2CollisionBox := range p2Boxes {
			if p1CollisionBox.IsOverlapping(p2CollisionBox) {
				resolvePair(p1, p2, p1CollisionBox, p2CollisionBox)
				return
			}
		}
	}
}

func resolvePair(p1, p2 *animation.StateMachine, p1CollisionBox, p2CollisionBox types.Rect) {

	// center of the collision box
	ax, ay := p1CollisionBox.Center()
	bx, by := p2CollisionBox.Center()

	centerA := types.Vector2{X: ax, Y: ay}
	centerB := types.Vector2{X: bx, Y: by}

	delta := centerB.Sub(centerA)

	overlapX := math.Min(p1CollisionBox.Right(), p2CollisionBox.Right()) - math.Max(p1CollisionBox.X, p2CollisionBox.X)
	overlapY := math.Min(p1CollisionBox.Bottom(), p2CollisionBox.Bottom()) - math.Max(p1CollisionBox.Y, p2CollisionBox.Y)
	if overlapX <= 0 || overlapY <= 0 {
		return
	}

	var separationValue types.Vector2

	// resolve in the direction of least penetration 🥵
	if overlapX < overlapY {
		if delta.X > 0 {
			separationValue = types.Vector2{X: -overlapX, Y: 0}
		} else {
			separationValue = types.Vector2{X: overlapX, Y: 0}
		}

		resolveWithVelocity(p1, p2, separationValue, true)

	} else {
		if delta.Y > 0 {
			separationValue = types.Vector2{X: 0, Y: -overlapY}
		} else {
			separationValue = types.Vector2{X: 0, Y: overlapY}
		}

		resolveWithVelocity(p1, p2, separationValue, false)
	}
}

// collisionBoxesInWorld returns all collision boxes of the active frame in
// world coordinates, in slice order (SPEC §6.6).
func collisionBoxesInWorld(sm *animation.StateMachine) []types.Rect {
	if sm == nil || sm.AnimPlayer == nil {
		return nil
	}

	frameData := sm.AnimPlayer.ActiveFrameData()
	if frameData == nil {
		return nil
	}

	boxes := frameData.Boxes[types.Collision]
	out := make([]types.Rect, 0, len(boxes))
	for _, box := range boxes {
		if world, ok := boxInWorldCoordinates(box, sm); ok {
			out = append(out, world)
		}
	}
	return out
}

func resolveWithVelocity(a, b *animation.StateMachine, separationValue types.Vector2, isX bool) {
	// magnitude da velocidade
	la := math.Hypot(a.Velocity.X, a.Velocity.Y)
	lb := math.Hypot(b.Velocity.X, b.Velocity.Y)

	total := la + lb
	fa := 0.5
	fb := 0.5
	if total > 0 {
		fa = la / total
		fb = lb / total
	}

	// aplica separação
	aMove := separationValue.Mul(fa)
	bMove := separationValue.Mul(fb)

	a.Position.X += aMove.X
	a.Position.Y += aMove.Y

	b.Position.X -= bMove.X
	b.Position.Y -= bMove.Y

	applyPushVelocity(a, b, isX)
}

func applyPushVelocity(a, b *animation.StateMachine, isX bool) {
	if isX {
		ax := math.Abs(a.Velocity.X)
		bx := math.Abs(b.Velocity.X)

		if ax > bx {
			b.Velocity.X = a.Velocity.X
			a.Velocity.X *= 0.9
		} else if bx > ax {
			a.Velocity.X = b.Velocity.X
			b.Velocity.X *= 0.9
		}
		return
	}

	ay := math.Abs(a.Velocity.Y)
	by := math.Abs(b.Velocity.Y)

	if ay > by {
		b.Velocity.Y = a.Velocity.Y
		a.Velocity.Y *= 0.9
	} else if by > ay {
		a.Velocity.Y = b.Velocity.Y
		b.Velocity.Y *= 0.9
	}
}
