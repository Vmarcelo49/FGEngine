package logic

import (
	"fgengine/input"
	"fgengine/player"
)

// Having a gamestate that is updated by inputs is useful for replays and netplay
func UpdateByInputs(inputs []input.GameInput, players []*player.Player) {
	for playerNum, playerInput := range inputs {
		if playerInput.IsPressed(input.Right) {
			players[playerNum].Character.Position.X += 2
		}
		if playerInput.IsPressed(input.Left) {
			players[playerNum].Character.Position.X -= 2
		}
	}

}
