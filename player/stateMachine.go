package player

import "fgengine/state"

type StateMachine struct {
	State         state.State
	PreviousState state.State
}

func (sm *StateMachine) AddState(stateToAdd state.State) {
	sm.PreviousState = sm.State
	sm.State |= stateToAdd
}

func (sm *StateMachine) SetState(newState state.State) {
	sm.PreviousState = sm.State
	sm.State = newState
}

// HasState checks if the state machine has ALL the specified state flags
func (sm *StateMachine) HasState(stateToCheck state.State) bool {
	return (sm.State & stateToCheck) == stateToCheck
}

// HasAnyState checks if the state machine has ANY of the specified state flags
func (sm *StateMachine) HasAnyState(stateToCheck state.State) bool {
	return (sm.State & stateToCheck) != 0
}

func (sm *StateMachine) RemoveState(stateToRemove state.State) {
	sm.State &= ^stateToRemove
	// guardrails to prevent charaters flying or being stuck
	if stateToRemove == state.StateAirborne {
		sm.State |= state.StateGrounded
		return
	}
	if stateToRemove == state.StateGrounded {
		sm.State |= state.StateAirborne
		return
	}
}

func (sm *StateMachine) IsAttacking() bool {
	return sm.HasAnyState(state.StateAttack | state.StateSpecialAttack | state.StateSuperAttack)
}
