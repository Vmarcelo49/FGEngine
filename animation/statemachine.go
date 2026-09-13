package animation

import (
	"fgengine/constants"
	"fgengine/types"
	"math"
)

const (
	horizontalFriction  = 0.80
	minHorizontalSpeed  = 0.05
	maxHorizontalSpeedX = 999

	maxVerticalSpeedY = 10.0
)

type StateMachine struct {
	//ActiveState         State
	//PreviousState       State
	HP                  int           `toml:"-"`
	MaxHP               int           `toml:"-"`
	Position            types.Vector2 `toml:"-"`
	Velocity            types.Vector2 `toml:"-"`
	IgnoreGravityFrames int           `toml:"-"`
	IsFacingLeft        Orientation   `toml:"-"`
	// StunFrames counts down hitstun/blockstun, loaded from the attack's
	// hitstun/blockstun value itself (SPEC §7.6). Ticked in §4.1 intake.
	StunFrames int `toml:"-"`
	// Launch memory (SPEC §7.5): set from the latest launching hit in
	// applyHit (overwrite semantics), consumed or cleared on ground
	// contact, cleared as stale while grounded without landing.
	KnockdownPending  bool `toml:"-"`
	WallBouncePending bool `toml:"-"`
	GroundBounceArmed bool `toml:"-"`
	GroundBounceUsed  bool `toml:"-"`

	AnimPlayer *AnimationPlayer `toml:"-"`
}

type Orientation bool

const (
	Right = false
	Left  = true
)

func (sm *StateMachine) IsAirborne() bool {
	return sm.Position.Y < constants.GroundLevelY
}

// ApplyVelocity applies movement deltas from the current frame data.
func (sm *StateMachine) ApplyVelocity() {
	frameData := sm.AnimPlayer.ActiveFrameData()
	incVelX := frameData.IncVelocityX
	if sm.IsFacingLeft {
		incVelX = -incVelX
	}
	sm.Velocity.X += incVelX
	sm.Velocity.Y += frameData.IncVelocityY

}

func (sm *StateMachine) ApplyPhysics() {
	if !sm.IsAirborne() {
		if sm.Velocity.X > maxHorizontalSpeedX {
			sm.Velocity.X = maxHorizontalSpeedX
		} else if sm.Velocity.X < -maxHorizontalSpeedX {
			sm.Velocity.X = -maxHorizontalSpeedX
		}

		sm.Velocity.X *= horizontalFriction
		if math.Abs(sm.Velocity.X) < minHorizontalSpeed {
			sm.Velocity.X = 0
		}
	}
	// Apply simple gravity while in the air.
	if sm.IgnoreGravityFrames > 0 {
		sm.IgnoreGravityFrames--
	} else if sm.Position.Y < constants.GroundLevelY || sm.Velocity.Y < 0 {
		sm.Velocity.Y += constants.Gravity
		if sm.Velocity.Y > maxVerticalSpeedY {
			sm.Velocity.Y = maxVerticalSpeedY
		}
	}

	// Integrate velocity into world position once per frame.
	sm.Position.X += sm.Velocity.X
	sm.Position.Y += sm.Velocity.Y

	// Keep character inside world bounds.
	if sm.Position.X < 0 {
		sm.Position.X = 0
		// Wall bounce: reflect with damping instead of stopping. The
		// pending flag is the significance signal (no speed threshold);
		// it survives wall bounces and clears on ground contact (§7.5).
		if sm.WallBouncePending && sm.Velocity.X < 0 {
			sm.Velocity.X = -sm.Velocity.X * 0.6
		} else {
			sm.Velocity.X = 0
		}
	} else if sm.Position.X > constants.WorldWidth {
		sm.Position.X = constants.WorldWidth
		if sm.WallBouncePending && sm.Velocity.X > 0 {
			sm.Velocity.X = -sm.Velocity.X * 0.6
		} else {
			sm.Velocity.X = 0
		}
	}

	if sm.Position.Y < 0 {
		sm.Position.Y = 0
		if sm.Velocity.Y < 0 {
			sm.Velocity.Y = 0
		}
	} else if sm.Position.Y > constants.GroundLevelY {
		// Ground bounce, once per launch: reflect with damping and rest
		// 1px above ground so the airborne invariant holds and no landing
		// event fires (§7.5).
		if sm.GroundBounceArmed && !sm.GroundBounceUsed && sm.Velocity.Y > 0 {
			sm.Position.Y = constants.GroundLevelY - 1
			sm.Velocity.Y = -sm.Velocity.Y * 0.4
			sm.GroundBounceUsed = true
		} else {
			sm.Position.Y = constants.GroundLevelY
			sm.Velocity.Y = 0
		}
	}
}
