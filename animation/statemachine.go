package animation

import (
	"fgengine/constants"
	"fgengine/input"
	"fgengine/types"
	"math"
)

const (
	horizontalFriction  = 0.80
	minHorizontalSpeed  = 0.05
	maxHorizontalSpeedX = 6.0

	maxVerticalSpeedY = 10.0
)

type StateMachine struct {
	//ActiveState         State
	//PreviousState       State
	HP                  int               `yaml:"-"`
	Position            types.Vector2     `yaml:"-"`
	Velocity            types.Vector2     `yaml:"-"`
	IgnoreGravityFrames int               `yaml:"-"`
	InputHistory        []input.GameInput `yaml:"-"`
	Facing              Orientation       `yaml:"-"`

	ActiveAnim *AnimationPlayer `yaml:"activeAnim"`
}

type Orientation bool

const (
	Right = false
	Left  = true
)

type State struct {
	Name string
}

func (sm *StateMachine) Update(inputs input.GameInput) {
	if sm.ActiveAnim == nil {
		return
	}

	// Update input history
	sm.InputHistory = append(sm.InputHistory, inputs)
	if len(sm.InputHistory) > constants.MaxInputHistory {
		sm.InputHistory = sm.InputHistory[1:] // remove oldest input
	}

	// Check special moves
	specialCommand := input.CheckSpecialMove(sm.InputHistory)
	// if sm.IsActable(){}
	if specialCommand != "" {
		// should make checks if the special move can be performed (e.g. not in the middle of another move)
		sm.ActiveAnim.SetAnimation(specialCommand, false)
		/*
			prevState := sm.ActiveState
			sm.PreviousState = prevState
			sm.ActiveAnim.SetStateByAnimation(specialCommand)
		*/
	}
	directionalInput := inputs & 0b1111 // first 4 bits are directional inputs
	switch directionalInput {
	case input.Left:
		if sm.Facing == Right {
			sm.ActiveAnim.SetAnimation("4", true)
		} else {
			sm.ActiveAnim.SetAnimation("6", true)
		}
	case input.Right:
		if sm.Facing == Left {
			sm.ActiveAnim.SetAnimation("4", true)
		} else {
			sm.ActiveAnim.SetAnimation("6", true)
		}
	case input.Up:
		// check if grounded before allowing jump
		// check()
		sm.ActiveAnim.SetAnimation("jump", false)
	default:
		sm.ActiveAnim.SetAnimation("idle", true)
	}

	// get stuff from the animation frame data and apply it to the state machine (e.g. velocity changes, hitboxes, etc.)
	frameData := sm.ActiveAnim.ActiveFrameData()
	if frameData != nil {
		// apply velocity changes from frame data
		sm.Velocity.X += frameData.IncVelocityX
		sm.Velocity.Y += frameData.IncVelocityY

		// Play audio if specified there too
		// audio.Play(frameData.CommonAudioID, frameData.UniqueAudioID)
	}
	sm.applyFrictionGravity(directionalInput)

	sm.ActiveAnim.Update()
}

func (sm *StateMachine) applyFrictionGravity(directionalInput input.GameInput) {
	if sm.Velocity.X > maxHorizontalSpeedX {
		sm.Velocity.X = maxHorizontalSpeedX
	} else if sm.Velocity.X < -maxHorizontalSpeedX {
		sm.Velocity.X = -maxHorizontalSpeedX
	}

	sm.Velocity.X *= horizontalFriction
	if math.Abs(sm.Velocity.X) < minHorizontalSpeed {
		sm.Velocity.X = 0
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
