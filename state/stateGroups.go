package state

import "fgengine/input"

var UnactableStates = StateOnHitsun | StateGrabbed | StateWinAnimation | StateTimeoutLoseAnimation | StateRoundStartAnimation | StateDowned | StateRecovery | StateGrabbing | StateGrabTech | StateAttack | StateSpecialAttack | StateSuperAttack

var MovementStates = StateWalk | StateDash | StateJump | StateFalling

func (sm *StateMachine) HandleInput(playerInput input.GameInput) {
	if sm.HasAnyState(UnactableStates) {
		return
	}
}
