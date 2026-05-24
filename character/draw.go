package character

import (
	"fgengine/animation"
	"fgengine/graphics"

	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (c *Character) Draw(screen *ebiten.Image, camera *graphics.Camera) {
	var img *ebiten.Image
	sprite := c.Sprite()
	if sprite == nil {
		img = graphics.LoadImage("") // loads a placeholder image
	} else {
		img = graphics.LoadImage(sprite.ImagePath)
	}

	op := &ebiten.DrawImageOptions{}

	anchorOffset := types.Vector2{X: 25, Y: 100} // info for the placeholder sprite
	if sprite != nil {
		anchorOffset = types.Vector2{X: sprite.Anchor.X, Y: sprite.Anchor.Y}
		if anchorOffset.X == 0 && anchorOffset.Y == 0 {
			anchorOffset = sprite.Anchor
		}
	}

	// Calcular posição na tela (posição real do personagem)
	var screenPos types.Vector2
	state := c.StateMachine
	animName := "none"

	if state != nil {
		if camera != nil {
			screenPos = camera.WorldToScreen(state.Position)
		}

		// Aplicar deslocamento para compensar o anchor point
		screenPos.X -= anchorOffset.X
		screenPos.Y -= anchorOffset.Y

		// Handle horizontal flip when facing left
		if state.IsFacingLeft == animation.Left {
			op.GeoM.Scale(-1, 1) // Flip horizontal
			// Compensate for the flip by translating by 2x the anchor point
			op.GeoM.Translate(2*anchorOffset.X, 0)
		}

		graphics.CameraTransform(op, camera, types.Vector2{X: 1, Y: 1}, screenPos)
		screen.DrawImage(img, op)

		// Debug info on top of the character
		animName = state.ActiveAnim.ActiveAnimationName()
	}
	ebitenutil.DebugPrintAt(screen, animName, int(screenPos.X)+int(anchorOffset.X), int(screenPos.Y))
}
