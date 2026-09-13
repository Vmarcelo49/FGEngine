package animation

import (
	"fgengine/types"
)

// FrameData represents static frame-specific data that varies per animation frame
type FrameData struct {
	//stuff being used so far, TODO remove this when done with the refactor
	IncVelocityX float64 `toml:"changeXSpeed,omitempty"`
	IncVelocityY float64 `toml:"changeYSpeed,omitempty"`

	CancelTypes []string `toml:"cancelTypes,omitempty"` // list of cancel types that can be used during this frame, e.g. "jump", "attack", "dash", etc.

	Boxes    map[types.BoxType][]types.Rect `toml:"boxes,omitempty"`
	Duration int                            `toml:"duration"`

	SpriteIndex int `toml:"spriteIndex,omitempty"` // index of the sprite to display for this frame
	// unused stuff that will be used later
	IsRecovery bool `toml:"isRecovery,omitempty"` // whether this frame is a recovery frame (can't cancel out of it)
	IsHoldable bool `toml:"isHoldable,omitempty"` // whether this frame can be held by holding the input

	AnimationSwitch string `toml:"animationSwitch,omitempty"` // switch to this animation after this frame ends
	Priority        int    `toml:"priority,omitempty"`        // maybe used in trades
	Damage          int    `toml:"damage,omitempty"`
	Hitstun         int    `toml:"hitstun,omitempty"`
	Blockstun       int    `toml:"blockstun,omitempty"`
	Pushback        int    `toml:"pushback,omitempty"`
	Knockback       int    `toml:"knockback,omitempty"`
	Knockup         int    `toml:"knockup,omitempty"`

	CanHardKnockdown bool `toml:"canHardKnockdown,omitempty"`
	CanWallBounce    bool `toml:"canWallBounce,omitempty"`
	CanGroundBounce  bool `toml:"canGroundBounce,omitempty"`
	CanOTG           bool `toml:"canOTG,omitempty"`
	CommonAudioID    int  `toml:"soundID,omitempty"` // sound effect ID, 0 means no sound
	UniqueAudioID    int  `toml:"uniqueSoundID,omitempty"`

	IsInvincible bool `toml:"isInvincible,omitempty"`
	HasArmor     bool `toml:"hasArmor,omitempty"`
}
