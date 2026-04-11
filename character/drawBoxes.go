package character

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/graphics"
	"fgengine/types"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	whitePixel *ebiten.Image
	once       sync.Once
	boxColors  = map[collision.BoxType]color.RGBA{
		collision.Collision: {R: 80, G: 80, B: 80, A: 32},
		collision.Hit:       {R: 100, G: 40, B: 40, A: 32},
		collision.Hurt:      {R: 40, G: 100, B: 40, A: 32},
	}
)

// whitePixel should never be unloaded and only is created once
func initWhitePixel() {
	once.Do(func() {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	})
}

func (c *Character) newBoxOpts(box types.Rect, boxType collision.BoxType, camera *graphics.Camera) *ebiten.DrawImageOptions {
	boxImgOptions := &ebiten.DrawImageOptions{}

	boxImgOptions.GeoM.Scale(box.W, box.H)

	worldPos := c.Position()

	boxWorldPos := types.Vector2{
		X: worldPos.X + box.X,
		Y: worldPos.Y + box.Y,
	}
	// If the character is facing left, it must be flipped horizontally
	if c.StateMachine.Facing == animation.Left {
		boxWorldPos.X = worldPos.X - box.X - box.W
	}

	screenPos := camera.WorldToScreen(boxWorldPos)
	// Aplicar deslocamento para compensar o anchor point
	screenPos.X -= c.Sprite().Anchor.X
	if c.StateMachine.Facing == animation.Left {
		screenPos.X += 2 * (c.Sprite().Anchor.X) // undo the anchor compensation if facing left and compensate in the opposite direction
	}
	screenPos.Y -= c.Sprite().Anchor.Y

	graphics.CameraTransform(boxImgOptions, camera, types.Vector2{X: 1, Y: 1}, screenPos)

	if color, exists := boxColors[boxType]; exists {
		boxImgOptions.ColorScale.ScaleWithColor(color)
	}

	return boxImgOptions
}

func (c *Character) DrawBoxes(screen *ebiten.Image, camera *graphics.Camera) {
	framedata := c.StateMachine.ActiveAnim.ActiveFrameData()
	if framedata == nil || len(framedata.Boxes) == 0 {
		return
	}
	initWhitePixel()
	for boxType, boxes := range framedata.Boxes {
		for _, box := range boxes {
			opts := c.newBoxOpts(box, boxType, camera)
			screen.DrawImage(whitePixel, opts)
		}
	}
}
