package editor

import (
	"fgengine/animation"
	"fgengine/graphics"
	"fgengine/types"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	editorWhitePixel *ebiten.Image
	editorWhiteOnce  sync.Once
	editorBoxColors  = map[types.BoxType]color.RGBA{
		types.Collision: {R: 80, G: 80, B: 80, A: 32},
		types.Hit:       {R: 100, G: 40, B: 40, A: 32},
		types.Hurt:      {R: 40, G: 100, B: 40, A: 32},
	}
)

func initEditorWhitePixel() {
	editorWhiteOnce.Do(func() {
		editorWhitePixel = ebiten.NewImage(1, 1)
		editorWhitePixel.Fill(color.White)
	})
}

func (ed *CharacterEditor) drawCharacterPreview(screen *ebiten.Image) {
	if ed.char == nil {
		return
	}

	sprite := ed.char.Sprite()
	imgPath := ""
	if sprite != nil {
		imgPath = sprite.ImagePath
	}

	img := graphics.LoadImage(imgPath)
	if img == nil {
		return
	}

	imgW := img.Bounds().Dx()
	imgH := img.Bounds().Dy()
	scale := float64(ed.previewScale)
	if scale <= 0 {
		scale = 1
	}

	anchorX := float64(imgW) / 2
	anchorY := float64(imgH) / 2
	if sprite != nil && (sprite.Anchor.X != 0 || sprite.Anchor.Y != 0) {
		anchorX = sprite.Anchor.X
		anchorY = sprite.Anchor.Y
	}

	screenW := screen.Bounds().Dx()
	screenH := screen.Bounds().Dy()
	centerX := float64(screenW) / 2
	centerY := float64(screenH) / 2

	op := &ebiten.DrawImageOptions{}
	topLeftY := centerY - anchorY*scale
	if ed.char.StateMachine != nil && ed.char.StateMachine.Facing == animation.Left {
		topLeftX := centerX - (float64(imgW)-anchorX)*scale
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(float64(imgW)*scale+topLeftX, topLeftY)
	} else {
		topLeftX := centerX - anchorX*scale
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(topLeftX, topLeftY)
	}
	screen.DrawImage(img, op)

	ed.drawCharacterPreviewBoxes(screen, centerX, centerY, anchorX, anchorY, scale)
	ed.drawCharacterPreviewAnchor(screen, centerX, centerY, scale)
}

func (ed *CharacterEditor) drawCharacterPreviewBoxes(screen *ebiten.Image, centerX, centerY, anchorX, anchorY, scale float64) {
	frameData := ed.currentFrameData()
	if frameData == nil || len(frameData.Boxes) == 0 {
		return
	}

	initEditorWhitePixel()

	facingLeft := ed.char != nil && ed.char.StateMachine != nil && ed.char.StateMachine.Facing == animation.Left

	for boxType, boxes := range frameData.Boxes {
		for _, box := range boxes {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(box.W*scale, box.H*scale)

			boxX := centerX - anchorX*scale + box.X*scale
			if facingLeft {
				boxX = centerX + anchorX*scale - (box.X+box.W)*scale
			}
			boxY := centerY - anchorY*scale + box.Y*scale

			opts.GeoM.Translate(boxX, boxY)

			if col, ok := editorBoxColors[boxType]; ok {
				opts.ColorScale.ScaleWithColor(col)
			}

			screen.DrawImage(editorWhitePixel, opts)
		}
	}
}

func (ed *CharacterEditor) drawCharacterPreviewAnchor(screen *ebiten.Image, centerX, centerY, scale float64) {
	initEditorWhitePixel()

	lineLen := 10.0 * scale
	thickness := 2.0

	horizontal := &ebiten.DrawImageOptions{}
	horizontal.GeoM.Scale(2*lineLen, thickness)
	horizontal.GeoM.Translate(centerX-lineLen, centerY-thickness/2)
	horizontal.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	screen.DrawImage(editorWhitePixel, horizontal)

	vertical := &ebiten.DrawImageOptions{}
	vertical.GeoM.Scale(thickness, 2*lineLen)
	vertical.GeoM.Translate(centerX-thickness/2, centerY-lineLen)
	vertical.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	screen.DrawImage(editorWhitePixel, vertical)
}
