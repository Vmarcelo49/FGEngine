// Command generate_fighter_art creates the transparent atlases used by the
// shipped roster. A character with sequence source sheets gets those authored
// frames appended after the 8×4 base pose grid.
package main

import (
	"fmt"
	"image"
	"image/color"
	imagedraw "image/draw"
	"image/png"
	"os"
	"path/filepath"

	xdraw "golang.org/x/image/draw"
)

const tileSize = 128

var sequenceSheets = []string{
	"standing_attacks.png",
	"crouching_attacks.png",
	"aerial_attacks.png",
	"special_defense.png",
}

type point struct{ x, y int }

type palette struct {
	name                     string
	skin, cloth, shade, trim color.RGBA
	hair, glow, metal, eye   color.RGBA
	style                    string
}

var fighters = []palette{
	{name: "brasa", skin: rgba(216, 143, 98), cloth: rgba(145, 38, 34), shade: rgba(91, 29, 36), trim: rgba(255, 190, 83), hair: rgba(237, 73, 34), glow: rgba(255, 111, 43), metal: rgba(255, 218, 154), eye: rgba(255, 245, 225), style: "flame"},
	{name: "muralha", skin: rgba(175, 128, 91), cloth: rgba(55, 73, 88), shade: rgba(30, 42, 54), trim: rgba(216, 167, 86), hair: rgba(48, 48, 57), glow: rgba(242, 171, 68), metal: rgba(161, 184, 196), eye: rgba(255, 239, 207), style: "armor"},
	{name: "iris", skin: rgba(211, 154, 116), cloth: rgba(26, 105, 112), shade: rgba(24, 59, 86), trim: rgba(121, 229, 208), hair: rgba(47, 40, 90), glow: rgba(93, 235, 225), metal: rgba(221, 210, 255), eye: rgba(255, 250, 232), style: "mystic"},
	{name: "vulto", skin: rgba(142, 92, 73), cloth: rgba(35, 32, 53), shade: rgba(20, 21, 38), trim: rgba(146, 92, 211), hair: rgba(16, 17, 27), glow: rgba(187, 123, 255), metal: rgba(198, 194, 226), eye: rgba(225, 250, 255), style: "shadow"},
}

func rgba(r, g, b uint8) color.RGBA { return color.RGBA{R: r, G: g, B: b, A: 255} }

func main() {
	for _, fighter := range fighters {
		atlas := image.NewRGBA(image.Rect(0, 0, tileSize*8, tileSize*4))
		for pose := 0; pose < 32; pose++ {
			tile := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
			drawPose(tile, fighter, pose)
			x, y := (pose%8)*tileSize, (pose/8)*tileSize
			for py := 0; py < tileSize; py++ {
				for px := 0; px < tileSize; px++ {
					atlas.Set(x+px, y+py, tile.At(px, py))
				}
			}
		}
		var err error
		atlas, err = appendSequenceSheets(fighter.name, atlas)
		if err != nil {
			panic(err)
		}
		out := filepath.Join("assets", "characters", fighter.name, "atlas.png")
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			panic(err)
		}
		file, err := os.Create(out)
		if err != nil {
			panic(err)
		}
		if err := png.Encode(file, atlas); err != nil {
			file.Close()
			panic(err)
		}
		if err := file.Close(); err != nil {
			panic(err)
		}
	}
}

