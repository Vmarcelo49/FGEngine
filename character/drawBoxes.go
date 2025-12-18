package character

import (
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

func createBoxImageOptionsWithCamera(chara *Character, box types.Rect, boxType collision.BoxType, camera *graphics.Camera) *ebiten.DrawImageOptions {
	boxImgOptions := &ebiten.DrawImageOptions{}

	boxImgOptions.GeoM.Scale(box.W, box.H)

	worldPos := chara.Position()
	boxWorldPos := types.Vector2{
		X: worldPos.X + box.X,
		Y: worldPos.Y + box.Y,
	}

	screenPos := camera.WorldToScreen(boxWorldPos)

	graphics.CameraTransform(boxImgOptions, camera, types.Vector2{X: 1, Y: 1}, screenPos)

	if color, exists := boxColors[boxType]; exists {
		boxImgOptions.ColorScale.ScaleWithColor(color)
	}

	return boxImgOptions
}

type BoxDrawable struct {
	Character *Character
}

func (b *BoxDrawable) Draw(screen *ebiten.Image, camera *graphics.Camera) {
	framedata := b.Character.AnimationPlayer.ActiveFrameData()
	if framedata == nil || len(framedata.Boxes) == 0 {
		return
	}

	initWhitePixel()
	boxImgOptions := &ebiten.DrawImageOptions{}
	for boxType, boxes := range framedata.Boxes {
		for _, box := range boxes {
			boxImgOptions = createBoxImageOptionsWithCamera(b.Character, box, boxType, camera)
			screen.DrawImage(whitePixel, boxImgOptions)
		}
	}
}
