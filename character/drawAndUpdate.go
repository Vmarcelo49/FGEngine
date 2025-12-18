package character

import (
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/state"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Character) Draw(screen *ebiten.Image, camera *graphics.Camera) {
	sprite := c.Sprite()
	img := &ebiten.Image{}
	if sprite == nil {
		img = graphics.LoadImage("")
	} else {
		img = graphics.LoadImage(c.Sprite().ImagePath)
	}

	op := &ebiten.DrawImageOptions{}

	screenPos := camera.WorldToScreen(c.StateMachine.Position)
	graphics.CameraTransform(op, camera, types.Vector2{X: 1, Y: 1}, screenPos)
	screen.DrawImage(img, op)
}

func (c *Character) Update() {
	// Update animation frame timing first
	c.updateAnimation()

	// Get current frame properties for dynamic values
	frameData := c.AnimationPlayer.ActiveFrameData()
	if frameData == nil {
		panic("Frame properties should not be nil if animation is set")
	}

	switch c.StateMachine.ActiveState {
	case state.StateIdle:
		c.SetAnimation("idle")
	case state.StateWalk | state.StateForward:
		c.SetAnimation("walk")
	case state.StateDash | state.StateForward:
		c.SetAnimation("dash")
	default:
		c.SetAnimation("idle")
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

	sprite := c.Sprite()
	spriteH := 0.0
	if sprite != nil {
		spriteH = sprite.Rect.H
	}
	if c.StateMachine.Position.Y+spriteH >= constants.GroundLevelY {
		c.StateMachine.Position.Y = constants.GroundLevelY - spriteH
		c.StateMachine.Velocity.Y = 0
		c.StateMachine.AddState(state.StateGrounded)
		c.StateMachine.RemoveState(state.StateAirborne)
	}

	if c.StateMachine.Position.X < constants.World.X {
		c.StateMachine.Position.X = constants.World.X
	}
	if c.StateMachine.Position.X+float64(c.Sprite().Rect.W) > constants.World.Right() {
		c.StateMachine.Position.X = constants.World.Right() - float64(c.Sprite().Rect.W)
	}
}