func appendSequenceSheets(name string, base *image.RGBA) (*image.RGBA, error) {
	dir := filepath.Join("assets", "characters", name, "sequences")
	paths := make([]string, len(sequenceSheets))
	found := false
	for i, sheet := range sequenceSheets {
		paths[i] = filepath.Join(dir, sheet)
		if _, err := os.Stat(paths[i]); err == nil {
			found = true
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("stat sequence sheet %s: %w", paths[i], err)
		}
	}
	if !found {
		return base, nil
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			return nil, fmt.Errorf("incomplete sequence set for %s: %s: %w", name, path, err)
		}
	}

	cols := 8
	baseFrames := cols * 4
	sequenceFrames := len(paths) * 16
	rows := (baseFrames + sequenceFrames + cols - 1) / cols
	atlas := image.NewRGBA(image.Rect(0, 0, cols*tileSize, rows*tileSize))
	imagedraw.Draw(atlas, base.Bounds(), base, base.Bounds().Min, imagedraw.Src)

	frameIndex := baseFrames
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open sequence sheet %s: %w", path, err)
		}
		sheet, decodeErr := png.Decode(file)
		closeErr := file.Close()
		if decodeErr != nil {
			return nil, fmt.Errorf("decode sequence sheet %s: %w", path, decodeErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("close sequence sheet %s: %w", path, closeErr)
		}
		if sheet.Bounds().Dx() < 4 || sheet.Bounds().Dy() < 4 {
			return nil, fmt.Errorf("sequence sheet %s is too small: %v", path, sheet.Bounds())
		}

		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				srcRect := image.Rect(
					sheet.Bounds().Min.X+col*sheet.Bounds().Dx()/4,
					sheet.Bounds().Min.Y+row*sheet.Bounds().Dy()/4,
					sheet.Bounds().Min.X+(col+1)*sheet.Bounds().Dx()/4,
					sheet.Bounds().Min.Y+(row+1)*sheet.Bounds().Dy()/4,
				)
				x, y := (frameIndex%cols)*tileSize, (frameIndex/cols)*tileSize
				dstRect := image.Rect(x, y, x+tileSize, y+tileSize)
				xdraw.CatmullRom.Scale(atlas, dstRect, sheet, srcRect, xdraw.Src, nil)
				frameIndex++
			}
		}
	}
	return atlas, nil
}

