package player

import (
	"FGEngine/character"
	"FGEngine/collision"
	"FGEngine/input"
	"FGEngine/types"
)

type Player struct {
	ID           int
	Character    *character.Character
	InputManager *input.InputManager

	State *PlayerState
}

type PlayerState struct {
	Position            types.Vector2
	Velocity            types.Vector2
	HP                  int
	IgnoreGravityFrames int

	// Runtime animation state
	AnimationManager *AnimationManager
	StateMachine     *StateMachine

	// Current frame properties (derived from current animation frame)
	CurrentFrameProps *character.FrameProperties
}

// GetAllBoxes returns all boxes of the current player's sprite.
func (p *Player) GetAllBoxes() []collision.Box {
	if p.State.AnimationManager.CurrentSprite == nil {
		return []collision.Box{}
	}

	var boxes []collision.Box
	currentSprite := p.State.AnimationManager.CurrentSprite

	for _, boxRect := range currentSprite.CollisionBoxes {
		boxes = append(boxes, collision.Box{Rect: boxRect, BoxType: collision.Collision})
	}
	for _, boxRect := range currentSprite.HitBoxes {
		boxes = append(boxes, collision.Box{Rect: boxRect, BoxType: collision.Hit})
	}
	for _, boxRect := range currentSprite.HurtBoxes {
		boxes = append(boxes, collision.Box{Rect: boxRect, BoxType: collision.Hurt})
	}
	return boxes
}
