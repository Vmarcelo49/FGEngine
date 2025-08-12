package player

import (
	"FGEngine/animation"
	"FGEngine/character"
	"FGEngine/collision"
	"FGEngine/input"
	"FGEngine/types"
	"log"
)

type Player struct {
	PlayerState

	ID                 int
	Character          *character.Character
	InputManager       *input.InputManager
	AnimationComponent *animation.AnimationComponent
}

type PlayerState struct {
	Position            types.Vector2
	Velocity            types.Vector2
	HP                  int
	IgnoreGravityFrames int
	StateMachine        *StateMachine
	CurrentFrameProps   *character.FrameProperties // Current frame properties (derived from current animation frame)
}

// GetAllBoxes returns all boxes of the current player's sprite.
// Implements the Renderable interface
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

// GetPosition returns the player's position. Implements the Renderable interface
func (p *Player) GetPosition() types.Vector2 {
	return p.Position
}

// GetAnimationComponent returns the animation component. Implements the Renderable interface
func (p *Player) GetAnimationComponent() *animation.AnimationComponent {
	return p.AnimationComponent
}

// GetID returns the player's ID. Implements the Renderable interface
func (p *Player) GetID() int {
	return p.ID
}

// SetAnimation is a convenience method to set animations
func (p *Player) SetAnimation(animName string) {
	if p.AnimationComponent != nil {
		p.AnimationComponent.SetAnimation(animName)
	}
}

// Makes a player with helmet for debugging
func CreateDebugPlayer(animManager *animation.ComponentManager) *Player {
	character, err := character.LoadCharacter("./assets/characters/helmet.yaml")
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
		PlayerState:        PlayerState{},
	}
}
