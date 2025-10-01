package state

import (
	"fgengine/constants"
	"fgengine/input"
	"reflect"
)

var UnactableStates = StateOnHitsun | StateGrabbed | StateWinAnimation | StateTimeoutLoseAnimation | StateRoundStartAnimation | StateDowned | StateRecovery | StateGrabbing | StateGrabTech | StateAttack | StateSpecialAttack | StateSuperAttack

var MovementStates = StateWalk | StateDash | StateJump | StateFalling

func (sm *StateMachine) HandleInput(playerInput input.GameInput) {
	sm.InputHistory = append(sm.InputHistory, playerInput)
	// remove old inputs
	if len(sm.InputHistory) > constants.MaxInputHistory {
		sm.InputHistory = sm.InputHistory[len(sm.InputHistory)-constants.MaxInputHistory:]
	}
	if sm.HasAnyState(UnactableStates) {
		return
	}

	// check for input sequences first
	for _, seq := range InputSequences { // cooldown here probably would be good
		if DetectInputSequence(seq, sm.InputHistory) {
			// if detected, set the state accordingly
			if reflect.DeepEqual(seq, InputSequences["66"]) {
				sm.SetState(StateDash | StateForward)
				return
			}
			if reflect.DeepEqual(seq, InputSequences["236A"]) {
				sm.SetState(StateSpecialAttack)
				return
			}
		}
	}

	// then check for single inputs
	if playerInput == input.Right {
		sm.AddState(StateWalk | StateForward)
	}

	if playerInput == input.NoInput {
		sm.RemoveState(MovementStates) // remove movement states when no input
		sm.AddState(StateNeutral)
	}
}
