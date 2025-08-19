package animation

import (
	"fgengine/state"
)

// FrameProperties represents static frame-specific data that varies per animation frame
type FrameProperties struct {
	// State represents the character state during this frame
	State state.State `yaml:"state"`

	CancelType AttackCancelType `yaml:"cancelType"`
	Priority   int              `yaml:"priority"` // maybe used in trades

	Damage    int `yaml:"damage"`
	Hitstun   int `yaml:"hitstun"`
	Blockstun int `yaml:"blockstun"`
	Pushback  int `yaml:"pushback"`
	Knockback int `yaml:"knockback"`
	Knockup   int `yaml:"knockup"`

	// walk, jump, dash, etc.
	ChangeXSpeed int `yaml:"changeXSpeed"`
	ChangeYSpeed int `yaml:"changeYSpeed"`

	// attack flags
	CanHardKnockdown bool `yaml:"canHardKnockdown"`
	CanWallBounce    bool `yaml:"canWallBounce"`
	CanGroundBounce  bool `yaml:"canGroundBounce"`
	CanOTG           bool `yaml:"canOTG"`

	SoundID int `yaml:"soundID"` // sound effect ID, 0 means no sound

	MoveType MoveType `yaml:"moveType"`
	// hit properties
	HitType      HitType        `yaml:"hitType"`
	AnimPhase    AnimationPhase `yaml:"animPhase"`
	IsInvincible bool           `yaml:"isInvincible"`
	HasArmor     bool           `yaml:"hasArmor"`
}

func (p *FrameProperties) CanBeCounterHit() bool {
	if p.MoveType == NonAttack || p.AnimPhase == Recovery || p.IsInvincible {
		return false
	}
	return true
}

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
