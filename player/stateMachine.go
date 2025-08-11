package player

// Each State constant is a unique flag (bitmask). Use bitwise operations to compose and test states.
type State uint64

const (
	// core positions
	stateGrounded State = 1 << iota
	stateCrouching
	stateAirborne
	stateDowned

	// movement states
	stateWalk
	stateDash
	stateJump
	stateFalling

	// orientation
	stateForward
	stateBackward
	stateNeutral
	// recovery states
	stateRecovery

	// hit states
	stateOnHitsun // being hit
	// orientation for hit states
	stateHigh
	stateBody
	stateLow
	stateSweep

	// block states
	stateBlock

	// grab states
	stateGrabbing
	stateGrabbed
	stateGrabTech

	// win/lose/round start states
	stateWinAnimation
	stateTimeoutLoseAnimation
	stateRoundStartAnimation

	// attack states
	stateAttack
	stateSpecialAttack
	stateSuperAttack

	// buttons
	stateA
	stateB
	stateC

	// some other random stuff to complement specials
	modifier1
	modifier2
	modifier3
)

// composite states - combinations of base states
const (
	idle            = stateGrounded | stateNeutral
	walkForward     = stateGrounded | stateWalk | stateForward
	walkBackward    = stateGrounded | stateWalk | stateBackward
	dashForward     = stateGrounded | stateDash | stateForward
	dashBackward    = stateGrounded | stateDash | stateBackward
	airDashForward  = stateAirborne | stateDash | stateForward
	airDashBackward = stateAirborne | stateDash | stateBackward

	jumpNeutral  = stateAirborne | stateJump | stateNeutral
	jumpForward  = stateAirborne | stateJump | stateForward
	jumpBackward = stateAirborne | stateJump | stateBackward

	neutralFall = stateAirborne | stateFalling | stateNeutral // used for neutral and back jump
	forwardFall = stateAirborne | stateFalling | stateForward

	highHit      = stateGrounded | stateOnHitsun | stateHigh
	bodyHit      = stateGrounded | stateOnHitsun | stateBody
	lowHit       = stateGrounded | stateOnHitsun | stateLow
	airHit       = stateAirborne | stateOnHitsun
	sweepFall    = stateAirborne | stateOnHitsun | stateSweep | stateFalling
	crouchingHit = stateCrouching | stateOnHitsun

	// attacks
	attack5A    = stateGrounded | stateAttack | stateA
	attack2A    = stateCrouching | stateAttack | stateA
	attack5B    = stateGrounded | stateAttack | stateB
	attack2B    = stateCrouching | stateAttack | stateB
	attack5C    = stateGrounded | stateAttack | stateC
	attack2C    = stateCrouching | stateAttack | stateC
	attackJumpA = stateAirborne | stateAttack | stateA
	attackJumpB = stateAirborne | stateAttack | stateB
	attackJumpC = stateAirborne | stateAttack | stateC

	// special attacks
	special236A = stateGrounded | stateSpecialAttack | modifier1 | stateA
	special236B = stateGrounded | stateSpecialAttack | modifier1 | stateB
	special236C = stateGrounded | stateSpecialAttack | modifier1 | stateC
	special214A = stateGrounded | stateSpecialAttack | modifier2 | stateA
	special214B = stateGrounded | stateSpecialAttack | modifier2 | stateB
	special214C = stateGrounded | stateSpecialAttack | modifier2 | stateC

	// super attacks
	super236A = stateGrounded | stateSuperAttack | stateB | stateC
)

type StateMachine struct {
	state         State
	previousState State
}

func (sm *StateMachine) AddState(state State) {
	sm.previousState = sm.state
	sm.state |= state
}

func (sm *StateMachine) SetState(newState State) {
	sm.previousState = sm.state
	sm.state = newState
}

// HasState checks if the state machine has ALL the specified state flags
func (sm *StateMachine) HasState(state State) bool {
	return (sm.state & state) == state
}

// HasAnyState checks if the state machine has ANY of the specified state flags
func (sm *StateMachine) HasAnyState(state State) bool {
	return (sm.state & state) != 0
}

func (sm *StateMachine) RemoveState(state State) {
	sm.state &= ^state
	// guardrails to prevent charaters flying or being stuck
	if state == stateAirborne {
		sm.state |= stateGrounded
		return
	}
	if state == stateGrounded {
		sm.state |= stateAirborne
		return
	}
}

func (sm *StateMachine) IsAttacking() bool {
	return sm.HasAnyState(stateAttack | stateSpecialAttack | stateSuperAttack)
}
