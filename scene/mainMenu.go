package scene

import (
	"fgengine/input"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type menuOption struct {
	label  string
	signal SceneStatus
}

var mainMenuOptions = []menuOption{
	{label: "ARCADE  /  VS CPU", signal: SceneSelectCPU},
	{label: "LOCAL VERSUS", signal: SceneSelectVersus},
	{label: "CONTROLLERS", signal: SceneControllerSelect},
	{label: "QUIT", signal: SceneQuit},
}

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
		m.selected = (m.selected + 1) % len(mainMenuOptions)
	}
	if input.JustPressed(p1, prev, input.Up) {
		m.selected--
		if m.selected < 0 {
			m.selected = len(mainMenuOptions) - 1
		}
	}
	if input.JustPressed(p1, prev, input.A) {
		return mainMenuOptions[m.selected].signal
	}
	return SceneDontChange
}

func (m *MainMenuScene) Draw(screen *ebiten.Image) {
	drawMenuBackdrop(screen)
	ebitenutil.DebugPrintAt(screen, "FG  /  FIGHTING GROUNDS", 34, 31)
	ebitenutil.DebugPrintAt(screen, "A DATA-DRIVEN FIGHTING GAME", 36, 52)

	for i, opt := range mainMenuOptions {
		x, y := float32(34), float32(112+i*43)
		bg := color.RGBA{R: 19, G: 28, B: 49, A: 245}
		border := color.RGBA{R: 52, G: 68, B: 94, A: 255}
		if i == m.selected {
			bg = color.RGBA{R: 35, G: 73, B: 96, A: 255}
			border = color.RGBA{R: 87, G: 220, B: 214, A: 255}
		}
		vector.DrawFilledRect(screen, x, y, 222, 34, bg, false)
		vector.StrokeRect(screen, x, y, 222, 34, 1, border, false)
		if i == m.selected {
			vector.DrawFilledRect(screen, x, y, 4, 34, border, false)
		}
		ebitenutil.DebugPrintAt(screen, opt.label, int(x+16), int(y+13))
	}

	ebitenutil.DebugPrintAt(screen, "MOVE  W/S OR ARROWS", 34, 310)
	ebitenutil.DebugPrintAt(screen, "CONFIRM  U / GAMEPAD A", 34, 326)
	ebitenutil.DebugPrintAt(screen, "ESC  QUIT", 34, 342)
}

func drawMenuBackdrop(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 9, G: 14, B: 29, A: 255})
	vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{R: 9, G: 14, B: 29, A: 255}, false)
	vector.DrawFilledRect(screen, 340, 0, 300, 360, color.RGBA{R: 12, G: 27, B: 44, A: 255}, false)

	for i := 0; i < 6; i++ {
		y := float32(220 + i*23)
		vector.StrokeLine(screen, 280, y, 640, y-24, 1, color.RGBA{R: 28, G: 55, B: 73, A: 255}, false)
	}
	for i := 0; i < 7; i++ {
		x := float32(365 + i*46)
		vector.StrokeLine(screen, x, 0, x-90, 360, 1, color.RGBA{R: 22, G: 43, B: 62, A: 255}, false)
	}

	vector.DrawFilledRect(screen, 302, 0, 2, 360, color.RGBA{R: 246, G: 120, B: 65, A: 255}, false)
	vector.DrawFilledRect(screen, 306, 0, 1, 360, color.RGBA{R: 247, G: 190, B: 96, A: 200}, false)
	vector.StrokeLine(screen, 26, 91, 614, 91, 1, color.RGBA{R: 61, G: 80, B: 101, A: 255}, false)
}
