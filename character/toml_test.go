package character

import (
	"bytes"
	"fgengine/constants"
	"fgengine/types"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "..")
}

func placeholderPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "assets", "characters", "PlaceHolder.toml")
}

func TestPlaceholderLoadsAndValidates(t *testing.T) {
	chara, err := DecodeCharacterFile(placeholderPath(t))
	if err != nil {
		t.Fatalf("decode placeholder: %v", err)
	}
	if chara.Name != "PlaceHolder" {
		t.Fatalf("name = %q, want PlaceHolder", chara.Name)
	}
	if chara.Properties.MaxHP != 10000 {
		t.Fatalf("maxHP = %d, want 10000", chara.Properties.MaxHP)
	}
	if errs := chara.Validate(); len(errs) > 0 {
		t.Fatalf("placeholder invalid: %v", errs)
	}

	anims := chara.StateMachine.AnimPlayer.Animations
	for _, name := range append(append([]string{}, RequiredStates...), "236A") {
		if anims[name] == nil {
			t.Errorf("placeholder missing %q", name)
		}
	}

	// Named box keys must decode (collision/hit/hurt, SPEC §6.5).
	active := anims["A"].FrameData[1]
	if len(active.Boxes[types.Hit]) != 1 {
		t.Fatalf("A active frame should carry 1 hit box, got %d", len(active.Boxes[types.Hit]))
	}
	if len(active.Boxes[types.Hurt]) != 1 {
		t.Fatalf("A active frame should carry 1 hurt box, got %d", len(active.Boxes[types.Hurt]))
	}
	if active.Damage != 100 || active.Hitstun != 12 || active.Knockback != 2 {
		t.Fatalf("A active frame data wrong: %+v", active)
	}
}

// Save → load → save must be byte-stable, and saved files stay valid.
func TestRoundTripStable(t *testing.T) {
	chara, err := DecodeCharacterFile(placeholderPath(t))
	if err != nil {
		t.Fatalf("decode placeholder: %v", err)
	}

	enc1, err := EncodeTOML(chara)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	tmp := filepath.Join(t.TempDir(), "rt.toml")
	if err := os.WriteFile(tmp, enc1, 0644); err != nil {
		t.Fatal(err)
	}
	chara2, err := DecodeCharacterFile(tmp)
	if err != nil {
		t.Fatalf("re-decode: %v", err)
	}
	if errs := chara2.Validate(); len(errs) > 0 {
		t.Fatalf("re-decoded character invalid: %v", errs)
	}
	enc2, err := EncodeTOML(chara2)
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if !bytes.Equal(enc1, enc2) {
		t.Fatalf("round-trip unstable: %d vs %d bytes", len(enc1), len(enc2))
	}
}

// Full runtime load path: strict, validated, initialized from file data.
func TestLoadCharacterInitializes(t *testing.T) {
	t.Chdir(repoRoot(t))
	chara, err := LoadCharacter("PlaceHolder", 1)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if chara.StateMachine.HP != 10000 {
		t.Fatalf("HP = %d, want maxHP 10000 from file", chara.StateMachine.HP)
	}
	if chara.StateMachine.Position.Y != constants.GroundLevelY {
		t.Fatalf("spawn Y = %v, want grounded %v", chara.StateMachine.Position.Y, constants.GroundLevelY)
	}
}
