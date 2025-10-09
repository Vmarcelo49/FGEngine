package editor

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/collision"
	"fgengine/filepicker"
	"fgengine/types"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func exportCharacterToYAML(c *character.Character) error {
	assetsDir := "assets/characters"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets/characters directory: %w", err)
	}

	filename := fmt.Sprintf("%s.yaml", c.Name)
	path := filepath.Join(assetsDir, filename)

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

func loadCharacterFromYAML(characterName string) (*character.Character, error) {
	filename := fmt.Sprintf("%s.yaml", characterName)
	path := filepath.Join("assets/characters", filename)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open character file %s: %w", path, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	character := &character.Character{}
	if err := decoder.Decode(character); err != nil {
		return nil, fmt.Errorf("failed to decode character: %w", err)
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
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open character file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	character := &character.Character{}
	if err := decoder.Decode(character); err != nil {
		return nil, fmt.Errorf("failed to decode character")
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
		Name:      a.Name,
		FrameData: make([]animation.FrameData, len(a.FrameData)),
	}

	copy(animCopy.FrameData, a.FrameData)

	animCopy.Sprites = make([]*animation.Sprite, len(a.Sprites))
	for i, sprite := range a.Sprites {
		animCopy.Sprites[i] = deepCopySprite(sprite)
	}

	return animCopy
}

func deepCopySprite(source *animation.Sprite) *animation.Sprite {
	destination := &animation.Sprite{}
	destination.ImagePath = source.ImagePath
	destination.Rect = source.Rect

	destination.Boxes = make(map[collision.BoxType][]types.Rect)

	copyBoxes(source.Boxes, destination.Boxes)
	return destination
}

func copyBoxes(sourceBoxes map[collision.BoxType][]types.Rect, destBoxes map[collision.BoxType][]types.Rect) {
	for key, boxes := range sourceBoxes {
		boxesCopy := make([]types.Rect, len(boxes))
		copy(boxesCopy, boxes)
		destBoxes[key] = boxesCopy
	}
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

func listAvailableCharacters() ([]string, error) {
	assetsDir := "assets/characters"

	entries, err := os.ReadDir(assetsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read assets/characters directory: %w", err)
	}

	var characters []string
	for _, entry := range entries { // TODO, add proper validation of the yaml files
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
			name := strings.TrimSuffix(entry.Name(), ".yaml")
			characters = append(characters, name)
		}
	}

	return characters, nil
}
