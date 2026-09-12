package animation

import (
	"fgengine/types"
)

// FrameData represents static frame-specific data that varies per animation frame
type FrameData struct {
	//stuff being used so far, TODO remove this when done with the refactor
	IncVelocityX float64 `yaml:"changeXSpeed,omitempty" toml:"changeXSpeed,omitempty"`
	IncVelocityY float64 `yaml:"changeYSpeed,omitempty" toml:"changeYSpeed,omitempty"`

	CancelTypes []string `yaml:"cancelTypes,omitempty" toml:"cancelTypes,omitempty"` // list of cancel types that can be used during this frame, e.g. "jump", "attack", "dash", etc.

	Boxes    map[types.BoxType][]types.Rect `yaml:"boxes,omitempty" toml:"boxes,omitempty"`
	Duration int                            `yaml:"duration" toml:"duration"`

	SpriteIndex int `yaml:"spriteIndex,omitempty" toml:"spriteIndex,omitempty"` // index of the sprite to display for this frame
	// unused stuff that will be used later
	IsRecovery bool `yaml:"isRecovery,omitempty" toml:"isRecovery,omitempty"` // whether this frame is a recovery frame (can't cancel out of it)
	IsHoldable bool `yaml:"isHoldable,omitempty" toml:"isHoldable,omitempty"` // whether this frame can be held by holding the input

	AnimationSwitch string `yaml:"animationSwitch,omitempty" toml:"animationSwitch,omitempty"` // switch to this animation after this frame ends
	Priority        int    `yaml:"priority,omitempty" toml:"priority,omitempty"`                 // maybe used in trades
	Damage          int    `yaml:"damage,omitempty" toml:"damage,omitempty"`
	Hitstun         int    `yaml:"hitstun,omitempty" toml:"hitstun,omitempty"`
	Blockstun       int    `yaml:"blockstun,omitempty" toml:"blockstun,omitempty"`
	Pushback        int    `yaml:"pushback,omitempty" toml:"pushback,omitempty"`
	Knockback       int    `yaml:"knockback,omitempty" toml:"knockback,omitempty"`
	Knockup         int    `yaml:"knockup,omitempty" toml:"knockup,omitempty"`

	CanHardKnockdown bool `yaml:"canHardKnockdown,omitempty" toml:"canHardKnockdown,omitempty"`
	CanWallBounce    bool `yaml:"canWallBounce,omitempty" toml:"canWallBounce,omitempty"`
	CanGroundBounce  bool `yaml:"canGroundBounce,omitempty" toml:"canGroundBounce,omitempty"`
	CanOTG           bool `yaml:"canOTG,omitempty" toml:"canOTG,omitempty"`
	CommonAudioID    int  `yaml:"soundID,omitempty" toml:"soundID,omitempty"`               // sound effect ID, 0 means no sound
	UniqueAudioID    int  `yaml:"uniqueSoundID,omitempty" toml:"uniqueSoundID,omitempty"`

	IsInvincible bool `yaml:"isInvincible,omitempty" toml:"isInvincible,omitempty"`
	HasArmor     bool `yaml:"hasArmor,omitempty" toml:"hasArmor,omitempty"`
}
