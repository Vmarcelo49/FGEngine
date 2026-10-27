package scene

import (
	"fgengine/character"
	"fgengine/graphics"
	"fgengine/input"
	"image/color"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type rosterEntry struct {
	file string
	char *character.Character
}

type CharacterSelectScene struct {
	mode        MatchMode
	roster      []rosterEntry
	cursor      [2]int
	ready       [2]bool
	prevInputs  [2]input.GameInput
	cpuSelectP2 bool
}

func MakeCharacterSelectScene(mode MatchMode) *CharacterSelectScene {
	scene := &CharacterSelectScene{mode: mode}
	paths, err := filepath.Glob("assets/characters/*.toml")
	if err != nil {
		log.Printf("character roster scan failed: %v", err)
		return scene
	}
	for _, path := range paths {
		file := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if strings.EqualFold(file, "PlaceHolder") {
			continue
		}
		ch, err := character.LoadCharacter(file, 1)
		if err != nil {
			log.Printf("character roster skipped %q: %v", file, err)
			continue
		}
		scene.roster = append(scene.roster, rosterEntry{file: file, char: ch})
	}
	if len(scene.roster) > 1 {
		scene.cursor[1] = 1
	}
	return scene
}

func (s *CharacterSelectScene) Update(inputs [2]input.GameInput) SceneStatus {
	if len(s.roster) == 0 {
		previous := s.prevInputs[0]
		s.prevInputs = inputs
		if input.JustPressed(inputs[0], previous, input.B) {
			return SceneMainMenu
		}
		return SceneDontChange
	}

	if s.mode == MatchVersusCPU {
		current := inputs[0]
		previous := s.prevInputs[0]
		if input.JustPressed(current, previous, input.B) {
			if s.cpuSelectP2 {
				s.cpuSelectP2 = false
				s.ready[0] = false
			} else {
				s.prevInputs = inputs
				return SceneMainMenu
			}
		}
		cursorPlayer := 0
		if s.cpuSelectP2 {
			cursorPlayer = 1
		}
		s.moveCursor(cursorPlayer, current, previous)
		if input.JustPressed(current, previous, input.A) {
			if !s.cpuSelectP2 {
				s.ready[0] = true
				s.cpuSelectP2 = true
				if len(s.roster) > 1 && s.cursor[1] == s.cursor[0] {
					s.cursor[1] = (s.cursor[0] + 1) % len(s.roster)
				}
			} else {
				s.prevInputs = inputs
				return SceneStartMatch
			}
		}
		s.prevInputs = inputs
		return SceneDontChange
	}

	for player := range 2 {
		current := inputs[player]
		previous := s.prevInputs[player]
		if input.JustPressed(current, previous, input.B) {
			if s.ready[player] {
				s.ready[player] = false
			} else if player == 0 {
				s.prevInputs = inputs
				return SceneMainMenu
			}
		}
		if !s.ready[player] {
			s.moveCursor(player, current, previous)
			if input.JustPressed(current, previous, input.A) {
				s.ready[player] = true
			}
		}
	}
	s.prevInputs = inputs
	if s.ready[0] && s.ready[1] {
		return SceneStartMatch
	}
	return SceneDontChange
}

func (s *CharacterSelectScene) moveCursor(player int, current, previous input.GameInput) {
	index := s.cursor[player]
	column, row := index%2, index/2
	if input.JustPressed(current, previous, input.Left) {
		column = (column + 1) % 2
	}
	if input.JustPressed(current, previous, input.Right) {
		column = (column + 1) % 2
	}
	if input.JustPressed(current, previous, input.Up) || input.JustPressed(current, previous, input.Down) {
		row = (row + 1) % ((len(s.roster) + 1) / 2)
	}
	index = row*2 + column
	if index >= len(s.roster) {
		index = len(s.roster) - 1
	}
	s.cursor[player] = index
}

func (s *CharacterSelectScene) playerOneName() string {
	return s.roster[s.cursor[0]].file
}

func (s *CharacterSelectScene) playerTwoName() string {
	return s.roster[s.cursor[1]].file
}

func (s *CharacterSelectScene) Draw(screen *ebiten.Image) {
	drawMenuBackdrop(screen)
	vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{R: 8, G: 13, B: 25, A: 226}, false)
	ebitenutil.DebugPrintAt(screen, "CHOOSE YOUR FIGHTER", 24, 24)
	if s.mode == MatchVersusCPU {
		if s.cpuSelectP2 {
			ebitenutil.DebugPrintAt(screen, "PLAYER 1: PICK THE CPU OPPONENT", 24, 45)
		} else {
			ebitenutil.DebugPrintAt(screen, "PLAYER 1: PICK YOUR FIGHTER", 24, 45)
		}
	} else {
		ebitenutil.DebugPrintAt(screen, "LOCAL VERSUS  /  BOTH PLAYERS PICK", 24, 45)
	}

	for index, entry := range s.roster {
		column, row := index%2, index/2
		x := float32(24 + column*300)
		y := float32(76 + row*130)
		accent := accentColor(entry.char.Properties.AccentColor)
		selectedByP1 := s.cursor[0] == index
		selectedByP2 := s.mode == MatchLocalVersus && s.cursor[1] == index
		selectedByCPU := s.mode == MatchVersusCPU && s.cpuSelectP2 && s.cursor[1] == index

		bg := color.RGBA{R: 18, G: 26, B: 43, A: 245}
		vector.DrawFilledRect(screen, x, y, 282, 119, bg, false)
		vector.StrokeRect(screen, x, y, 282, 119, 1, color.RGBA{R: 44, G: 57, B: 81, A: 255}, false)
		if selectedByP1 {
			vector.StrokeRect(screen, x-2, y-2, 286, 123, 1, color.RGBA{R: 84, G: 219, B: 211, A: 255}, false)
		}
		if selectedByP2 || selectedByCPU {
			vector.StrokeRect(screen, x+2, y+2, 278, 115, 1, color.RGBA{R: 251, G: 137, B: 85, A: 255}, false)
		}
		vector.DrawFilledRect(screen, x, y, 4, 119, accent, false)

		if entry.char.Sprite() != nil {
			img := graphics.SpriteImage(entry.char.Sprite())
			if img != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(0.78, 0.78)
				op.GeoM.Translate(float64(x+4), float64(y+15))
				screen.DrawImage(img, op)
			}
		}

		textX := int(x + 105)
		ebitenutil.DebugPrintAt(screen, strings.ToUpper(entry.char.Name), textX, int(y+16))
		ebitenutil.DebugPrintAt(screen, entry.char.Properties.Role, textX, int(y+33))
		ebitenutil.DebugPrintAt(screen, entry.char.Properties.Tagline, textX, int(y+51))
		ebitenutil.DebugPrintAt(screen, "SPECIAL  "+entry.char.Properties.SpecialName, textX, int(y+76))
		ebitenutil.DebugPrintAt(screen, entry.char.Properties.SpecialInput, textX, int(y+92))

		if selectedByP1 {
			label := "P1"
			if s.mode == MatchVersusCPU && s.ready[0] {
				label = "P1 READY"
			} else if s.mode == MatchLocalVersus && s.ready[0] {
				label = "P1 READY"
			}
			drawTag(screen, x+7, y+4, label, color.RGBA{R: 84, G: 219, B: 211, A: 255})
		}
		if selectedByP2 || selectedByCPU {
			label := "P2"
			if selectedByCPU {
				label = "CPU"
			} else if s.ready[1] {
				label = "P2 READY"
			}
			drawTag(screen, x+220, y+4, label, color.RGBA{R: 251, G: 137, B: 85, A: 255})
		}
	}

	if len(s.roster) == 0 {
		ebitenutil.DebugPrintAt(screen, "NO PLAYABLE CHARACTERS FOUND", 24, 180)
	}
	ebitenutil.DebugPrintAt(screen, "P1 WASD U/I/O/J  P2 ARROWS N/M/,/.  PAD D-PAD + 4 BUTTONS", 24, 345)
}

func drawTag(screen *ebiten.Image, x, y float32, label string, accent color.RGBA) {
	w := float32(len(label)*7 + 8)
	vector.DrawFilledRect(screen, x, y, w, 14, color.RGBA{R: accent.R / 4, G: accent.G / 4, B: accent.B / 4, A: 255}, false)
	vector.StrokeRect(screen, x, y, w, 14, 1, accent, false)
	ebitenutil.DebugPrintAt(screen, label, int(x+4), int(y+3))
}

func accentColor(value string) color.RGBA {
	value = strings.TrimPrefix(value, "#")
	if len(value) != 6 {
		return color.RGBA{R: 93, G: 207, B: 206, A: 255}
	}
	parsed, err := strconv.ParseUint(value, 16, 32)
	if err != nil {
		return color.RGBA{R: 93, G: 207, B: 206, A: 255}
	}
	return color.RGBA{R: uint8(parsed >> 16), G: uint8(parsed >> 8), B: uint8(parsed), A: 255}
}
