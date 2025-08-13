package graphics

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	imageCache map[string]*ebiten.Image
	cacheMutex sync.RWMutex
)

func loadRenderableImage(renderable Renderable) *ebiten.Image {
	if imageCache == nil {
		imageCache = make(map[string]*ebiten.Image)
	}

	animComp := renderable.GetAnimationComponent()
	if animComp == nil || !animComp.IsValid() {
		return nil
	}

	imagePath := animComp.GetCurrentSprite().ImagePath

	cacheMutex.RLock()
	if img, exists := imageCache[imagePath]; exists {
		cacheMutex.RUnlock()
		return img
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if img, exists := imageCache[imagePath]; exists {
		return img
	}
	image, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Panic(err)
		return nil
	}
	imageCache[imagePath] = image
	return image
}

// ClearImageCache clears the image cache.
func ClearImageCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	imageCache = make(map[string]*ebiten.Image)
}
