//go:build !js && !wasm
// +build !js,!wasm

package editorimgui

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/types"
	"fmt"
	"image"
	"os"

	"github.com/sqweek/dialog"
)

func (g *Game) addSpritesFromFiles(paths []string) error {
	if g.ActiveAnimation() == nil {
		return fmt.Errorf("no active animation available")
	}
	if len(paths) == 0 {
		return fmt.Errorf("no files selected")
	}

	for _, path := range paths {
		sprite, err := loadSpriteFromImagePath(path)
		if err != nil {
			return fmt.Errorf("error creating sprite from file %s: %w", path, err)
		}
		g.ActiveAnimation().Sprites = append(g.ActiveAnimation().Sprites, sprite)
		newFrameData := animation.FrameData{
			Duration:    1,
			SpriteIndex: len(g.ActiveAnimation().Sprites) - 1,
		}
		g.ActiveAnimation().FrameData = append(g.ActiveAnimation().FrameData, newFrameData)
	}
	return nil
}

func (g *Game) newAnimationFileDialog() (*animation.Animation, error) {
	path, err := dialog.File().Filter("PNG Image", "png").Load()
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
