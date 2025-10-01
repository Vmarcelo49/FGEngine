package animation

import (
	"fgengine/state"
)

// FrameProperties represents static frame-specific data that varies per animation frame
type FrameProperties struct {
	Duration         int              `yaml:"duration"`
	AnimationSwitch  string           `yaml:"animationSwitch,omitempty"` // switch to this animation after this frame ends
	State            state.State      `yaml:"state,omitempty"`
	CancelType       AttackCancelType `yaml:"cancelType,omitempty"`
	Priority         int              `yaml:"priority,omitempty"` // maybe used in trades
	Damage           int              `yaml:"damage,omitempty"`
	Hitstun          int              `yaml:"hitstun,omitempty"`
	Blockstun        int              `yaml:"blockstun,omitempty"`
	Pushback         int              `yaml:"pushback,omitempty"`
	Knockback        int              `yaml:"knockback,omitempty"`
	Knockup          int              `yaml:"knockup,omitempty"`
	ChangeXSpeed     int              `yaml:"changeXSpeed,omitempty"`
	ChangeYSpeed     int              `yaml:"changeYSpeed,omitempty"`
	CanHardKnockdown bool             `yaml:"canHardKnockdown,omitempty"`
	CanWallBounce    bool             `yaml:"canWallBounce,omitempty"`
	CanGroundBounce  bool             `yaml:"canGroundBounce,omitempty"`
	CanOTG           bool             `yaml:"canOTG,omitempty"`
	CommonAudioID    int              `yaml:"soundID,omitempty"` // sound effect ID, 0 means no sound
	UniqueAudioID    int              `yaml:"uniqueSoundID,omitempty"`
	MoveType         MoveType         `yaml:"moveType,omitempty"`
	HitType          HitType          `yaml:"hitType,omitempty"`
	AnimPhase        AnimationPhase   `yaml:"animPhase,omitempty"`
	IsInvincible     bool             `yaml:"isInvincible,omitempty"`
	HasArmor         bool             `yaml:"hasArmor,omitempty"`
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
