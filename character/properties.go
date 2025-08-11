package character

import "FGEngine/types"

// FrameProperties represents static frame-specific data that varies per animation frame
type FrameProperties struct {
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
	HitType      hitType
	AnimPhase    AnimationPhase
	IsInvincible bool
	HasArmor     bool

	// Static hitboxes for this frame
	HitBoxes       []types.Rect
	HurtBoxes      []types.Rect
	CollisionBoxes []types.Rect
}

func (p *FrameProperties) CanBeCounterHit() bool {
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
