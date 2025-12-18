package animation

import (
	"fgengine/collision"
	"fgengine/state"
	"fgengine/types"
)

// FrameData represents static frame-specific data that varies per animation frame
type FrameData struct {
	Boxes            map[collision.BoxType][]types.Rect `yaml:"boxes,omitempty"`
	Duration         int                                `yaml:"duration"`
	SpriteIndex      int                                `yaml:"spriteIndex,omitempty"`     // index of the sprite to display for this frame
	AnimationSwitch  string                             `yaml:"animationSwitch,omitempty"` // switch to this animation after this frame ends
	State            state.State                        `yaml:"state,omitempty"`
	CancelType       state.AttackCancelType             `yaml:"cancelType,omitempty"`
	Priority         int                                `yaml:"priority,omitempty"` // maybe used in trades
	Damage           int                                `yaml:"damage,omitempty"`
	Hitstun          int                                `yaml:"hitstun,omitempty"`
	Blockstun        int                                `yaml:"blockstun,omitempty"`
	Pushback         int                                `yaml:"pushback,omitempty"`
	Knockback        int                                `yaml:"knockback,omitempty"`
	Knockup          int                                `yaml:"knockup,omitempty"`
	ChangeXSpeed     float64                            `yaml:"changeXSpeed,omitempty"`
	ChangeYSpeed     float64                            `yaml:"changeYSpeed,omitempty"`
	CanHardKnockdown bool                               `yaml:"canHardKnockdown,omitempty"`
	CanWallBounce    bool                               `yaml:"canWallBounce,omitempty"`
	CanGroundBounce  bool                               `yaml:"canGroundBounce,omitempty"`
	CanOTG           bool                               `yaml:"canOTG,omitempty"`
	CommonAudioID    int                                `yaml:"soundID,omitempty"` // sound effect ID, 0 means no sound
	UniqueAudioID    int                                `yaml:"uniqueSoundID,omitempty"`
	MoveType         state.MoveType                     `yaml:"moveType,omitempty"`
	HitType          state.HitType                      `yaml:"hitType,omitempty"`
	Phase            state.AnimationPhase               `yaml:"animPhase,omitempty"`
	IsInvincible     bool                               `yaml:"isInvincible,omitempty"`
	HasArmor         bool                               `yaml:"hasArmor,omitempty"`
}
