package state

// OrderedStates defines the canonical order for displaying individual state flags
var OrderedStates = []State{
	// Core positions
	StateGrounded,
	StateCrouching,
	StateAirborne,
	StateDowned,

	// Movement states
	StateWalk,
	StateDash,
	StateJump,
	StateFalling,

	// Orientation
	StateForward,
	StateBackward,
	StateNeutral,

	// Recovery states
	StateRecovery,

	// Hit states
	StateOnHitsun,
	StateHigh,
	StateBody,
	StateLow,
	StateSweep,

	// Block states
	StateBlock,

	// Grab states
	StateGrabbing,
	StateGrabbed,
	StateGrabTech,

	// Win/lose/round states
	StateWinAnimation,
	StateTimeoutLoseAnimation,
	StateRoundStartAnimation,

	// Attack states
	StateAttack,
	StateSpecialAttack,
	StateSuperAttack,

	// Button states
	StateA,
	StateB,
	StateC,

	// Modifier states
	Modifier1,
	Modifier2,
	Modifier3,
}