func drawPose(dst *image.RGBA, f palette, pose int) {
	// Body mechanics differ by pose while keeping a consistent side-facing
	// silhouette and the engine's bottom-center anchor at (64, 120).
	// Atlas: 0–1 idle, 2 walk, 3 crouch, 4–7 jump/landing, 8–10 jab,
	// 11–13 heavy, 14 low, 15 launcher, 16–18 special, 19–21 reactions,
	// 22–24 guard, 25 fall, 26 getup, 27 knockdown, 28 KO, 29 victory,
	// 30 backdash, 31 ranged flourish.
	shiftX, shiftY := 0, 0
	lean, crouch, airborne := 0, 0, false
	frontArm, rearArm, legSpread := 0, 0, 0
	if pose >= 2 {
		legSpread = 4
	}
	switch pose {
	case 2: // walk
		legSpread, frontArm, rearArm = 11, 12, -8
	case 3: // crouch
		crouch, legSpread = 18, 13
	case 4, 5, 6, 25: // jump and fall
		airborne, shiftY, lean, legSpread = true, -16, 5, 12
	case 8, 9, 10, 14, 15: // normals
		lean, frontArm, legSpread = 10, 20, 9
	case 11, 12, 13, 16, 17, 18: // heavy strikes and special
		lean, frontArm, rearArm, legSpread = 7, 18, 8, 14
	case 19, 20, 21: // reactions
		shiftX, lean, frontArm, rearArm = -7, -7, -15, -18
	case 22, 23, 24: // guard
		frontArm, rearArm, legSpread = 16, 13, 6
	case 7: // landing
		crouch, legSpread = 8, 13
	case 27, 28: // knockdown / KO
		shiftY, lean = 28, -15
	case 26: // get up
		crouch, lean = 12, -5
	case 29: // win
		frontArm, lean = 24, -3
	case 30: // backdash
		shiftX, lean, frontArm, legSpread = -9, -8, -10, 12
	case 31: // ranged flourish
		lean, frontArm, rearArm = 5, 22, -15
	}
	if pose == 1 {
		frontArm, legSpread = 5, 4
	}

	// Shadow keeps the cutout grounded without adding a background.
	ellipse(dst, 34+shiftX, 108+shiftY, 82, 14, color.RGBA{R: 8, G: 12, B: 20, A: 105})
	if airborne {
		// The jump shadow separates the airborne pose from the ground plane.
		ellipse(dst, 43, 109, 65, 10, color.RGBA{R: 8, G: 12, B: 20, A: 65})
	}

	// Legs, boots, and bright cuffs are layered first so the torso overlaps them.
	leftFoot := point{47 - legSpread + shiftX, 115 + shiftY - crouch/2}
	rightFoot := point{77 + legSpread + shiftX, 115 + shiftY - crouch/2}
	leftKnee := point{48 - legSpread/2 + shiftX, 91 + shiftY - crouch}
	rightKnee := point{67 + legSpread/2 + shiftX, 91 + shiftY - crouch}
	if pose == 9 || pose == 12 || pose == 14 || pose == 15 {
		rightKnee.x += 14
		rightFoot.x += 27
		rightKnee.y -= 8
		rightFoot.y -= 6
	}
	if airborne {
		leftKnee.y -= 4
		rightKnee.y += 8
		leftFoot.y -= 10
		rightFoot.y += 5
	}
	if pose == 27 || pose == 28 {
		poly(dst, []point{{20, 102}, {103, 102}, {109, 116}, {28, 118}}, f.shade)
		ellipse(dst, 39, 97, 22, 18, f.cloth)
		ellipse(dst, 67, 97, 24, 18, f.cloth)
		poly(dst, []point{{35, 107}, {27, 115}, {46, 118}, {52, 111}}, f.metal)
		poly(dst, []point{{68, 108}, {78, 115}, {95, 112}, {86, 106}}, f.metal)
	} else {
		line(dst, leftKnee, leftFoot, 12, f.shade)
		line(dst, rightKnee, rightFoot, 13, f.cloth)
		line(dst, leftKnee, leftFoot, 5, f.cloth)
		line(dst, rightKnee, rightFoot, 5, f.trim)
		ellipse(dst, leftFoot.x-9, leftFoot.y-4, 24, 10, f.metal)
		ellipse(dst, rightFoot.x-7, rightFoot.y-4, 25, 10, f.shade)
		line(dst, point{leftKnee.x, leftKnee.y + 11}, point{leftKnee.x + 8, leftKnee.y + 14}, 3, f.trim)
		line(dst, point{rightKnee.x, rightKnee.y + 11}, point{rightKnee.x + 8, rightKnee.y + 14}, 3, f.trim)
	}

	// Rear arm and torso.
	rearShoulder := point{46 + shiftX, 60 + shiftY - crouch/2}
	rearHand := point{39 + rearArm + shiftX, 79 + shiftY - crouch/2}
	line(dst, rearShoulder, rearHand, 10, f.shade)
	line(dst, rearShoulder, rearHand, 6, f.cloth)
	ellipse(dst, rearHand.x-4, rearHand.y-4, 10, 10, f.skin)

	torso := []point{{42 + shiftX, 51 + shiftY}, {62 + shiftX, 49 + shiftY}, {74 + lean + shiftX, 74 + shiftY - crouch}, {64 + shiftX, 86 + shiftY - crouch}, {43 + shiftX, 75 + shiftY - crouch}}
	poly(dst, torso, f.cloth)
	poly(dst, []point{{44 + shiftX, 54 + shiftY}, {61 + shiftX, 52 + shiftY}, {65 + lean + shiftX, 74 + shiftY - crouch}, {53 + shiftX, 78 + shiftY - crouch}}, f.shade)
	line(dst, point{55 + shiftX, 55 + shiftY}, point{62 + lean/2 + shiftX, 75 + shiftY - crouch}, 3, f.trim)
	line(dst, point{43 + shiftX, 73 + shiftY - crouch}, point{67 + shiftX, 76 + shiftY - crouch}, 5, f.metal)

	// Front arm pose varies strongly so the attack rows read at game scale.
	frontShoulder := point{62 + lean/2 + shiftX, 57 + shiftY - crouch/2}
	frontElbow := point{72 + frontArm/2 + shiftX, 63 + shiftY - crouch/2}
	frontHand := point{79 + frontArm + shiftX, 62 + shiftY - crouch/2}
	if pose == 22 || pose == 23 || pose == 24 {
		frontElbow, frontHand = point{75 + shiftX, 47 + shiftY}, point{70 + shiftX, 39 + shiftY}
	}
	if pose == 29 {
		frontElbow, frontHand = point{71 + shiftX, 29 + shiftY}, point{77 + shiftX, 13 + shiftY}
	}
	line(dst, frontShoulder, frontElbow, 12, f.shade)
	line(dst, frontShoulder, frontElbow, 7, f.cloth)
	line(dst, frontElbow, frontHand, 9, f.cloth)
	line(dst, frontElbow, frontHand, 4, f.trim)
	ellipse(dst, frontHand.x-5, frontHand.y-5, 13, 13, f.skin)
	ellipse(dst, frontHand.x-3, frontHand.y-3, 7, 7, f.metal)

	// Shoulder accents establish each silhouette even in the 640×360 viewport.
	drawStyle(dst, f, shiftX, shiftY, pose, lean)

	// Head and face.
	headX, headY := 53+lean/4+shiftX, 34+shiftY-crouch/3
	ellipse(dst, headX-12, headY-14, 27, 30, f.skin)
	poly(dst, []point{{headX + 2, headY - 3}, {headX + 13, headY + 1}, {headX + 2, headY + 5}}, f.shade)
	ellipse(dst, headX+4, headY-4, 4, 3, f.eye)
	ellipse(dst, headX+5, headY-4, 2, 2, rgba(34, 30, 35))
	line(dst, point{headX + 7, headY + 5}, point{headX + 11, headY + 4}, 2, f.shade)
	drawHair(dst, f, headX, headY)

	// Signature glow. Iris has a clear ranged orb; other fighters show a small
	// accent spark on special frames only.
	if pose == 16 || pose == 17 || pose == 18 || pose == 31 {
		if f.style == "mystic" || pose == 31 {
			orb := point{frontHand.x + 15, frontHand.y - 1}
			ellipse(dst, orb.x-11, orb.y-11, 24, 24, color.RGBA{R: f.glow.R, G: f.glow.G, B: f.glow.B, A: 80})
			ellipse(dst, orb.x-6, orb.y-6, 14, 14, f.glow)
			ellipse(dst, orb.x-3, orb.y-3, 7, 7, rgba(242, 255, 255))
		} else {
			for i := 0; i < 3; i++ {
				poly(dst, []point{{frontHand.x + 4 + i*4, frontHand.y - 1}, {frontHand.x + 14 + i*3, frontHand.y - 9}, {frontHand.x + 10 + i*2, frontHand.y + 2}}, f.glow)
			}
		}
	}
}

