package character

import "FGEngine/types"

// info that matters when the game is running
type CharacterState struct {
	Position            types.Vector2
	Velocity            types.Vector2
	HP                  int // 10000
	IgnoreGravityFrames int // some moves ignore gravity for a few frames
	StateMachine        StateMachine
}

type Properties struct {
	CancelType attackCancelType
	Priority   int // maybe used in trades

	Damage    int
	Hitstun   int
	Blockstun int
	Pushback  int
	Knockback int
	Knockup   int

	// walk, jump, dash, etc.
	ChangeXSpeed int
	ChangeYSpeed int

	// attack flags
	CanHardKnockdown bool
	CanWallBounce    bool
	CanGroundBounce  bool
	CanOTG           bool

	SoundID int // sound effect ID, 0 means no sound

	MoveType MoveType
	// hit properties
	HitType        hitType
	AnimPhase      AnimationPhase
	IsInvincible   bool
	HasArmor       bool
	State          State
	HitBoxes       []types.Rect
	HurtBoxes      []types.Rect
	CollisionBoxes []types.Rect
}

func (p *Properties) CanBeCounterHit() bool {
	if p.MoveType == NonAttack || p.AnimPhase == Recovery || p.IsInvincible {
		return false
	}
	return true
}

type attackCancelType int

const (
	CancelAll           attackCancelType = iota // when in idle?
	CancelNone                                  // while taking damage, no commands do anything
	CancelNormalAttack                          // from idle and other normals
	CancelSpecialAttack                         // from idle and normals
	CancelSuper                                 // from idle, normals and specials
	CancelJump                                  // from idle and normals
)

type hitType int

const (
	Medium hitType = iota
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
