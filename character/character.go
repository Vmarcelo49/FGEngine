package character

import (
	"fgengine/animation"
	"fgengine/constants"
	"fgengine/types"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

type Character struct {
	Name         string                  `yaml:"name"`
	StateMachine *animation.StateMachine `yaml:"stateMachine"`
}

func LoadCharacter(name string, playerSide int) (*Character, error) {
	chara, err := loadCharacterByName(name)
	if err != nil {
		return nil, err
	}
	chara.initialize(playerSide)
	return chara, nil
}

func loadCharacterByName(name string) (*Character, error) {
	filePath := "./assets/characters/" + name + ".yaml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read character file: %w", err)
	}

	character := &Character{
		Name: name,
	}
	if err := yaml.Unmarshal(data, character); err != nil {
		return nil, fmt.Errorf("failed to unmarshal character data: %w", err)
	}

	if character.StateMachine == nil || character.StateMachine.AnimPlayer == nil {
		return nil, fmt.Errorf("character file is missing stateMachine.activeAnim")
	}
	if character.StateMachine.AnimPlayer.Animations == nil {
		return nil, fmt.Errorf("character file is missing stateMachine.activeAnim.animations")
	}

	// Keep runtime animation names in sync with the map keys.
	for animName, anim := range character.StateMachine.AnimPlayer.Animations {
		if anim == nil {
			continue
		}
		anim.Name = animName

		for _, sprite := range anim.Sprites {
			if sprite.ImagePath != "" {
				sprite.ImagePath = resolveRelativePath(sprite.ImagePath, filePath)
			}
		}
	}
	return character, nil
}

func (c *Character) initialize(playerSide int) {
	if c.StateMachine == nil {
		c.StateMachine = new(animation.StateMachine{})
	}
	if c.StateMachine.AnimPlayer == nil {
		c.StateMachine.AnimPlayer = new(animation.AnimationPlayer{})
	}

	var initialX float64
	var facing animation.Orientation
	switch playerSide {
	case 1:
		initialX = constants.WorldWidth / 4
		facing = animation.Right
	case 2:
		initialX = 3 * constants.WorldWidth / 4
		facing = animation.Left
	}

	c.StateMachine.HP = 10000
	c.StateMachine.Position = types.Vector2{X: initialX, Y: constants.WorldHeight / 2}
	c.StateMachine.IsFacingLeft = facing
	c.StateMachine.Velocity = types.Vector2{}
	c.StateMachine.IgnoreGravityFrames = 0

	setInitialAnimation(c.StateMachine.AnimPlayer)

}

func setInitialAnimation(player *animation.AnimationPlayer) {
	if player == nil || len(player.Animations) == 0 || player.ActiveAnimation != nil {
		return
	}

	if _, ok := player.Animations["idle"]; ok {
		player.SetAnimation("idle")
		return
	}

	animNames := make([]string, 0, len(player.Animations))
	for name := range player.Animations {
		animNames = append(animNames, name)
	}
	slices.Sort(animNames)
	player.SetAnimation(animNames[0])
}

// resolveRelativePath converts a relative path to an absolute path based on a reference path
func resolveRelativePath(relativePath, referencePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	referenceDir := filepath.Dir(referencePath)
	return filepath.Clean(filepath.Join(referenceDir, relativePath))
}

func (c *Character) Position() types.Vector2 {
	return c.StateMachine.Position
}

func (c *Character) Sprite() *animation.Sprite {
	return c.StateMachine.AnimPlayer.ActiveSprite()
}
