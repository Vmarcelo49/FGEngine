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
	HP                  int           `yaml:"-"`
	Position            types.Vector2 `yaml:"-"`
	Velocity            types.Vector2 `yaml:"-"`
	IgnoreGravityFrames int           `yaml:"-"`
	IsFacingLeft        Orientation   `yaml:"-"`

	AnimPlayer *AnimationPlayer `yaml:"activeAnim"`
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
		sm.Velocity.X = 0
	} else if sm.Position.X > constants.WorldWidth {
		sm.Position.X = constants.WorldWidth
		sm.Velocity.X = 0
	}

	if sm.Position.Y < 0 {
		sm.Position.Y = 0
		if sm.Velocity.Y < 0 {
			sm.Velocity.Y = 0
		}
	} else if sm.Position.Y > constants.GroundLevelY {
		sm.Position.Y = constants.GroundLevelY
		sm.Velocity.Y = 0
	}
}
