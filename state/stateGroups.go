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

	// reset queued actions each frame
	sm.MoveInput = 0
	sm.JumpRequested = false
	sm.AttackRequested = false
	sm.DashRequested = false
	sm.ClearState(StateA | StateB | StateC)

	if sm.IsInactable() {
		return // this is terrible
	}

	// check for input sequences first
	for key, seq := range input.InputSequences { // cooldown here probably would be good
		if input.DetectInputSequence(seq, sm.InputHistory) {
			if reflect.DeepEqual(seq, input.InputSequences[key]) {
				sm.DashRequested = true
			}
		}
	}

	// directional inputs (world space)
	if playerInput.IsPressed(input.Left) {
		sm.MoveInput = -1
	}
	if playerInput.IsPressed(input.Right) {
		sm.MoveInput = 1
	}

	// button states
	if playerInput.IsPressed(input.A) {
		sm.AttackRequested = true
		sm.AddState(StateA)
	}
	if playerInput.IsPressed(input.B) {
		sm.AddState(StateB)
	}
	if playerInput.IsPressed(input.C) {
		sm.AddState(StateC)
	}

	if playerInput.IsPressed(input.Up) && sm.HasState(StateGrounded) {
		sm.JumpRequested = true
	}
}
