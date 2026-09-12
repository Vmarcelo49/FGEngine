package character

import (
	"fgengine/animation"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempImg(t *testing.T) string {
	t.Helper()
	img := filepath.Join(t.TempDir(), "sprite.png")
	if err := os.WriteFile(img, []byte("fakepng"), 0644); err != nil {
		t.Fatal(err)
	}
	return img
}

func minAnim(img string, withLoop bool) *animation.Animation {
	anim := &animation.Animation{
		Sprites:   []*animation.Sprite{{ImagePath: img}},
		FrameData: []animation.FrameData{{Duration: 2}},
	}
	if withLoop {
		anim.LoopFrames = &animation.LoopFrame{Start: 0, End: 0}
	}
	return anim
}

// minCharacter builds a fully valid character: full required set, no specials.
func minCharacter(t *testing.T) *Character {
	t.Helper()
	img := writeTempImg(t)
	anims := make(map[string]*animation.Animation, len(RequiredStates))
	for _, name := range RequiredStates {
		anims[name] = minAnim(img, true)
	}
	return &Character{
		Name:       "Test",
		Properties: CharacterProperties{MaxHP: 10000},
		StateMachine: &animation.StateMachine{
			AnimPlayer: &animation.AnimationPlayer{Animations: anims},
		},
	}
}

func TestValidateAcceptsMinimalCharacter(t *testing.T) {
	c := minCharacter(t)
	if errs := c.Validate(); len(errs) > 0 {
		t.Fatalf("valid character rejected: %v", errs)
	}
}

func TestValidateRejects(t *testing.T) {
	cases := map[string]struct {
		mutate  func(t *testing.T, c *Character)
		wantSub string
	}{
		"empty name": {
			mutate:  func(t *testing.T, c *Character) { c.Name = "" },
			wantSub: "name",
		},
		"zero maxHP": {
			mutate:  func(t *testing.T, c *Character) { c.Properties.MaxHP = 0 },
			wantSub: "maxHP",
		},
		"nil state machine": {
			mutate:  func(t *testing.T, c *Character) { c.StateMachine = nil },
			wantSub: "activeAnim",
		},
		"missing required state": {
			mutate: func(t *testing.T, c *Character) {
				delete(c.StateMachine.AnimPlayer.Animations, "idle")
			},
			wantSub: `"idle"`,
		},
		"unknown cancel target": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["A"].FrameData[0].CancelTypes = []string{"Nope"}
			},
			wantSub: `"Nope"`,
		},
		"unknown animationSwitch target": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["A"].FrameData[0].AnimationSwitch = "Nope"
			},
			wantSub: `"Nope"`,
		},
		"zero duration": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["A"].FrameData[0].Duration = 0
			},
			wantSub: "duration",
		},
		"spriteIndex out of range": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["A"].FrameData[0].SpriteIndex = 9
			},
			wantSub: "spriteIndex",
		},
		"empty imgPath": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["A"].Sprites[0].ImagePath = ""
			},
			wantSub: "imgPath missing",
		},
		"imgPath not on disk": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["A"].Sprites[0].ImagePath = filepath.Join(t.TempDir(), "ghost.png")
			},
			wantSub: "not on disk",
		},
		"loopFrames out of range": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["6"].LoopFrames = &animation.LoopFrame{Start: 0, End: 99}
			},
			wantSub: "loopFrames",
		},
		"reaction without loopFrames": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["hurt"].LoopFrames = nil
			},
			wantSub: "must define loopFrames",
		},
		"reaction with cancelTypes": {
			mutate: func(t *testing.T, c *Character) {
				c.StateMachine.AnimPlayer.Animations["hurt"].FrameData[0].CancelTypes = []string{"A"}
			},
			wantSub: "must not carry cancelTypes",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := minCharacter(t)
			tc.mutate(t, c)
			errs := c.Validate()
			if len(errs) == 0 {
				t.Fatal("expected validation errors, got none")
			}
			joined := ""
			for _, err := range errs {
				joined += err.Error() + "\n"
			}
			if !strings.Contains(joined, tc.wantSub) {
				t.Fatalf("no error mentions %q, got:\n%s", tc.wantSub, joined)
			}
		})
	}
}

// Special motions are exempt: absent is fine, present-and-valid is fine.
func TestValidateSpecialsExempt(t *testing.T) {
	c := minCharacter(t)
	img := c.StateMachine.AnimPlayer.Animations["A"].Sprites[0].ImagePath
	c.StateMachine.AnimPlayer.Animations["236A"] = minAnim(img, false)
	if errs := c.Validate(); len(errs) > 0 {
		t.Fatalf("character with one valid special rejected: %v", errs)
	}
}

func TestIsSpecialMotion(t *testing.T) {
	for _, name := range []string{"236A", "214D", "22B", "642C", "623A", "423D", "246A"} {
		if !IsSpecialMotion(name) {
			t.Errorf("expected %s to classify as special motion", name)
		}
	}
	for _, name := range []string{"66", "44", "A", "idle", "hurt", "236", "236X", "", "6", "airBlock"} {
		if IsSpecialMotion(name) {
			t.Errorf("expected %s to NOT classify as special motion", name)
		}
	}
}

// Strict decoding must reject unknown fields at load.
func TestLoadRejectsUnknownFields(t *testing.T) {
	dir := t.TempDir()
	img := writeTempImg(t)
	content := "name: Strict\nproperties:\n  maxHP: 100\nstateMachine:\n  activeAnim:\n    animations:\n      idle:\n        sprites:\n        - imgPath: " + img + "\n        framedata:\n        - duration: 2\ntypo_field: 1\n"
	charDir := filepath.Join(dir, "assets", "characters")
	if err := os.MkdirAll(charDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(charDir, "Strict.yaml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	if _, err := loadCharacterByName("Strict"); err == nil {
		t.Fatal("expected strict-decode error for unknown field, got nil")
	} else if !strings.Contains(err.Error(), "typo_field") {
		t.Fatalf("error should name the unknown field, got: %v", err)
	}
}
