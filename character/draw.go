package character

import (
	"fgengine/graphics"

	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Character) Draw(screen *ebiten.Image, camera *graphics.Camera) {
	img := &ebiten.Image{}
	if c.Sprite() == nil {
		img = graphics.LoadImage("") // loads a placeholder image
	} else {
		img = graphics.LoadImage(c.Sprite().ImagePath)
	}

	op := &ebiten.DrawImageOptions{}

	// Obter dados do sprite atual
	sprite := c.Sprite()
	anchorOffset := types.Vector2{X: 32, Y: 128} // info for the placeholder sprite
	if sprite != nil {
		anchorOffset = types.Vector2{X: sprite.AnchorX, Y: sprite.AnchorY}
		if anchorOffset.X == 0 && anchorOffset.Y == 0 {
			anchorOffset = sprite.Anchor
		}
	}

	// Calcular posição na tela (posição real do personagem)
	screenPos := camera.WorldToScreen(c.StateMachine.Position)

	// Aplicar deslocamento para compensar o anchor point
	screenPos.X -= anchorOffset.X
	screenPos.Y -= anchorOffset.Y

	graphics.CameraTransform(op, camera, types.Vector2{X: 1, Y: 1}, screenPos)
	screen.DrawImage(img, op)
}
