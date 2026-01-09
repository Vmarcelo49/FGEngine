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
	HitstunFrames        int
	Position             types.Vector2
	Velocity             types.Vector2
	IgnoreGravityFrames  int
	InputHistory         []input.GameInput
	CharacterOrientation FacingDirection
	MoveInput            int  // -1 left, 0 idle, 1 right
	JumpRequested        bool // queued jump until processed by character update
	AttackRequested      bool // queued attack until processed by character update
	DashRequested        bool // queued dash until processed by character update
}

// ClearState removes flags without toggling grounded/airborne like RemoveState does.
func (sm *StateMachine) ClearState(flags State) {
	sm.PreviousState = sm.ActiveState
	sm.ActiveState &^= flags
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
	sm.PreviousState = sm.ActiveState
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

func (sm *StateMachine) IsInactable() bool {
	return sm.HasAnyState(StateOnHitsun | StateGrabbed | StateWinAnimation | StateTimeoutLoseAnimation | StateRoundStartAnimation | StateDowned | StateRecovery | StateGrabbing | StateGrabTech | StateAttack | StateSpecialAttack | StateSuperAttack)
}

var MovementStates = StateWalk | StateDash | StateJump | StateFalling

func (sm *StateMachine) IsMoving() bool {
	return sm.HasAnyState(MovementStates)
}
