package graphics

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/types"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawCharacter(screen *ebiten.Image, c *character.Character, camera *Camera) {
	sprite := c.Sprite()
	img := SpriteImage(sprite)

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

		CameraTransform(op, camera, types.Vector2{X: 1, Y: 1}, screenPos)
		screen.DrawImage(img, op)

	}
}

// SpriteImage returns either the full source image or the authored atlas tile
// for a character sprite. Character art files can therefore share a compact
// atlas without changing the simulation data contract.
func SpriteImage(sprite *animation.Sprite) *ebiten.Image {
	if sprite == nil {
		return LoadImage("")
	}
	img := LoadImage(sprite.ImagePath)
	if sprite.Rect.W <= 0 || sprite.Rect.H <= 0 {
		return img
	}
	rect := image.Rect(
		int(sprite.Rect.X),
		int(sprite.Rect.Y),
		int(sprite.Rect.X+sprite.Rect.W),
		int(sprite.Rect.Y+sprite.Rect.H),
	)
	if rect.Empty() || !rect.In(img.Bounds()) {
		return img
	}
	return img.SubImage(rect).(*ebiten.Image)
}
