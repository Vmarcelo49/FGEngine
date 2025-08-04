package character

import "FGEngine/types"

// info that matters when the game is running
type CharacterState struct {
	Pos                 types.Vector2
	Vel                 types.Vector2
	HP                  int // 10000
	IgnoreGravityframes int // some moves ignore gravity for a few frames
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
	if p.MoveType == nonAttack || p.AnimPhase == recovery || p.IsInvincible {
		return false
	}
	return true
}

type attackCancelType int

const (
	cancelAll           attackCancelType = iota // when in idle?
	cancelNone                                  // while taking damage, no commands do anything
	cancelNormalAttack                          // from idle and other normals
	cancelSpecialAttack                         // from idle and normals
	cancelSuper                                 // from idle, normals and specials
	cancelJump                                  // from idle and normals
)

type hitType int

const (
	mid hitType = iota
	overhead
	low
	unblockable
)

type AnimationPhase int

const (
	startup AnimationPhase = iota
	active
	recovery
)

type MoveType int

const (
	nonAttack MoveType = iota
	normalAttack
	specialAttack
	superAttack
	grabAttack
	// maybe add a type for projectiles
)
