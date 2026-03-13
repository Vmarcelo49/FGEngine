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

	screenPos := camera.WorldToScreen(c.StateMachine.Position)
	graphics.CameraTransform(op, camera, types.Vector2{X: 1, Y: 1}, screenPos)
	screen.DrawImage(img, op)
}
