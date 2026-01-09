package logic

import (
	"fgengine/input"
	"fgengine/player"
	"fgengine/state"
)

// Having a gamestate that is updated by inputs is useful for replays and netplay
func UpdateByInputs(inputs []input.GameInput, players []*player.Player) {
	for playerNum, playerInput := range inputs {
		players[playerNum].Character.StateMachine.HandleInput(playerInput)
		players[playerNum].Character.Update()
	}
}

func UpdateFacings(players []*player.Player) {
	if len(players) < 2 {
		return
	}
	p1 := players[0].Character
	p2 := players[1].Character
	if p1.Position().X <= p2.Position().X {
		p1.StateMachine.CharacterOrientation = state.Right
		p2.StateMachine.CharacterOrientation = state.Left
	} else {
		p1.StateMachine.CharacterOrientation = state.Left
		p2.StateMachine.CharacterOrientation = state.Right
	}
}
