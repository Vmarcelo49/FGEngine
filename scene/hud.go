package scene

import (
	"fgengine/character"
	"image/color"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HUD layout in screen space (viewport 640x360). All coordinates are
// constants below; the bars scale by HP fraction.
const (
	hudBarY   = 10.0
	hudBarH   = 14.0
	hudBarW   = 250.0
	hudMargin = 16.0

	hudPipY = 28.0
	hudPipW = 10.0
	hudPipH = 6.0

	hudTimerY = 8
	hudFontW  = 7 // approx advance of the debug font, for centering
)

var (
	hudBorder = color.RGBA{R: 240, G: 240, B: 240, A: 255}
	hudEmpty  = color.RGBA{R: 120, G: 20, B: 20, A: 255}
	hudFull   = color.RGBA{R: 250, G: 200, B: 40, A: 255}
	hudLow    = color.RGBA{R: 220, G: 40, B: 40, A: 255}
	hudPipOff = color.RGBA{R: 60, G: 60, B: 60, A: 255}
)

// drawHUD renders HP bars, round pips, the timer, and round-end text in
// screen space, last (SPEC §7.8). Presentation only: reads snapshot state
// through gameplay.HUDState, never mutates it.
func (g *GameplayScene) drawHUD(screen *ebiten.Image) {
	state := g.gamestate.HUD()
	p1Accent := accentColor(g.gamestate.Characters[0].Properties.AccentColor)
	p2Accent := accentColor(g.gamestate.Characters[1].Properties.AccentColor)

	drawHPBar(screen, hudMargin, state.P1HP, state.P1MaxHP, false, p1Accent)
	drawHPBar(screen, 640-hudMargin-hudBarW, state.P2HP, state.P2MaxHP, true, p2Accent)
	drawNameplate(screen, g.gamestate.Characters[0], int(hudMargin), 34, false)
	drawNameplate(screen, g.gamestate.Characters[1], int(640-hudMargin), 34, true)

	drawPips(screen, 288, state.P1Wins, false)
	drawPips(screen, 352, state.P2Wins, true)

	timer := strconv.Itoa(state.TimerSeconds)
	ebitenutil.DebugPrintAt(screen, timer, 320-len(timer)*hudFontW/2, hudTimerY)

	if text := g.tr(state.CenterTextKey); text != "" {
		ebitenutil.DebugPrintAt(screen, text, 320-len(text)*hudFontW/2, 170)
	}
}

// drawHPBar draws one mirrored HP bar. If mirror is set the fill drains
// from the right edge (P2); otherwise from the left (P1).
func drawHPBar(screen *ebiten.Image, x float64, hp, maxHP int, mirror bool, accent color.RGBA) {
	frac := 0.0
	if maxHP > 0 && hp > 0 {
		frac = float64(hp) / float64(maxHP)
		if frac > 1 {
			frac = 1
		}
	}
	fill := accent
	if frac < 0.3 {
		fill = hudLow
	}
	fw := float32(hudBarW * frac)
	vector.DrawFilledRect(screen, float32(x), hudBarY, hudBarW, hudBarH, hudEmpty, false)
	if fw > 0 {
		fx := float32(x)
		if mirror {
			fx = float32(x+hudBarW) - fw
		}
		vector.DrawFilledRect(screen, fx, hudBarY, fw, hudBarH, fill, false)
		vector.DrawFilledRect(screen, fx, hudBarY, fw, 2, color.RGBA{R: 255, G: 255, B: 255, A: 70}, false)
	}
	vector.StrokeRect(screen, float32(x), hudBarY, hudBarW, hudBarH, 1, hudBorder, false)
}

func drawNameplate(screen *ebiten.Image, fighter *character.Character, edge, y int, rightAligned bool) {
	if fighter == nil {
		return
	}
	name := strings.ToUpper(fighter.Name)
	x := edge
	if rightAligned {
		x -= len(name) * hudFontW
	}
	vector.DrawFilledRect(screen, float32(x-3), float32(y-3), float32(len(name)*hudFontW+6), 16, color.RGBA{R: 12, G: 18, B: 31, A: 225}, false)
	ebitenutil.DebugPrintAt(screen, name, x, y)
}

// drawPips draws round-win pips: left-aligned (P1) or right-aligned (P2).
func drawPips(screen *ebiten.Image, edge float64, wins int, rightAlign bool) {
	for i := 0; i < 2; i++ {
		x := edge + float64(i)*(hudPipW+4)
		if rightAlign {
			x = edge - float64(i+1)*(hudPipW+4) + 4
		}
		clr := hudPipOff
		if i < wins {
			clr = hudFull
		}
		vector.DrawFilledRect(screen, float32(x), hudPipY, hudPipW, hudPipH, clr, false)
	}
}

// tr resolves a HUD text key with key-fallback (never crashes on a missing
// locale entry; the raw key showing is the signal).
func (g *GameplayScene) tr(key string) string {
	if key == "" {
		return ""
	}
	if s, ok := g.texts[key]; ok {
		return s
	}
	return key
}
