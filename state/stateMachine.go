package state

import (
	"fgengine/input"
	"fgengine/types"
)

type FacingDirection bool

const (
	Right FacingDirection = false
	Left  FacingDirection = true
)

type StateMachine struct {
	ActiveState          State
	PreviousState        State
	HP                   int
	Position             types.Vector2
	Velocity             types.Vector2
	IgnoreGravityFrames  int
	InputHistory         []input.GameInput
	CharacterOrientation FacingDirection
}

func (sm *StateMachine) AddState(stateToAdd State) {
	sm.PreviousState = sm.ActiveState
	sm.ActiveState |= stateToAdd
}

func (sm *StateMachine) SetState(newState State) {
	sm.PreviousState = sm.ActiveState
	sm.ActiveState = newState
}

// HasState checks if the state machine has ALL the specified state flags
func (sm *StateMachine) HasState(stateToCheck State) bool {
	return (sm.ActiveState & stateToCheck) == stateToCheck
}

// HasAnyState checks if the state machine has ANY of the specified state flags
func (sm *StateMachine) HasAnyState(stateToCheck State) bool {
	return (sm.ActiveState & stateToCheck) != 0
}

func (sm *StateMachine) RemoveState(stateToRemove State) {
	sm.ActiveState &= ^stateToRemove
	// guardrails to prevent charaters flying or being stuck
	if stateToRemove == StateAirborne {
		sm.ActiveState |= StateGrounded
		return
	}
	if stateToRemove == StateGrounded {
		sm.ActiveState |= StateAirborne
		return
	}
}

// sample of how to check if the character is attacking or other new functions may be added in this package
func (sm *StateMachine) IsAttacking() bool {
	return sm.HasAnyState(StateAttack | StateSpecialAttack | StateSuperAttack)
}
