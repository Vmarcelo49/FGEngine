package state

type StateMachine struct {
	State         State
	PreviousState State
}

func (sm *StateMachine) AddState(stateToAdd State) {
	sm.PreviousState = sm.State
	sm.State |= stateToAdd
}

func (sm *StateMachine) SetState(newState State) {
	sm.PreviousState = sm.State
	sm.State = newState
}

// HasState checks if the state machine has ALL the specified state flags
func (sm *StateMachine) HasState(stateToCheck State) bool {
	return (sm.State & stateToCheck) == stateToCheck
}

// HasAnyState checks if the state machine has ANY of the specified state flags
func (sm *StateMachine) HasAnyState(stateToCheck State) bool {
	return (sm.State & stateToCheck) != 0
}

func (sm *StateMachine) RemoveState(stateToRemove State) {
	sm.State &= ^stateToRemove
	// guardrails to prevent charaters flying or being stuck
	if stateToRemove == StateAirborne {
		sm.State |= StateGrounded
		return
	}
	if stateToRemove == StateGrounded {
		sm.State |= StateAirborne
		return
	}
}

// sample of how to check if the character is attacking or other new functions may be added in this package
func (sm *StateMachine) IsAttacking() bool {
	return sm.HasAnyState(StateAttack | StateSpecialAttack | StateSuperAttack)
}
