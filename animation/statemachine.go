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

// Initial states should include: idle, walk, jump, attack, hitstun.
type State struct {
	Name    string
	OnEnter func(sm *StateMachine)
	OnExit  func(sm *StateMachine)
}

// Update follows the general flow of([X] means done):
// 1. Update input history and check for special move commands. [X]
// 2. Update the active animation based on current inputs and state.
// 3. Apply velocity changes from the current animation frame data. [X]
// 4. Apply friction and gravity to update the position and velocity of the character. [X]
// Then it should do in the future:
// 5. Check for state transitions based on the new position, velocity, and inputs.
// 6. Handle hitboxes and collisions based on the current animation frame data.
// 7. Play any audio associated with the current animation frame.
func (sm *StateMachine) Update(inputs input.GameInput) {
	if sm.ActiveAnim == nil {
		return
	}

	// Update input history
	sm.InputHistory = append(sm.InputHistory, inputs)
	if len(sm.InputHistory) > constants.MaxInputHistory {
		sm.InputHistory = sm.InputHistory[1:] // remove oldest input
	}

	detectedAnimations := []string{}
	// Check special moves
	specialCommand := input.CheckSpecialMove(sm.InputHistory)
	// if sm.IsActable(){}
	if specialCommand != "" {
		detectedAnimations = append(detectedAnimations, specialCommand)
	}

	directionalInput := inputs & 0b1111 // first 4 bits are directional inputs
	switch directionalInput {
	case input.Left:
		if sm.Facing == Right {
			detectedAnimations = append(detectedAnimations, "4")
		} else {
			detectedAnimations = append(detectedAnimations, "6")
		}
	case input.Right:
		if sm.Facing == Left {
			detectedAnimations = append(detectedAnimations, "4")
		} else {
			detectedAnimations = append(detectedAnimations, "6")
		}
	case input.Up:
		detectedAnimations = append(detectedAnimations, "jump")
	}

	// get stuff from the animation frame data and apply it to the state machine (e.g. velocity changes, hitboxes, etc.)
	frameData := sm.ActiveAnim.ActiveFrameData()
	if frameData != nil {
		// apply velocity changes from frame data
		sm.Velocity.X += frameData.IncVelocityX
		sm.Velocity.Y += frameData.IncVelocityY

		// Process cancel routes declared by the active frame.
		if len(detectedAnimations) > 0 {
			frameData.switchToAnim(detectedAnimations, sm)
		}

		// Play audio if specified there too
		// audio.Play(frameData.CommonAudioID, frameData.UniqueAudioID)
	}

	// Return to idle if nothing detected and current animation is finished.
	if len(detectedAnimations) == 0 && sm.ActiveAnim.IsFinished() {
		sm.ActiveAnim.SetAnimation("idle")
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
