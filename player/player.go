package player

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/input"
	"log"
)

type Player struct {
	Character    *character.Character
	InputManager *input.InputManager
}

// Makes a player with helmet for debugging
func CreateDebugPlayer(animManager *animation.AnimationRegistry) *Player {
	character, err := character.LoadCharacter(character.Helmet)
	if err != nil {
		log.Fatal(err)
	}
	p1InputManager := input.NewInputManager()
	p1InputManager.AssignGamepadID(0) // TODO, this should check for available gamepads and return an error if none found

	return &Player{
		Character:    character,
		InputManager: p1InputManager,
	}
}
