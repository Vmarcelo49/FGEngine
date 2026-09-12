package gameplay

import (
	"fgengine/animation"
	"fgengine/types"
)

// ResolveHits checks both hit directions (trade-capable) and applies the F1
// subset of §7.2: damage, knockback/knockup velocity, pushback displacement,
// and the one-hit-per-frame ledger. Order is P1→P2 then P2→P1 (§4.1 step 4).
// Reaction-state entry + StunFrames (F2), guard (F3), invincibility/armor
// (F2), OTG (F4), and KO handling (F5) are not applied here.
func (g *GameState) ResolveHits() {
	g.resolveDirection(0, 1)
	g.resolveDirection(1, 0)
}

func (g *GameState) resolveDirection(ai, di int) {
	atk := g.Characters[ai].StateMachine
	def := g.Characters[di].StateMachine
	if atk == nil || def == nil || atk.AnimPlayer == nil || def.AnimPlayer == nil {
		return
	}
	afd := atk.AnimPlayer.ActiveFrameData()
	dfd := def.AnimPlayer.ActiveFrameData()
	if afd == nil || dfd == nil {
		return
	}
	if len(afd.Boxes[types.Hit]) == 0 || len(dfd.Boxes[types.Hurt]) == 0 {
		return
	}
	anim, frame := "", 0
	if atk.AnimPlayer.ActiveAnimation != nil {
		anim = atk.AnimPlayer.ActiveAnimation.Name
		frame = atk.AnimPlayer.FrameIndex
	}
	for _, hitBox := range afd.Boxes[types.Hit] {
		hitBoxWorld, ok := boxInWorldCoordinates(hitBox, atk)
		if !ok {
			continue
		}
		for _, hurtBox := range dfd.Boxes[types.Hurt] {
			hurtBoxWorld, ok := boxInWorldCoordinates(hurtBox, def)
			if !ok {
				continue
			}
			if !hitBoxWorld.IsOverlapping(hurtBoxWorld) {
				continue
			}
			if g.HasConnected(ai, di, anim, frame) {
				continue
			}
			g.applyHit(ai, di, anim, frame)
			return
		}
	}
}

// applyHit resolves one landed hit per §7.2 (F1 subset).
func (g *GameState) applyHit(ai, di int, anim string, frame int) {
	atk := g.Characters[ai].StateMachine
	def := g.Characters[di].StateMachine
	afd := atk.AnimPlayer.ActiveFrameData()

	// Away from the attacker in world X, by position (robust to facing).
	dir := 1.0
	if def.Position.X < atk.Position.X {
		dir = -1.0
	}

	// Step 3: damage (no clamp — F2; no KO handling — F5).
	def.HP -= afd.Damage

	// Step 5, velocity only: a hit overwrites momentum, then impulses add.
	def.Velocity.X = 0
	def.Velocity.Y = 0
	def.Velocity.X += dir * float64(afd.Knockback)
	def.Velocity.Y += -float64(afd.Knockup)

	// Step 6: one-time pushback displacement apart along X (integer
	// division on the attacker's share, SPEC §7.2).
	def.Position.X += dir * float64(afd.Pushback)
	atk.Position.X -= dir * float64(afd.Pushback/2)

	g.RecordConnect(ai, di, anim, frame)
}

func boxInWorldCoordinates(box types.Rect, sm *animation.StateMachine) (types.Rect, bool) {
	if sm == nil || sm.AnimPlayer == nil {
		return types.Rect{}, false
	}

	sprite := sm.AnimPlayer.ActiveSprite()
	anchor := types.Vector2{}
	if sprite != nil {
		anchor = sprite.Anchor
	}

	worldX := sm.Position.X + box.X - anchor.X
	if sm.IsFacingLeft == animation.Left {
		worldX = sm.Position.X - box.X - box.W + anchor.X
	}

	worldY := sm.Position.Y + box.Y - anchor.Y

	return types.Rect{X: worldX, Y: worldY, W: box.W, H: box.H}, true
}
