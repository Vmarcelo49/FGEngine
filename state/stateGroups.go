package state

import (
	"fgengine/constants"
	"fgengine/input"
	"reflect"
)

func (sm *StateMachine) HandleInput(playerInput input.GameInput) {
	sm.InputHistory = append(sm.InputHistory, playerInput)
	// remove old inputs
	if len(sm.InputHistory) > constants.MaxInputHistory {
		sm.InputHistory = sm.InputHistory[len(sm.InputHistory)-constants.MaxInputHistory:]
	}
	if sm.IsInactable() {
		return
	}

	// check for input sequences first
	for key, seq := range input.InputSequences { // cooldown here probably would be good
		if input.DetectInputSequence(seq, sm.InputHistory) {
			if reflect.DeepEqual(seq, input.InputSequences[key]) {
				sm.SetState(StateDash | StateForward)
				return
			}
		}
	}

	// then check for single inputs
	if playerInput == input.Right {
		sm.AddState(StateWalkForward)
	}

	if playerInput == input.NoInput {
		sm.RemoveState(MovementStates) // remove movement states when no input
		sm.AddState(StateNeutral)
	}
}
