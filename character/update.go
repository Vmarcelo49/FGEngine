package character

import (
	"fgengine/constants"
	"fgengine/state"
)

func (c *Character) Update() {
	// Update animation frame timing first
	c.updateAnimation()

	// Get current frame properties for dynamic values
	frameData := c.AnimationPlayer.GetActiveFrameData()
	if frameData == nil {
		panic("Frame properties should not be nil if animation is set")
	}

	switch c.StateMachine.ActiveState {
	case state.StateIdle:
		c.setAnimation("idle")
	case state.StateWalk | state.StateForward:
		c.setAnimation("walk")
	case state.StateDash | state.StateForward:
		c.setAnimation("dash")
	default:
		c.setAnimation("idle")
	}

	// Apply frame-specific velocity changes if available
	c.StateMachine.Velocity.X = float64(frameData.ChangeXSpeed)
	c.StateMachine.Velocity.Y = float64(frameData.ChangeYSpeed)

	// Update position based on velocity
	c.StateMachine.Position.X += c.StateMachine.Velocity.X
	c.StateMachine.Position.Y += c.StateMachine.Velocity.Y

	// Apply friction
	if c.StateMachine.Velocity.X > 0 {
		c.StateMachine.Velocity.X -= float64(c.Friction)
		if c.StateMachine.Velocity.X < 0 {
			c.StateMachine.Velocity.X = 0
		}
	} else if c.StateMachine.Velocity.X < 0 {
		c.StateMachine.Velocity.X += float64(c.Friction)
		if c.StateMachine.Velocity.X > 0 {
			c.StateMachine.Velocity.X = 0
		}
	}

	// Apply gravity only if airborne and not ignoring gravity frames
	if c.StateMachine.HasState(state.StateAirborne) && c.StateMachine.IgnoreGravityFrames <= 0 {
		c.StateMachine.Velocity.Y += constants.Gravity // example gravity value
	} else {
		c.StateMachine.IgnoreGravityFrames--
	}

	if c.StateMachine.Position.Y >= constants.GroundLevelY {
		c.StateMachine.Position.Y = constants.GroundLevelY
		c.StateMachine.Velocity.Y = 0
		c.StateMachine.AddState(state.StateGrounded)
		c.StateMachine.RemoveState(state.StateAirborne)
	}

	if c.StateMachine.Position.X < constants.World.X {
		c.StateMachine.Position.X = constants.World.X
	}
	if c.StateMachine.Position.X+float64(c.GetSprite().Rect.W) > constants.World.Right() {
		c.StateMachine.Position.X = constants.World.Right() - float64(c.GetSprite().Rect.W)
	}
}
