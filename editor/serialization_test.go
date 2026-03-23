package editor

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCharacterFromPathUsesCurrentStateMachineSchema(t *testing.T) {
	tmp := t.TempDir()
	charDir := filepath.Join(tmp, "assets", "characters")
	if err := os.MkdirAll(charDir, 0755); err != nil {
		t.Fatalf("failed to create character dir: %v", err)
	}

	yamlPath := filepath.Join(charDir, "SchemaTest.yaml")
	yamlContent := `name: SchemaTest
stateMachine:
  activeAnim:
    animations:
      idle:
        name: idle
        sprites:
          - imgPath: ../common/idle.png
            rect:
              w: 100
              h: 200
        framedata:
          - duration: 1
            spriteIndex: 0
`
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml: %v", err)
	}

	chara, err := loadCharacterFromPath(yamlPath)
	if err != nil {
		t.Fatalf("loadCharacterFromPath returned error: %v", err)
	}

	if chara.StateMachine == nil || chara.StateMachine.ActiveAnim == nil {
		t.Fatalf("state machine was not initialized")
	}

	idle := chara.StateMachine.ActiveAnim.Animations["idle"]
	if idle == nil {
		t.Fatalf("idle animation not found")
	}

	if len(idle.Sprites) != 1 {
		t.Fatalf("expected one sprite, got %d", len(idle.Sprites))
	}

	expectedPath := filepath.Clean(filepath.Join(tmp, "assets", "common", "idle.png"))
	if idle.Sprites[0].ImagePath != expectedPath {
		t.Fatalf("unexpected sprite path: got %q want %q", idle.Sprites[0].ImagePath, expectedPath)
	}
}

func TestExportCharacterToYAMLWritesRelativeSpritePathsAndRestoresMemory(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to change cwd: %v", err)
	}

	absSpritePath := filepath.Join(tmp, "assets", "common", "idle.png")
	if err := os.MkdirAll(filepath.Dir(absSpritePath), 0755); err != nil {
		t.Fatalf("failed to create sprite dir: %v", err)
	}
	if err := os.WriteFile(absSpritePath, []byte("fake"), 0644); err != nil {
		t.Fatalf("failed to write sprite file: %v", err)
	}

	sprite := &animation.Sprite{ImagePath: absSpritePath, Rect: types.Rect{W: 64, H: 64}}
	chara := &character.Character{
		Name: "ExportSchemaTest",
		StateMachine: &animation.StateMachine{
			ActiveAnim: &animation.AnimationPlayer{
				Animations: map[string]*animation.Animation{
					"idle": {
						Name:      "idle",
						Sprites:   []*animation.Sprite{sprite},
						FrameData: []animation.FrameData{{Duration: 1, SpriteIndex: 0}},
					},
				},
			},
		},
	}

	if err := exportCharacterToYAML(chara); err != nil {
		t.Fatalf("exportCharacterToYAML returned error: %v", err)
	}

	if sprite.ImagePath != absSpritePath {
		t.Fatalf("sprite path should be restored after export: got %q want %q", sprite.ImagePath, absSpritePath)
	}

	savedPath := filepath.Join(tmp, "assets", "characters", "ExportSchemaTest.yaml")
	content, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("failed to read exported yaml: %v", err)
	}

	if !strings.Contains(string(content), "imgPath: ../common/idle.png") {
		t.Fatalf("exported yaml did not contain relative sprite path; yaml:\n%s", string(content))
	}
}
