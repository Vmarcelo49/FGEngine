package player

import (
	"fgengine/character"
	"fgengine/input"
	"log"
)

type Player struct {
	Character *character.Character
	Input     *input.InputManager
}

// Makes a player with helmet for debugging
func NewDebugPlayer() *Player {
	character, err := character.LoadCharacter(character.Helmet)
	if err != nil {
		log.Fatal(err)
	}
	p1InputManager := input.NewInputManager()
	p1InputManager.AssignGamepadID(0) // TODO, this should check for available gamepads and return an error if none found

	return &Player{
		Character: character,
		Input:     p1InputManager,
	}
}

// DO NOT USE THIS YET, UNFINISHED
func NewPlayer() *Player {
	player := &Player{}

	return player
}

// NewPlayerWithInput builds a player using the debug character and the provided input.
func NewPlayerWithInput(manager *input.InputManager) (*Player, error) {
	chara, err := character.LoadCharacter(character.Helmet)
	if err != nil {
		return nil, err
	}
	if manager == nil {
		manager = input.NewInputManager()
	}

	return &Player{
		Character: chara,
		Input:     manager,
	}, nil
}
