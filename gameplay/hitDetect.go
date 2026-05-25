package gameplay

import (
	"fgengine/animation"
	"fgengine/types"
)

func CheckHits(p1, p2 *animation.StateMachine) {
	checkhit(p1, p2)
	checkhit(p2, p1)
}

func checkhit(thisPlayer, otherPlayer *animation.StateMachine) bool {
	if thisPlayer == nil || otherPlayer == nil || thisPlayer.AnimPlayer == nil || otherPlayer.AnimPlayer == nil {
		return false
	}

	thisFrameData := thisPlayer.AnimPlayer.ActiveFrameData()
	otherFrameData := otherPlayer.AnimPlayer.ActiveFrameData()
	if thisFrameData == nil || otherFrameData == nil {
		return false
	}

	for _, hitBox := range thisFrameData.Boxes[types.Hit] {
		hitBoxWorld, ok := boxInWorldCoordinates(hitBox, thisPlayer)
		if !ok {
			continue
		}

		for _, hurtBox := range otherFrameData.Boxes[types.Hurt] {
			hurtBoxWorld, ok := boxInWorldCoordinates(hurtBox, otherPlayer)
			if !ok {
				continue
			}

			if hitBoxWorld.IsOverlapping(hurtBoxWorld) {
				// fmt.Println("Hit detected!")
				// Apply hitstun or other effects here
				return true
			}
		}
	}
	return false
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
