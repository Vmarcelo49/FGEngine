package gameplay

import (
	"fgengine/animation"
	"fgengine/input"
	"fgengine/types"
	"slices"
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
	// §7.6: ko takes none.
	if def.AnimPlayer.ActiveAnimation != nil && def.AnimPlayer.ActiveAnimation.Name == "ko" {
		return
	}
	// §7.2 step 1: invincibility.
	if dfd.IsInvincible {
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
			if blockState, guarded := g.checkGuard(di); guarded {
				g.applyBlock(ai, di, anim, frame, blockState)
				return
			}
			g.applyHit(ai, di, anim, frame)
			return
		}
	}
}

// applyHit resolves one landed hit per §7.2 (F1 core + F2 steps: clamp,
// armor, reaction entry, KO; guard is F3, OTG is F4, KO flow is F5).
func (g *GameState) applyHit(ai, di int, anim string, frame int) {
	atk := g.Characters[ai].StateMachine
	def := g.Characters[di].StateMachine
	afd := atk.AnimPlayer.ActiveFrameData()
	dfd := def.AnimPlayer.ActiveFrameData()
	dir := awayDir(g, ai, di)

	// Step 3: damage, clamped to [0, maxHP].
	def.HP -= afd.Damage
	if def.HP < 0 {
		def.HP = 0
	} else if def.HP > def.MaxHP {
		def.HP = def.MaxHP
	}

	// Step 4: armor absorbs the entire reaction (no hitstun, no
	// impulses, no pushback, no state change) — but the connect is still
	// recorded below.
	if dfd.HasArmor {
		g.RecordConnect(ai, di, anim, frame)
		return
	}

	// Step 5 velocity: a hit overwrites momentum, then impulses add.
	def.Velocity.X = 0
	def.Velocity.Y = 0
	def.Velocity.X += dir * float64(afd.Knockback)
	def.Velocity.Y += -float64(afd.Knockup)

	// Step 6: one-time pushback displacement apart along X (integer
	// division on the attacker's share, SPEC §7.2).
	def.Position.X += dir * float64(afd.Pushback)
	atk.Position.X -= dir * float64(afd.Pushback/2)

	// Step 5 reaction: KO takes precedence over hitstun; otherwise enter
	// the selected hitstun state with the attack's hitstun value.
	if def.HP <= 0 {
		def.StunFrames = 0
		def.AnimPlayer.SetAnimation("ko")
	} else {
		def.AnimPlayer.SetAnimation(selectHitstun(def))
		def.StunFrames = afd.Hitstun
	}

	g.RecordConnect(ai, di, anim, frame)
}

// selectHitstun picks the reaction state by defender posture (§6.7):
// airborne → airHurt, crouch states → crouchHurt, otherwise hurt.
func selectHitstun(def *animation.StateMachine) string {
	if def.IsAirborne() {
		return "airHurt"
	}
	if def.AnimPlayer.ActiveAnimation != nil {
		switch def.AnimPlayer.ActiveAnimation.Name {
		case "1", "2", "3":
			return "crouchHurt"
		}
	}
	return "hurt"
}

// Guardable postures (§7.4). Canonical names only (P2).
var groundGuardable = []string{"idle", "4", "1", "2", "3", "blockHit", "crouchBlock"}
var airGuardable = []string{"7", "8", "9", "fall", "airBlock"}

// guardHeld reports holding-back on the hit frame: the facing-corrected
// current input holds Left (§7.4). It reads the history tail, which the
// intake step (§4.1 step 2) guarantees is the current frame's input by the
// time hit detection (§4.1 step 4) runs.
func guardHeld(g *GameState, di int) bool {
	hist := g.inputHist[di]
	if len(hist) == 0 {
		return false
	}
	in := hist[len(hist)-1]
	if g.Characters[di].StateMachine.IsFacingLeft == animation.Left {
		if in&input.Left != 0 {
			in = (in &^ input.Left) | input.Right
		} else if in&input.Right != 0 {
			in = (in &^ input.Right) | input.Left
		}
	}
	return in&input.Left != 0
}

// checkGuard evaluates §7.4: back held + guardable posture. Returns the
// blockstun state to enter ("" = not guarded).
func (g *GameState) checkGuard(di int) (string, bool) {
	def := g.Characters[di].StateMachine
	if def == nil || def.AnimPlayer == nil || def.AnimPlayer.ActiveAnimation == nil {
		return "", false
	}
	if !guardHeld(g, di) {
		return "", false
	}
	name := def.AnimPlayer.ActiveAnimation.Name
	if def.IsAirborne() {
		if slices.Contains(airGuardable, name) {
			return "airBlock", true
		}
		return "", false
	}
	if !slices.Contains(groundGuardable, name) {
		return "", false
	}
	if name == "1" || name == "2" || name == "3" {
		return "crouchBlock", true
	}
	return "blockHit", true
}

// applyBlock resolves a guarded hit: no damage (chip 0, SPEC §7.4), the
// selected blockstun state with the attack's blockstun value, and the
// usual pushback displacement. The connect is recorded.
func (g *GameState) applyBlock(ai, di int, anim string, frame int, state string) {
	atk := g.Characters[ai].StateMachine
	def := g.Characters[di].StateMachine
	afd := atk.AnimPlayer.ActiveFrameData()
	dir := awayDir(g, ai, di)

	def.AnimPlayer.SetAnimation(state)
	def.StunFrames = afd.Blockstun

	def.Position.X += dir * float64(afd.Pushback)
	atk.Position.X -= dir * float64(afd.Pushback/2)

	g.RecordConnect(ai, di, anim, frame)
}

// awayDir returns +1 when the defender stands right of the attacker
// (impulses push them apart), -1 otherwise.
func awayDir(g *GameState, ai, di int) float64 {
	if g.Characters[di].StateMachine.Position.X < g.Characters[ai].StateMachine.Position.X {
		return -1.0
	}
	return 1.0
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