func drawStyle(dst *image.RGBA, f palette, dx, dy, pose, lean int) {
	switch f.style {
	case "flame":
		// Ember scarf and flame-shaped hair.
		poly(dst, []point{{39 + dx, 58 + dy}, {27 + dx, 71 + dy}, {42 + dx, 70 + dy}, {35 + dx, 84 + dy}, {53 + dx, 75 + dy}}, f.glow)
		poly(dst, []point{{43 + dx, 27 + dy}, {39 + dx, 16 + dy}, {49 + dx, 22 + dy}, {54 + dx, 9 + dy}, {59 + dx, 20 + dy}, {68 + dx, 14 + dy}, {66 + dx, 30 + dy}}, f.hair)
		line(dst, point{44 + dx, 87 + dy}, point{55 + dx, 88 + dy}, 4, f.trim)
	case "armor":
		// Oversized pauldron, crest and engraved plate.
		ellipse(dst, 37+dx, 50+dy, 27, 25, f.metal)
		poly(dst, []point{{42 + dx, 52 + dy}, {62 + dx, 52 + dy}, {58 + dx, 72 + dy}, {43 + dx, 70 + dy}}, f.cloth)
		poly(dst, []point{{50 + dx, 54 + dy}, {56 + dx, 61 + dy}, {50 + dx, 68 + dy}, {45 + dx, 61 + dy}}, f.trim)
		poly(dst, []point{{45 + dx, 20 + dy}, {57 + dx, 12 + dy}, {69 + dx, 21 + dy}, {66 + dx, 31 + dy}, {48 + dx, 31 + dy}}, f.hair)
		poly(dst, []point{{49 + dx, 19 + dy}, {57 + dx, 3 + dy}, {65 + dx, 20 + dy}}, f.trim)
	case "mystic":
		// Hood, long ribbon and layered robe.
		poly(dst, []point{{42 + dx, 29 + dy}, {49 + dx, 13 + dy}, {66 + dx, 16 + dy}, {73 + dx, 30 + dy}, {64 + dx, 40 + dy}, {45 + dx, 38 + dy}}, f.hair)
		poly(dst, []point{{43 + dx, 72 + dy}, {28 + dx, 94 + dy}, {42 + dx, 93 + dy}, {50 + dx, 79 + dy}}, f.trim)
		poly(dst, []point{{49 + dx, 54 + dy}, {61 + dx, 57 + dy}, {57 + dx, 78 + dy}, {47 + dx, 75 + dy}}, f.glow)
	case "shadow":
		// Mask, high collar and long split scarf.
		poly(dst, []point{{37 + dx, 51 + dy}, {61 + dx, 51 + dy}, {71 + dx, 62 + dy}, {63 + dx, 72 + dy}, {40 + dx, 69 + dy}}, f.shade)
		poly(dst, []point{{46 + dx, 33 + dy}, {65 + dx, 35 + dy}, {66 + dx, 43 + dy}, {46 + dx, 43 + dy}}, f.cloth)
		line(dst, point{46 + dx, 37 + dy}, point{61 + dx, 39 + dy}, 2, f.eye)
		poly(dst, []point{{38 + dx, 57 + dy}, {21 + dx, 69 + dy}, {43 + dx, 67 + dy}}, f.trim)
		poly(dst, []point{{62 + dx, 59 + dy}, {82 + dx, 72 + dy}, {65 + dx, 68 + dy}}, f.trim)
	}

	if pose == 17 && f.style == "armor" {
		ellipse(dst, 82+dx, 49+dy, 22, 22, color.RGBA{R: f.glow.R, G: f.glow.G, B: f.glow.B, A: 130})
	}
}

