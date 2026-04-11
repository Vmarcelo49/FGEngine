package gameplay

import (
	"fgengine/character"
	"fgengine/input"
)

type GameState struct {
	Characters [2]*character.Character
}

func (g *GameState) Update(inputs [2]input.GameInput) {
	g.Characters[0].StateMachine.Update(inputs[0], g.Characters[1].StateMachine)
	g.Characters[1].StateMachine.Update(inputs[1], g.Characters[0].StateMachine)
}

/*
char.update()
  func update()
    1. check input / state machine
    2. check physics (gravity, friction, velocity)
    3. check collisions
    4. check animation
*/
