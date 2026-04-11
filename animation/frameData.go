package animation

import (
	"fgengine/collision"
	"fgengine/types"
)

// FrameData represents static frame-specific data that varies per animation frame
type FrameData struct {
	//stuff being used so far, TODO remove this when done with the refactor
	IncVelocityX float64 `yaml:"changeXSpeed,omitempty"`
	IncVelocityY float64 `yaml:"changeYSpeed,omitempty"`

	CancelTypes []string `yaml:"cancelTypes,omitempty"` // list of cancel types that can be used during this frame, e.g. "jump", "attack", "dash", etc.

	// unused stuff that will be used later
	Boxes           map[collision.BoxType][]types.Rect `yaml:"boxes,omitempty"`
	Duration        int                                `yaml:"duration"`
	SpriteIndex     int                                `yaml:"spriteIndex,omitempty"`     // index of the sprite to display for this frame
	AnimationSwitch string                             `yaml:"animationSwitch,omitempty"` // switch to this animation after this frame ends
	//State            State                              `yaml:"state,omitempty"`
	//CancelType       state.AttackCancelType `yaml:"cancelType,omitempty"`
	//MoveType         MoveType       `yaml:"moveType,omitempty"`
	//HitType          HitType        `yaml:"hitType,omitempty"`
	//Phase            AnimationPhase `yaml:"animPhase,omitempty"`
	Priority  int `yaml:"priority,omitempty"` // maybe used in trades
	Damage    int `yaml:"damage,omitempty"`
	Hitstun   int `yaml:"hitstun,omitempty"`
	Blockstun int `yaml:"blockstun,omitempty"`
	Pushback  int `yaml:"pushback,omitempty"`
	Knockback int `yaml:"knockback,omitempty"`
	Knockup   int `yaml:"knockup,omitempty"`

	CanHardKnockdown bool `yaml:"canHardKnockdown,omitempty"`
	CanWallBounce    bool `yaml:"canWallBounce,omitempty"`
	CanGroundBounce  bool `yaml:"canGroundBounce,omitempty"`
	CanOTG           bool `yaml:"canOTG,omitempty"`
	CommonAudioID    int  `yaml:"soundID,omitempty"` // sound effect ID, 0 means no sound
	UniqueAudioID    int  `yaml:"uniqueSoundID,omitempty"`

	IsInvincible bool `yaml:"isInvincible,omitempty"`
	HasArmor     bool `yaml:"hasArmor,omitempty"`
}

func (fd *FrameData) switchToAnim(detectedAnimations string, sm *StateMachine) {
	if detectedAnimations == "" {
		return
	}

	if sm.ActiveAnim != nil && sm.ActiveAnim.ActiveAnimationName() == detectedAnimations {
		return
	}

	if len(fd.CancelTypes) == 0 {
		return
	}

	// Prevent jump-start animations while already airborne.
	if (detectedAnimations == "7" || detectedAnimations == "8" || detectedAnimations == "9") && sm.IsAirborne() {
		return
	}

	if fd.CancelTypes[0] == "any" {
		sm.ActiveAnim.SetAnimation(detectedAnimations)
		return
	}

	for _, cancelType := range fd.CancelTypes {
		if cancelType == detectedAnimations {
			sm.ActiveAnim.SetAnimation(detectedAnimations)
			return
		}
	}
	// it would be better to each animation having its own identifier instead of relying on the name
}