func drawHair(dst *image.RGBA, f palette, x, y int) {
	if f.style == "armor" || f.style == "mystic" {
		return
	}
	if f.style == "shadow" {
		poly(dst, []point{{x - 13, y - 6}, {x - 9, y - 17}, {x + 5, y - 17}, {x + 14, y - 7}, {x + 10, y + 2}, {x - 11, y + 2}}, f.hair)
		return
	}
	poly(dst, []point{{x - 13, y - 7}, {x - 11, y - 18}, {x - 3, y - 12}, {x + 2, y - 23}, {x + 7, y - 12}, {x + 14, y - 16}, {x + 12, y - 4}, {x + 3, y - 1}}, f.hair)
}

func poly(dst *image.RGBA, points []point, c color.RGBA) {
	if len(points) < 3 {
		return
	}
	minY, maxY := points[0].y, points[0].y
	minX, maxX := points[0].x, points[0].x
	for _, p := range points[1:] {
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
	}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			inside := false
			for i, p := range points {
				q := points[(i+1)%len(points)]
				if (p.y > y) != (q.y > y) && float64(x) < float64(q.x-p.x)*float64(y-p.y)/float64(q.y-p.y)+float64(p.x) {
					inside = !inside
				}
			}
			if inside {
				set(dst, x, y, c)
			}
		}
	}
}

func line(dst *image.RGBA, a, b point, width int, c color.RGBA) {
	dx, dy := b.x-a.x, b.y-a.y
	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}
	if steps == 0 {
		ellipse(dst, a.x-width/2, a.y-width/2, width, width, c)
		return
	}
	for i := 0; i <= steps; i++ {
		x := a.x + dx*i/steps
		y := a.y + dy*i/steps
		ellipse(dst, x-width/2, y-width/2, width, width, c)
	}
}

func ellipse(dst *image.RGBA, x, y, w, h int, c color.RGBA) {
	if w <= 0 || h <= 0 {
		return
	}
	rx, ry := float64(w)/2, float64(h)/2
	cx, cy := float64(x)+rx, float64(y)+ry
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			dx, dy := (float64(px)+.5-cx)/rx, (float64(py)+.5-cy)/ry
			if dx*dx+dy*dy <= 1 {
				set(dst, px, py, c)
			}
		}
	}
}

func set(dst *image.RGBA, x, y int, c color.RGBA) {
	if image.Pt(x, y).In(dst.Bounds()) {
		dst.SetRGBA(x, y, c)
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
