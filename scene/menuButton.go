package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type Button struct {
	Background *ebiten.Image
	Face       font.Face
	Label      string
	X, Y       int
	Hovered    bool
	Pressed    bool
}

func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.X), float64(b.Y))

	// Escurece no hover/pressed
	if b.Hovered {
		op.ColorScale.Scale(0.92, 0.92, 0.92, 1)
	}
	if b.Pressed {
		op.ColorScale.Scale(0.85, 0.85, 0.85, 1)
	}

	// Sombra simples (desenho duplicado, deslocado e preto translúcido)
	shadow := &ebiten.DrawImageOptions{}
	shadow.GeoM.Translate(float64(b.X+2), float64(b.Y+2))
	shadow.ColorScale.Scale(0, 0, 0, 0.35)
	screen.DrawImage(b.Background, shadow)

	// Fundo
	screen.DrawImage(b.Background, op)

	// Texto
	//paddingX, paddingY := 12, 10
	//opts := text.DrawOptions{}
	//text.Draw(screen, b.Label, face, opts)
}
