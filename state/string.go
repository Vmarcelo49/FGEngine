package state

var StateNames = map[State]string{
	// Core positions
	StateGrounded:  "grounded",
	StateCrouching: "crouching",
	StateAirborne:  "airborne",
	StateDowned:    "downed",

	// Movement states
	StateWalk:    "walking",
	StateDash:    "dashing",
	StateJump:    "jumping",
	StateFalling: "falling",

	// Orientation
	StateForward:  "forward",
	StateBackward: "backward",
	StateNeutral:  "neutral",

	// Recovery states
	StateRecovery: "recovery",

	// Hit states
	StateOnHitsun: "hitstun",
	StateHigh:     "high",
	StateBody:     "body",
	StateLow:      "low",
	StateSweep:    "sweep",

	// Block states
	StateBlock: "blocking",

	// Grab states
	StateGrabbing: "grabbing",
	StateGrabbed:  "grabbed",
	StateGrabTech: "grabTech",

	// Win/lose/round states
	StateWinAnimation:         "winAnimation",
	StateTimeoutLoseAnimation: "timeoutLoseAnimation",
	StateRoundStartAnimation:  "roundStartAnimation",

	// Attack states
	StateAttack:        "attack",
	StateSpecialAttack: "specialAttack",
	StateSuperAttack:   "superAttack",

	// Button states
	StateA: "a",
	StateB: "b",
	StateC: "c",

	// Misc states
	Modifier1: "modifier1",
	Modifier2: "modifier2",
	Modifier3: "modifier3",
}

// CompositeStateNames maps composite states to their string representations
var CompositeStateNames = map[State]string{
	// Basic states
	StateIdle: "idle",

	// Movement states
	StateWalkForward:     "walkForward",
	StateWalkBackward:    "walkBackward",
	StateDashForward:     "dashForward",
	StateDashBackward:    "dashBackward",
	StateAirDashForward:  "airDashForward",
	StateAirDashBackward: "airDashBackward",

	// Jump states
	StateJumpNeutral:  "jumpNeutral",
	StateJumpForward:  "jumpForward",
	StateJumpBackward: "jumpBackward",
	StateNeutralFall:  "neutralFall",
	StateForwardFall:  "forwardFall",

	// Hit states
	StateHighHit:      "highHit",
	StateBodyHit:      "bodyHit",
	StateLowHit:       "lowHit",
	StateAirHit:       "airHit",
	StateSweepFall:    "sweepFall",
	StateCrouchingHit: "crouchingHit",

	// Attack states
	StateAttack5A:    "attack5A",
	StateAttack2A:    "attack2A",
	StateAttack5B:    "attack5B",
	StateAttack2B:    "attack2B",
	StateAttack5C:    "attack5C",
	StateAttack2C:    "attack2C",
	StateAttackJumpA: "attackJumpA",
	StateAttackJumpB: "attackJumpB",
	StateAttackJumpC: "attackJumpC",

	// Special attack states
	StateSpecial236A: "special236A",
	StateSpecial236B: "special236B",
	StateSpecial236C: "special236C",
	StateSpecial214A: "special214A",
	StateSpecial214B: "special214B",
	StateSpecial214C: "special214C",

	// Super attack states
	StateSuper236A: "super236A",
}

// String returns the string representation of a state.
// It first checks composite states, then individual states.
func (s State) String() string {
	// Check composite states first (they're more specific)
	if name, ok := CompositeStateNames[s]; ok {
		return name
	}
	// Check individual states
	if name, ok := StateNames[s]; ok {
		return name
	}
	return "unknown"
}
