package gameplay

import (
	"fgengine/character"
	"fgengine/input"
)

type GameState struct {
	Characters [2]*character.Character
}

func (g *GameState) Update(inputs [2]input.GameInput) {
	p1 := g.Characters[0]
	p2 := g.Characters[1]

	p1.Update(inputs[0])
	p2.Update(inputs[1])
}

/*
char.update()
  func update()
    1. check input / state machine
    2. check physics (gravity, friction, velocity)
    3. check collisions
    4. check animation
*/
