package player

import (
	"FGEngine/character"
	"FGEngine/collision"
	"FGEngine/input"
	"FGEngine/types"
	"log"
)

type Player struct {
	ID           int
	Character    *character.Character
	InputManager *input.InputManager

	State            *PlayerState
	AnimationManager *AnimationManager
}

type PlayerState struct {
	Position            types.Vector2
	Velocity            types.Vector2
	HP                  int
	IgnoreGravityFrames int

	// Runtime animation state

	StateMachine *StateMachine

	// Current frame properties (derived from current animation frame)
	CurrentFrameProps *character.FrameProperties
}

// GetAllBoxes returns all boxes of the current player's sprite.
func (p *Player) GetAllBoxes() []collision.Box {
	if p.AnimationManager.CurrentSprite == nil {
		return []collision.Box{}
	}

	var boxes []collision.Box
	currentSprite := p.AnimationManager.CurrentSprite

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

// Makes a player with helmet for debugging
func CreateDebugPlayer() *Player {
	character, err := character.LoadCharacter("./assets/characters/helmet.yaml")
	if err != nil {
		log.Fatal(err)
	}
	p1InputManager := input.NewInputManager()
	p1InputManager.AssignGamepadID(0) // TODO, this should check for available gamepads and return an error if none found
	return &Player{
		ID:           0,
		Character:    character,
		InputManager: p1InputManager,
		State:        &PlayerState{},
	}
}
