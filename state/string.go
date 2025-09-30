package state

var StateNames = map[State]string{
	StateGrounded:  "grounded",
	StateCrouching: "crouching",
	StateAirborne:  "airborne",
	StateDowned:    "downed",

	StateWalk:    "walking",
	StateDash:    "dashing",
	StateJump:    "jumping",
	StateFalling: "falling",

	StateForward:  "forward",
	StateBackward: "backward",
	StateNeutral:  "neutral",

	StateRecovery: "recovery",

	StateOnHitsun: "hitstun",
	StateHigh:     "high",
	StateBody:     "body",
	StateLow:      "low",
	StateSweep:    "sweep",

	StateBlock: "blocking",

	StateGrabbing: "grabbing",
	StateGrabbed:  "grabbed",
	StateGrabTech: "grabTech",

	StateWinAnimation:         "winAnimation",
	StateTimeoutLoseAnimation: "timeoutLoseAnimation",
	StateRoundStartAnimation:  "roundStartAnimation",

	StateAttack:        "attack",
	StateSpecialAttack: "specialAttack",
	StateSuperAttack:   "superAttack",

	StateA: "a",
	StateB: "b",
	StateC: "c",

	Modifier1: "modifier1",
	Modifier2: "modifier2",
	Modifier3: "modifier3",
}

// CompositeStateNames maps composite states to their string representations
var CompositeStateNames = map[State]string{
	StateIdle: "idle",

	StateWalkForward:     "walkForward",
	StateWalkBackward:    "walkBackward",
	StateDashForward:     "dashForward",
	StateDashBackward:    "dashBackward",
	StateAirDashForward:  "airDashForward",
	StateAirDashBackward: "airDashBackward",

	StateJumpNeutral:  "jumpNeutral",
	StateJumpForward:  "jumpForward",
	StateJumpBackward: "jumpBackward",
	StateNeutralFall:  "neutralFall",
	StateForwardFall:  "forwardFall",

	StateHighHit:      "highHit",
	StateBodyHit:      "bodyHit",
	StateLowHit:       "lowHit",
	StateAirHit:       "airHit",
	StateSweepFall:    "sweepFall",
	StateCrouchingHit: "crouchingHit",

	StateAttack5A:    "attack5A",
	StateAttack2A:    "attack2A",
	StateAttack5B:    "attack5B",
	StateAttack2B:    "attack2B",
	StateAttack5C:    "attack5C",
	StateAttack2C:    "attack2C",
	StateAttackJumpA: "attackJumpA",
	StateAttackJumpB: "attackJumpB",
	StateAttackJumpC: "attackJumpC",

	StateSpecial236A: "special236A",
	StateSpecial236B: "special236B",
	StateSpecial236C: "special236C",
	StateSpecial214A: "special214A",
	StateSpecial214B: "special214B",
	StateSpecial214C: "special214C",

	StateSuper236A: "super236A",
}

// String returns the string representation of a state.
// It first checks composite states, then individual states.
func (s State) String() string {
	if name, ok := CompositeStateNames[s]; ok {
		return name
	}
	if name, ok := StateNames[s]; ok {
		return name
	}
	return "unknown"
}

// Strings returns all active state flags as a slice of strings.
func (s State) Strings() []string {
	var activeStates []string

	if name, ok := CompositeStateNames[s]; ok {
		activeStates = append(activeStates, name)
	} else {
		for state, name := range StateNames {
			if s&state != 0 {
				activeStates = append(activeStates, name)
			}
		}
	}

	if len(activeStates) == 0 {
		return []string{"unknown"}
	}

	return activeStates
}
