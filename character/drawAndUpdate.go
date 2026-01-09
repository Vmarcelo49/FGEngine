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
	sm := c.StateMachine
	if sm.HitstunFrames > 0 {
		sm.HitstunFrames--
		if sm.HitstunFrames == 0 {
			sm.ClearState(state.StateOnHitsun)
		}
	}

	// attack queue
	if sm.AttackRequested && !sm.IsAttacking() {
		sm.AddState(state.StateAttack)
		c.AttackHasHit = false
		c.AnimationPlayer.AnimationQueue = []string{"idle"}
		c.ensureAnimation("attack", false)
	}
	sm.AttackRequested = false

	movementSpeed := 3.2
	dashSpeed := 6.0
	speed := movementSpeed

	if sm.DashRequested {
		sm.AddState(state.StateDash)
		speed = dashSpeed
	} else {
		sm.ClearState(state.StateDash)
	}

	if sm.MoveInput != 0 {
		sm.ClearState(state.StateNeutral | state.StateBackward | state.StateForward)
		sm.AddState(state.StateWalk)
		if sm.MoveInput > 0 {
			if sm.CharacterOrientation == state.Right {
				sm.AddState(state.StateForward)
			} else {
				sm.AddState(state.StateBackward)
			}
		} else {
			if sm.CharacterOrientation == state.Left {
				sm.AddState(state.StateForward)
			} else {
				sm.AddState(state.StateBackward)
			}
		}
		sm.Velocity.X = float64(sm.MoveInput) * speed
	} else {
		sm.ClearState(state.StateWalk | state.StateDash | state.StateBackward | state.StateForward)
		sm.AddState(state.StateNeutral)
		if sm.Velocity.X > 0 {
			sm.Velocity.X -= c.Friction
			if sm.Velocity.X < 0 {
				sm.Velocity.X = 0
			}
		} else if sm.Velocity.X < 0 {
			sm.Velocity.X += c.Friction
			if sm.Velocity.X > 0 {
				sm.Velocity.X = 0
			}
		}
	}
	sm.DashRequested = false

	if sm.JumpRequested && sm.HasState(state.StateGrounded) {
		sm.ClearState(state.StateGrounded | state.StateNeutral | state.StateCrouching)
		sm.AddState(state.StateAirborne | state.StateJump)
		sm.Velocity.Y = -c.JumpHeight
	}
	sm.JumpRequested = false

	// Apply gravity only if airborne and not ignoring gravity frames
	if sm.HasState(state.StateAirborne) && sm.IgnoreGravityFrames <= 0 {
		sm.Velocity.Y += constants.Gravity
	} else {
		sm.IgnoreGravityFrames--
	}

	// Advance animation after state decisions
	c.updateAnimation()

	frameData := c.AnimationPlayer.ActiveFrameData()
	if frameData != nil {
		if frameData.ChangeXSpeed != 0 {
			sm.Velocity.X = float64(frameData.ChangeXSpeed)
		}
		if frameData.ChangeYSpeed != 0 {
			sm.Velocity.Y = float64(frameData.ChangeYSpeed)
		}
	}

	sm.Position.X += sm.Velocity.X
	sm.Position.Y += sm.Velocity.Y

	if sm.HasState(state.StateAirborne) {
		if sm.Velocity.Y > 0 {
			sm.AddState(state.StateFalling)
			sm.ClearState(state.StateJump)
		} else {
			sm.ClearState(state.StateFalling)
		}
	}

	sprite := c.Sprite()
	spriteH := 0.0
	spriteW := 0.0
	if sprite != nil {
		spriteH = sprite.Rect.H
		spriteW = sprite.Rect.W
	}
	if sm.Position.Y+spriteH >= constants.GroundLevelY {
		sm.Position.Y = constants.GroundLevelY - spriteH
		sm.Velocity.Y = 0
		sm.ClearState(state.StateAirborne | state.StateJump | state.StateFalling)
		sm.AddState(state.StateGrounded | state.StateNeutral)
	}

	if sm.Position.X < constants.World.X {
		sm.Position.X = constants.World.X
	}
	if sm.Position.X+spriteW > constants.World.Right() {
		sm.Position.X = constants.World.Right() - spriteW
	}

	// Animation selection based on resulting state
	switch {
	case sm.HasState(state.StateOnHitsun):
		c.ensureAnimation("fall", true)
	case sm.IsAttacking():
		c.ensureAnimation("attack", false)
		// return to idle when attack animation finishes
		if c.AnimationPlayer.IsFinished() {
			sm.ClearState(state.StateAttack | state.StateSpecialAttack | state.StateSuperAttack | state.StateA | state.StateB | state.StateC)
			c.AttackHasHit = false
			c.ensureAnimation("idle", true)
		}
	case sm.HasState(state.StateAirborne):
		if sm.HasState(state.StateFalling) {
			c.ensureAnimation("fall", true)
		} else {
			c.ensureAnimation("jump", true)
		}
	case sm.HasState(state.StateDash):
		c.ensureAnimation("dash", true)
	case sm.HasState(state.StateWalk):
		c.ensureAnimation("walk", true)
	default:
		c.ensureAnimation("idle", true)
	}
}
