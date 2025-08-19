package state

// Each State constant is a unique flag (bitmask). Use bitwise operations to compose and test states.
type State uint64

const (
	// core positions
	StateGrounded State = 1 << iota
	StateCrouching
	StateAirborne
	StateDowned

	// movement states
	StateWalk
	StateDash
	StateJump
	StateFalling

	// orientation
	StateForward
	StateBackward
	StateNeutral
	// recovery states
	StateRecovery

	// hit states
	StateOnHitsun // being hit
	// orientation for hit states
	StateHigh
	StateBody
	StateLow
	StateSweep

	// block states
	StateBlock

	// grab states
	StateGrabbing
	StateGrabbed
	StateGrabTech

	// win/lose/round start states
	StateWinAnimation
	StateTimeoutLoseAnimation
	StateRoundStartAnimation

	// attack states
	StateAttack
	StateSpecialAttack
	StateSuperAttack

	// buttons
	StateA
	StateB
	StateC

	// some other random stuff to complement specials
	Modifier1
	Modifier2
	Modifier3
)

// composite states - combinations of base states
const (
	StateIdle            = StateGrounded | StateNeutral
	StateWalkForward     = StateGrounded | StateWalk | StateForward
	StateWalkBackward    = StateGrounded | StateWalk | StateBackward
	StateDashForward     = StateGrounded | StateDash | StateForward
	StateDashBackward    = StateGrounded | StateDash | StateBackward
	StateAirDashForward  = StateAirborne | StateDash | StateForward
	StateAirDashBackward = StateAirborne | StateDash | StateBackward

	StateJumpNeutral  = StateAirborne | StateJump | StateNeutral
	StateJumpForward  = StateAirborne | StateJump | StateForward
	StateJumpBackward = StateAirborne | StateJump | StateBackward

	StateNeutralFall = StateAirborne | StateFalling | StateNeutral // used for neutral and back jump
	StateForwardFall = StateAirborne | StateFalling | StateForward

	StateHighHit      = StateGrounded | StateOnHitsun | StateHigh
	StateBodyHit      = StateGrounded | StateOnHitsun | StateBody
	StateLowHit       = StateGrounded | StateOnHitsun | StateLow
	StateAirHit       = StateAirborne | StateOnHitsun
	StateSweepFall    = StateAirborne | StateOnHitsun | StateSweep | StateFalling
	StateCrouchingHit = StateCrouching | StateOnHitsun

	// attacks
	StateAttack5A    = StateGrounded | StateAttack | StateA
	StateAttack2A    = StateCrouching | StateAttack | StateA
	StateAttack5B    = StateGrounded | StateAttack | StateB
	StateAttack2B    = StateCrouching | StateAttack | StateB
	StateAttack5C    = StateGrounded | StateAttack | StateC
	StateAttack2C    = StateCrouching | StateAttack | StateC
	StateAttackJumpA = StateAirborne | StateAttack | StateA
	StateAttackJumpB = StateAirborne | StateAttack | StateB
	StateAttackJumpC = StateAirborne | StateAttack | StateC

	// special attacks
	StateSpecial236A = StateGrounded | StateSpecialAttack | Modifier1 | StateA
	StateSpecial236B = StateGrounded | StateSpecialAttack | Modifier1 | StateB
	StateSpecial236C = StateGrounded | StateSpecialAttack | Modifier1 | StateC
	StateSpecial214A = StateGrounded | StateSpecialAttack | Modifier2 | StateA
	StateSpecial214B = StateGrounded | StateSpecialAttack | Modifier2 | StateB
	StateSpecial214C = StateGrounded | StateSpecialAttack | Modifier2 | StateC

	// super attacks
	StateSuper236A = StateGrounded | StateSuperAttack | StateB | StateC
)
