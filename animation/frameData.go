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

func (fd *FrameData) switchToAnim(detectedAnimations []string, sm *StateMachine) {
	if len(fd.CancelTypes) == 0 {
		return
	}
	if fd.CancelTypes[0] == "any" {
		anim := filterAnimations(detectedAnimations)
		sm.ActiveAnim.SetAnimation(anim)
	}
	for _, cancelType := range fd.CancelTypes {
		for _, detectedAnim := range detectedAnimations {
			if cancelType == detectedAnim {
				sm.ActiveAnim.SetAnimation(detectedAnim)
				return
			}
		}
	}
	// it would be better to each animation having its own identifier instead of relying on the name
}

// placeholder that returns the last detected animation, we have to make a priority system later
func filterAnimations(detectedAnimations []string) string {
	if len(detectedAnimations) == 0 {
		return ""
	}
	return detectedAnimations[len(detectedAnimations)-1]
}
