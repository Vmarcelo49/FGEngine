package character

import (
	"crypto/sha256"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestPlayableRosterSpriteRects(t *testing.T) {
	t.Chdir(repoRoot(t))
	for _, fighter := range []string{"Brasa", "Muralha", "Iris", "Vulto"} {
		t.Run(fighter, func(t *testing.T) {
			c, err := LoadCharacter(fighter, 1)
			if err != nil {
				t.Fatal(err)
			}
			images := map[string]image.Image{}
			checked := map[struct {
				path string
				rect image.Rectangle
			}]bool{}
			for name, anim := range c.StateMachine.AnimPlayer.Animations {
				for _, sprite := range anim.Sprites {
					if sprite == nil {
						t.Fatalf("%s contains a nil sprite", name)
					}
					img := images[sprite.ImagePath]
					if img == nil {
						f, err := os.Open(sprite.ImagePath)
						if err != nil {
							t.Fatal(err)
						}
						img, err = png.Decode(f)
						closeErr := f.Close()
						if err != nil {
							t.Fatalf("decode %s: %v", sprite.ImagePath, err)
						}
						if closeErr != nil {
							t.Fatal(closeErr)
						}
						images[sprite.ImagePath] = img
					}
					r := sprite.Rect
					rect := image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))
					if rect.Empty() || !rect.In(img.Bounds()) {
						t.Fatalf("%s sprite rectangle %v outside %v", name, rect, img.Bounds())
					}
					key := struct {
						path string
						rect image.Rectangle
					}{sprite.ImagePath, rect}
					if checked[key] {
						continue
					}
					checked[key] = true
					nonempty := false
					for y := rect.Min.Y; y < rect.Max.Y && !nonempty; y++ {
						for x := rect.Min.X; x < rect.Max.X; x++ {
							_, _, _, alpha := img.At(x, y).RGBA()
							if alpha != 0 {
								nonempty = true
								break
							}
						}
					}
					if !nonempty {
						t.Fatalf("%s sprite %v is fully transparent", name, rect)
					}
				}
			}
		})
	}
}

func TestPlayableRosterDoesNotShareSequenceSheets(t *testing.T) {
	t.Chdir(repoRoot(t))
	seen := map[[sha256.Size]byte]string{}
	for _, fighter := range []string{"brasa", "muralha", "iris", "vulto"} {
		for _, sheet := range []string{"standing_attacks.png", "crouching_attacks.png", "aerial_attacks.png", "special_defense.png"} {
			path := filepath.Join("assets", "characters", fighter, "sequences", sheet)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			sum := sha256.Sum256(data)
			if earlier, ok := seen[sum]; ok {
				t.Errorf("%s duplicates %s", path, earlier)
			}
			seen[sum] = path
		}
	}
}
