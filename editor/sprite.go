package editor

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/types"
	"fmt"
	"image"
	"os"

	"github.com/sqweek/dialog"
)

// Opens a PNG file, appends a new sprite to the active animation
func (e *EditorManager) addSpriteByFile(path string) error {
	if e == nil || e.activeAnimation == nil {
		return fmt.Errorf("no active animation available")
	}
	sprite, err := newSpriteFromImage(path)
	if err != nil {
		return fmt.Errorf("error creating sprite from image: %w", err)
	}
	e.activeAnimation.Sprites = append(e.activeAnimation.Sprites, sprite)
	newFrameData := animation.FrameData{
		Duration:    1, // Default duration
		SpriteIndex: len(e.activeAnimation.Sprites) - 1,
	}
	e.activeAnimation.FrameData = append(e.activeAnimation.FrameData, newFrameData)
	return nil
}

func (e *EditorManager) newAnimationFileDialog() (*animation.Animation, error) {
	path, err := dialog.File().Filter(".png Image", "png").Load()
	if err != nil {
		return nil, err
	}
	sprite, err := newSpriteFromImage(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create sprite from image: %w", err)
	}
	anim := &animation.Animation{
		Sprites: []*animation.Sprite{sprite},
		FrameData: []animation.FrameData{{
			Duration:    60, // Default duration of 60 frames (1 second at 60fps)
			SpriteIndex: 0,
		}},
	}
	return anim, nil
}

func newSpriteFromImage(path string) (*animation.Sprite, error) {
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
	err = makeDefaultBoxes(sprite)
	if err != nil {
		return nil, fmt.Errorf("failed to make default boxes: %w", err)
	}
	return sprite, nil
}

func makeDefaultBoxes(sprite *animation.Sprite) error {
	if sprite == nil {
		return fmt.Errorf("sprite is nil")
	}
	if sprite.Boxes == nil {
		sprite.Boxes = make(map[collision.BoxType][]types.Rect)
	}
	sprite.Boxes[collision.Hurt] = []types.Rect{
		{X: 0, Y: 0, W: float64(sprite.Rect.W), H: float64(sprite.Rect.H)},
	}
	sprite.Boxes[collision.Collision] = []types.Rect{
		{X: 0, Y: 0, W: float64(sprite.Rect.W) / 2, H: float64(sprite.Rect.H) / 2},
	}
	return nil
}
