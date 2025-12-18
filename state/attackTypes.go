package state

type AttackCancelType int

const (
	CancelAll           AttackCancelType = iota // when in idle?
	CancelNone                                  // while taking damage, no commands do anything
	CancelNormalAttack                          // from idle and other normals
	CancelSpecialAttack                         // from idle and normals
	CancelSuper                                 // from idle, normals and specials
	CancelJump                                  // from idle and normals
)

type HitType int

const (
	Medium HitType = iota
	Overhead
	Low
	Unblockable
)

type AnimationPhase int

const (
	Startup AnimationPhase = iota
	Active
	Recovery
)

type MoveType int

const (
	NonAttack MoveType = iota
	NormalAttack
	SpecialAttack
	SuperAttack
	GrabAttack
	// maybe add a type for projectiles
)
