package animation

import (
	"fgengine/constants"
	"fgengine/input"
	"fgengine/types"
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

func (sm *StateMachine) Update(inputs input.GameInput) {
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
			sm.ActiveAnim.SetAnimation("4", false)
		} else {
			sm.ActiveAnim.SetAnimation("6", false)
		}
	case input.Right:
		if sm.Facing == Left {
			sm.ActiveAnim.SetAnimation("6", false)
		} else {
			sm.ActiveAnim.SetAnimation("4", false)
		}
	case input.Up:
		// check if grounded before allowing jump
		// check()
		sm.ActiveAnim.SetAnimation("jump", false)
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
	// check gravity and friction here

}
