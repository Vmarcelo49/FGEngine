package editor

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/filepicker"
	"fgengine/state"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func exportCharacterToYAML(c *character.Character) error {
	assetsDir := "assets/characters"
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

	for _, anim := range c.Animations {
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

func resolveRelativePath(relativePath, referencePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	referenceDir := filepath.Dir(referencePath)
	return filepath.Clean(filepath.Join(referenceDir, relativePath))
}

func loadCharacterFromPath(path string) (*character.Character, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open character file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	char := &character.Character{
		Animations:      make(map[string]*animation.Animation),
		StateMachine:    &state.StateMachine{},
		AnimationPlayer: &animation.AnimationPlayer{},
	}

	if err := decoder.Decode(char); err != nil {
		return nil, fmt.Errorf("failed to decode character: %w", err)
	}

	for _, anim := range char.Animations {
		for _, sprite := range anim.Sprites {
			if sprite.ImagePath != "" {
				sprite.ImagePath = resolveRelativePath(sprite.ImagePath, path)
			}
		}
	}

	char.FilePath = path
	return char, nil
}

func loadCharacterFromYAML(characterName string) (*character.Character, error) {
	path := filepath.Join("assets/characters", fmt.Sprintf("%s.yaml", characterName))
	return loadCharacterFromPath(path)
}

func loadCharacterFromYAMLDialog() (*character.Character, error) {
	picker := filepicker.GetFilePicker()
	filter := filepicker.FileFilter{
		Description: ".yaml",
		Extensions:  []string{"yaml"},
	}

	path, err := picker.LoadFile(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to load character: user cancelled")
	}

	return loadCharacterFromPath(path)
}

func exportAnimationToYaml(source *animation.Animation, path string) error {
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

func loadAnimationFromYAML() (animation.Animation, error) {
	picker := filepicker.GetFilePicker()
	filter := filepicker.FileFilter{
		Description: ".yaml",
		Extensions:  []string{"yaml"},
	}

	path, err := picker.LoadFile(filter)
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
