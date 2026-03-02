package logic

import (
	"fgengine/input"
	"fgengine/player"
)

// Having a gamestate that is updated by inputs is useful for replays and netplay
func UpdateByInputs(inputs []input.GameInput, players []*player.Player) {
	for playerNum, playerInput := range inputs {
		players[playerNum].Character.StateMachine.HandleInput(playerInput)
		players[playerNum].Character.Update()
	}
}
