package scene

import (
	"fgengine/input"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var menuOptions = []string{"Play", "Options", "Exit"}

type MainMenuScene struct {
	selected   int
	prevInputs [2]input.GameInput
}

func MakeMainMenuScene() Scene {
	return &MainMenuScene{}
}

func (m *MainMenuScene) Update(inputs [2]input.GameInput) SceneStatus {
	p1 := inputs[0]
	prev := m.prevInputs[0]
	defer func() { m.prevInputs = inputs }()

	if input.JustPressed(p1, prev, input.Down) {
		m.selected++
		if m.selected >= len(menuOptions) {
			m.selected = 0
		}
	}
	if input.JustPressed(p1, prev, input.Up) {
		m.selected--
		if m.selected < 0 {
			m.selected = len(menuOptions) - 1
		}
	}
	if input.JustPressed(p1, prev, input.A) {
		switch m.selected {
		case 0: // Play
			return Scene2
		case 1: // Options
			return Scene1
		case 2: // Exit
			return SceneDontChange
		}
	}
	return SceneDontChange
}

func (m *MainMenuScene) Draw(screen *ebiten.Image) {
	titleX := 75
	ebitenutil.DebugPrintAt(screen, "MAIN MENU", int(titleX), 40)

	boxW := float32(120.0)
	boxH := float32(30.0)
	startY := float32(100.0)
	spacing := float32(40.0)

	for i, opt := range menuOptions {
		x := float32(40)
		y := float32(startY + float32(i)*spacing)

		bgColor := color.RGBA{R: 60, G: 60, B: 60, A: 255}
		if i == m.selected {
			bgColor = color.RGBA{R: 100, G: 149, B: 237, A: 255}
		}
		vector.FillRect(screen, x, y, boxW, boxH, bgColor, false)

		textX := int(x + boxW/2 - float32(len(opt)*3))
		textY := int(y + boxH/2 - 6)
		ebitenutil.DebugPrintAt(screen, opt, textX, textY)
	}
}
