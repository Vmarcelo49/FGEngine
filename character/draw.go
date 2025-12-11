package character

import (
	"fgengine/graphics"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Character) Draw(screen *ebiten.Image, camera *graphics.Camera) error {
	img := graphics.LoadImage(c.GetSprite().ImagePath)
	op := &ebiten.DrawImageOptions{}

	screenPos := camera.WorldToScreen(c.StateMachine.Position)
	graphics.ApplyCameraTransformNEW(op, camera, types.Vector2{X: 1, Y: 1}, screenPos)
	screen.DrawImage(img, op)
	return nil
}
