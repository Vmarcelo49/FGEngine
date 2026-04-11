package character

import (
	"fgengine/animation"
	"fgengine/types"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

const jumpHeight = 4
const jumpDistance = 5

func TestMakeCharacter(t *testing.T) {
	char := &Character{}

	char.Name = "PlaceHolder"
	char.StateMachine = &animation.StateMachine{
		HP:       100,
		Position: types.Vector2{X: 0, Y: 0},
		Velocity: types.Vector2{X: 0, Y: 0},
		Facing:   animation.Right,
	}

	idleSprite := animation.Sprite{
		ImagePath: "../assets/common/idle.png",
		Rect:      types.Rect{W: 256, H: 256},
		Anchor:    types.Vector2{X: 127, Y: 177},
	}
	walkSprite := animation.Sprite{
		ImagePath: "../assets/common/walk.png",
		Rect:      types.Rect{W: 256, H: 256},
		Anchor:    types.Vector2{X: 127, Y: 177},
	}
	fdIdle := animation.FrameData{
		Duration:    6,
		SpriteIndex: 0,
		CancelTypes: []string{"any"},
	}
	fdWalk := animation.FrameData{
		Duration:     6,
		SpriteIndex:  0,
		IncVelocityX: 2,
		CancelTypes:  []string{"any"},
	}
	fdWalkBack := fdWalk
	fdWalkBack.IncVelocityX = -2
	idleAnim := &animation.Animation{
		Name:      "idle",
		Sprites:   []*animation.Sprite{&idleSprite},
		FrameData: []animation.FrameData{fdIdle},
	}
	walkAnim := &animation.Animation{
		Name:      "6",
		Sprites:   []*animation.Sprite{&walkSprite},
		FrameData: []animation.FrameData{fdWalk},
	}
	walkBackAnim := &animation.Animation{
		Name:      "4",
		Sprites:   []*animation.Sprite{&walkSprite},
		FrameData: []animation.FrameData{fdWalkBack},
	}

	// Jump animations using idle sprite
	fdJumpBack := animation.FrameData{
		Duration:     6,
		SpriteIndex:  0,
		IncVelocityX: -jumpDistance,
		IncVelocityY: -jumpHeight,
		CancelTypes:  []string{},
	}
	fdJumpNeutral := animation.FrameData{
		Duration:     6,
		SpriteIndex:  0,
		IncVelocityY: -jumpHeight,
		CancelTypes:  []string{},
	}
	fdJumpForward := animation.FrameData{
		Duration:     6,
		SpriteIndex:  0,
		IncVelocityX: jumpDistance,
		IncVelocityY: -jumpHeight,
		CancelTypes:  []string{},
	}
	fdFall := animation.FrameData{
		Duration:    6,
		SpriteIndex: 0,
		CancelTypes: []string{},
	}
	fdLanding := animation.FrameData{
		Duration:    4,
		SpriteIndex: 0,
		CancelTypes: []string{},
	}

	jumpBackAnim := &animation.Animation{
		Name:      "7",
		Sprites:   []*animation.Sprite{&idleSprite},
		FrameData: []animation.FrameData{fdJumpBack},
	}
	jumpNeutralAnim := &animation.Animation{
		Name:      "8",
		Sprites:   []*animation.Sprite{&idleSprite},
		FrameData: []animation.FrameData{fdJumpNeutral},
	}
	jumpForwardAnim := &animation.Animation{
		Name:      "9",
		Sprites:   []*animation.Sprite{&idleSprite},
		FrameData: []animation.FrameData{fdJumpForward},
	}
	fallAnim := &animation.Animation{
		Name:      "fall",
		Sprites:   []*animation.Sprite{&idleSprite},
		FrameData: []animation.FrameData{fdFall},
	}
	landingAnim := &animation.Animation{
		Name:      "landing",
		Sprites:   []*animation.Sprite{&idleSprite},
		FrameData: []animation.FrameData{fdLanding},
	}

	char.StateMachine.ActiveAnim = &animation.AnimationPlayer{
		Animations: map[string]*animation.Animation{
			"idle":    idleAnim,
			"6":       walkAnim,
			"4":       walkBackAnim,
			"7":       jumpBackAnim,
			"8":       jumpNeutralAnim,
			"9":       jumpForwardAnim,
			"fall":    fallAnim,
			"landing": landingAnim,
		},
	}
	char.StateMachine.ActiveAnim.SetAnimation("idle")

	err := exportCharacterToYAML(char)
	if err != nil {
		t.Fatalf("Failed to export character to YAML: %v", err)
	}
}

func exportCharacterToYAML(c *Character) error {
	assetsDir := "../assets/characters"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets/characters directory: %w", err)
	}

	filename := fmt.Sprintf("%s.yaml", c.Name)
	path := filepath.Join(assetsDir, filename)

	originalPaths := make(map[*animation.Sprite]string)
	defer func() {
		for sprite, originalPath := range originalPaths {
			sprite.ImagePath = originalPath
		}
	}()

	animations := map[string]*animation.Animation{}
	if c.StateMachine != nil && c.StateMachine.ActiveAnim != nil && c.StateMachine.ActiveAnim.Animations != nil {
		animations = c.StateMachine.ActiveAnim.Animations
	}

	for _, anim := range animations {
		for _, sprite := range anim.Sprites {
			if sprite.ImagePath != "" {
				originalPaths[sprite] = sprite.ImagePath
				sprite.ImagePath = makeRelativePath(sprite.ImagePath, path)
			}
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create YAML file: %w", err)
	}
	defer file.Close()

	yamlInfo, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal character to YAML: %w", err)
	}

	if _, err = file.Write(yamlInfo); err != nil {
		return fmt.Errorf("failed to write YAML to file: %w", err)
	}

	return nil
}

func makeRelativePath(absolutePath, referencePath string) string {
	referenceDir := filepath.Dir(referencePath)

	absPath, err := filepath.Abs(absolutePath)
	if err != nil {
		return absolutePath
	}

	absReferenceDir, err := filepath.Abs(referenceDir)
	if err != nil {
		return absolutePath
	}

	relativePath, err := filepath.Rel(absReferenceDir, absPath)
	if err != nil {
		return absolutePath
	}

	return relativePath
}
