package editor

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/filepicker"
	"fgengine/types"
	"fmt"
	"image"
	"os"
)

// Opens a PNG file, appends a new sprite to the active animation
// TODO: refactor to accept file data instead of path for better cross-platform support
func (g *Game) addSpriteByFile(path string) error {
	if g.getActiveAnimation() == nil {
		return fmt.Errorf("no active animation available")
	}
	sprite, err := loadSpriteFromImagePath(path)
	if err != nil {
		return fmt.Errorf("error creating sprite from file: %w", err)
	}
	g.getActiveAnimation().Sprites = append(g.getActiveAnimation().Sprites, sprite)
	newFrameData := animation.FrameData{
		Duration:    1,
		SpriteIndex: len(g.getActiveAnimation().Sprites) - 1, // Index of the newly added sprite
	}
	g.getActiveAnimation().FrameData = append(g.getActiveAnimation().FrameData, newFrameData)
	return nil
}

func (g *Game) newAnimationFileDialog() (*animation.Animation, error) {
	picker := filepicker.GetFilePicker()
	filter := filepicker.FileFilter{
		Description: ".png Image",
		Extensions:  []string{"png"},
	}

	path, err := picker.LoadFile(filter)
	if err != nil {
		return nil, err
	}
	sprite, err := loadSpriteFromImagePath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create sprite from image: %w", err)
	}
	anim := &animation.Animation{
		Name:    "newAnimation",
		Sprites: []*animation.Sprite{sprite},
		FrameData: []animation.FrameData{{
			Duration: 1,
		}},
	}
	anim.FrameData[0].Boxes = make(map[collision.BoxType][]types.Rect)

	anim.FrameData[0].Boxes[collision.Collision] = []types.Rect{
		{X: 0, Y: 0, W: sprite.Rect.W, H: sprite.Rect.H},
	}
	anim.FrameData[0].Boxes[collision.Hurt] = []types.Rect{
		{X: 0, Y: 0, W: sprite.Rect.W, H: sprite.Rect.H},
	}
	return anim, nil
}

func loadSpriteFromImagePath(path string) (*animation.Sprite, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	decodedImage, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	imageWidth := decodedImage.Bounds().Dx()
	imageHeight := decodedImage.Bounds().Dy()

	sprite := &animation.Sprite{
		ImagePath: path,
		Rect:      types.Rect{X: 0, Y: 0, W: float64(imageWidth), H: float64(imageHeight)},
	}
	return sprite, nil
}
