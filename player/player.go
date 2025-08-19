package player

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/input"
	"fgengine/types"
	"fmt"
	"log"
)

type Player struct {
	// values that maybe make sense to be here
	ID           int
	Character    *animation.Character
	InputManager *input.InputManager
	HP           int

	// wtf is wrong with me values
	AnimationComponent  *animation.AnimationManager
	Position            types.Vector2
	Velocity            types.Vector2
	IgnoreGravityFrames int
	StateMachine        *StateMachine
}

// GetAllBoxes returns all boxes of the current player's sprite.
func (p *Player) GetAllBoxes() []collision.Box {
	if p.AnimationComponent == nil || !p.AnimationComponent.IsValid() {
		return []collision.Box{}
	}

	var boxes []collision.Box
	currentSprite := p.AnimationComponent.GetCurrentSprite()

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

// SetAnimation is a convenience method to set animations
func (p *Player) SetAnimation(animName string) {
	if p.AnimationComponent != nil {
		fmt.Println("Stop using player.SetAnimation()! use AnimationComponent.SetAnimation() instead.")
		p.AnimationComponent.SetAnimation(animName)
	}
}

// Makes a player with helmet for debugging
func CreateDebugPlayer(animManager *animation.AnimationRegistry) *Player {
	character, err := animation.LoadCharacter("./assets/characters/helmet.yaml")
	if err != nil {
		log.Fatal(err)
	}
	p1InputManager := input.NewInputManager()
	p1InputManager.AssignGamepadID(0) // TODO, this should check for available gamepads and return an error if none found

	// Create animation component through the manager
	animComponent := animManager.CreateComponent(character)

	return &Player{
		ID:                 0,
		Character:          character,
		InputManager:       p1InputManager,
		AnimationComponent: animComponent,
	}
}
