package editor

import (
	"errors"
	"fgengine/animation"
	"fgengine/types"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sqweek/dialog"
	"gopkg.in/yaml.v2"
)

func ExportAnimationToYaml(source *animation.Animation, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create YAML file: %w", err)
	}
	defer file.Close()

	yamlInfo, err := yaml.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal animation to YAML: %w", err)
	}

	if _, err = file.Write(yamlInfo); err != nil {
		return fmt.Errorf("failed to write YAML to file: %w", err)
	}

	return nil
}

func LoadAnimationFromYAML() (animation.Animation, error) {
	path, err := dialog.File().Filter(".yaml", "yaml").Load()
	if err != nil {
		return animation.Animation{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return animation.Animation{}, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	anim := animation.Animation{}

	if err = decoder.Decode(&anim); err != nil {
		return animation.Animation{}, err
	}
	return anim, nil
}

func ExportCharacterToYaml(c *animation.Character, path string) error {
	tempCharacter := *c
	tempCharacter.Animations = make(map[string]*animation.Animation)

	for name, anim := range c.Animations {
		tempAnim := deepCopyAnimation(anim)
		for _, sprite := range tempAnim.Sprites {
			if sprite.ImagePath != "" {
				sprite.ImagePath = makeRelativePath(sprite.ImagePath, path)
			}
		}

		tempCharacter.Animations[name] = tempAnim
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create YAML file: %w", err)
	}
	defer file.Close()

	yamlInfo, err := yaml.Marshal(&tempCharacter)
	if err != nil {
		return fmt.Errorf("failed to marshal character to YAML: %w", err)
	}

	if _, err = file.Write(yamlInfo); err != nil {
		return fmt.Errorf("failed to write YAML to file: %w", err)
	}

	return nil
}

// Path utility functions for YAML serialization

// makeRelativePath converts an absolute path to a relative path based on a reference path
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

// resolveRelativePath converts a relative path to an absolute path based on a reference path
func resolveRelativePath(relativePath, referencePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	referenceDir := filepath.Dir(referencePath)
	return filepath.Clean(filepath.Join(referenceDir, relativePath))
}

func LoadCharacterFromYAML() (*animation.Character, error) {
	path, err := dialog.File().Filter(".yaml", "yaml").Load()
	if err != nil {
		return nil, errors.New("failed to load character: user cancelled")
	}
	file, err := os.Open(path)
	if err != nil {
		dialog.Message("Failed to open file: %s", err.Error()).Error()
		return nil, errors.New("failed to open character file")
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	character := &animation.Character{}
	if err := decoder.Decode(character); err != nil {
		dialog.Message("Failed to decode character: %s", err.Error()).Error()
		return nil, errors.New("failed to decode character")
	}

	// Convert relative paths to absolute paths based on YAML file location
	for _, anim := range character.Animations {
		for _, sprite := range anim.Sprites {
			if sprite.ImagePath != "" {
				sprite.ImagePath = resolveRelativePath(sprite.ImagePath, path)
			}
		}
	}

	character.FilePath = path
	return character, nil
}

func deepCopyAnimation(a *animation.Animation) *animation.Animation {
	animCopy := &animation.Animation{
		Name: a.Name,
		Prop: make([]animation.FrameProperties, len(a.Prop)),
	}

	copy(animCopy.Prop, a.Prop)

	animCopy.Sprites = make([]*animation.Sprite, len(a.Sprites))
	for i, sprite := range a.Sprites {
		animCopy.Sprites[i] = deepCopySprite(sprite)
	}

	return animCopy
}

func deepCopySprite(source *animation.Sprite) *animation.Sprite {
	other := &animation.Sprite{}
	other.ImagePath = source.ImagePath
	other.Duration = source.Duration
	other.SourceSize = source.SourceSize
	other.Anchor = source.Anchor

	other.CollisionBoxes = make([]types.Rect, len(source.CollisionBoxes))
	copy(other.CollisionBoxes, source.CollisionBoxes)

	other.HurtBoxes = make([]types.Rect, len(source.HurtBoxes))
	copy(other.HurtBoxes, source.HurtBoxes)

	other.HitBoxes = make([]types.Rect, len(source.HitBoxes))
	copy(other.HitBoxes, source.HitBoxes)

	return other
}
