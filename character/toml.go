package character

import (
	"bytes"
	"errors"
	"fgengine/animation"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// tomlCharacterFile mirrors the flat SPEC §6.2 file layout. The runtime
// structs stay nested (Character → StateMachine → AnimPlayer), so this aux
// type bridges file layout and memory layout in both directions.
type tomlCharacterFile struct {
	Name       string                          `toml:"name"`
	Properties CharacterProperties             `toml:"properties"`
	Animations map[string]*animation.Animation `toml:"animations"`
}

// DecodeCharacterFile reads path, strict-decodes the flat TOML layout,
// resolves sprite paths against the file location, and syncs animation
// names. No validation, no init — the editor opens work-in-progress files
// through this path. Runtime loading must follow with Validate.
func DecodeCharacterFile(path string) (*Character, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read character file: %w", err)
	}

	var tc tomlCharacterFile
	decoder := toml.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&tc); err != nil {
		var strict *toml.StrictMissingError
		if errors.As(err, &strict) {
			return nil, fmt.Errorf("failed to decode character data (unknown fields):\n%s", strict.String())
		}
		return nil, fmt.Errorf("failed to decode character data: %w", err)
	}

	chara := &Character{
		Name:       tc.Name,
		Properties: tc.Properties,
	}
	if tc.Animations != nil {
		chara.StateMachine = &animation.StateMachine{
			AnimPlayer: &animation.AnimationPlayer{Animations: tc.Animations},
		}
		for animName, anim := range tc.Animations {
			if anim == nil {
				continue
			}
			anim.Name = animName
			for _, sprite := range anim.Sprites {
				if sprite != nil && sprite.ImagePath != "" {
					sprite.ImagePath = resolveRelativePath(sprite.ImagePath, path)
				}
			}
		}
	}
	return chara, nil
}

// EncodeTOML serializes the character to the flat SPEC §6.2 layout.
func EncodeTOML(c *Character) ([]byte, error) {
	tc := tomlCharacterFile{
		Name:       c.Name,
		Properties: c.Properties,
	}
	if c.StateMachine != nil && c.StateMachine.AnimPlayer != nil {
		tc.Animations = c.StateMachine.AnimPlayer.Animations
	}
	return toml.Marshal(tc)
}
