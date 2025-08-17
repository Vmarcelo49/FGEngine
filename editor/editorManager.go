package editor

import (
	"fgengine/animation"
	"fgengine/types"
	"fmt"
	"image"
	"os"

	_ "image/png"

	"github.com/sqweek/dialog"
)

type EditorManager struct {
	activeAnimation *animation.Animation
	frameCount      int
	frameIndex      int
	playingAnim     bool

	// UI related
	uiPrevAnimationName string
	logBuf              string
	logUpdated          bool
}

func (e *EditorManager) SetActiveAnimation(anim *animation.Animation) {
	e.activeAnimation = anim
	e.frameIndex = 0
	e.playingAnim = false
}

// Opens a PNG file, appends a new sprite to the active animation
func (e *EditorManager) addSpriteByFile(path string) error {
	sprite, err := newSpriteFromImage(path)
	if err != nil {
		return fmt.Errorf("error creating sprite from image: %w", err)
	}
	e.activeAnimation.Sprites = append(e.activeAnimation.Sprites, sprite) // current animation must be not nil
	e.activeAnimation.Prop = append(e.activeAnimation.Prop, animation.FrameProperties{})
	return nil
}

// newAnimationFileDialog opens a file dialog to select a PNG file and creates a new animation
func (e *EditorManager) newAnimationFileDialog() (*animation.Animation, error) {
	path, err := dialog.File().Filter(".png Image", "png").Load()
	if err != nil {
		return nil, err
	}
	sprite, err := newSpriteFromImage(path) // TODO, make paths be relative in here
	if err != nil {
		return nil, fmt.Errorf("failed to create sprite from image: %w", err)
	}
	sprite.HurtBoxes = []types.Rect{
		{X: 0, Y: 0, W: float64(sprite.SourceSize.W), H: float64(sprite.SourceSize.H)},
	}
	sprite.CollisionBoxes = []types.Rect{
		{X: 0, Y: 0, W: float64(sprite.SourceSize.W) / 2, H: float64(sprite.SourceSize.H) / 2},
	}
	anim := &animation.Animation{
		Sprites: []*animation.Sprite{sprite},
		Prop:    []animation.FrameProperties{{}},
	}
	e.activeAnimation = anim
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

	return &animation.Sprite{
		ImagePath: path,
		SourceSize: types.Rect{
			W: float64(imageWidth),
			H: float64(imageHeight),
		},
		Anchor: types.Rect{
			X: float64(imageWidth) / 2,
			Y: float64(imageHeight) / 2,
		},
	}, nil
}

func (em *EditorManager) getCurrentSprite() *animation.Sprite {
	if em.activeAnimation == nil {
		return nil
	}
	if em.frameIndex < 0 || em.frameIndex >= len(em.activeAnimation.Sprites) {
		return nil
	}
	return em.activeAnimation.Sprites[em.frameIndex]
}
