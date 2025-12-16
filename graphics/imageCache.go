package graphics

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	imageCache = make(map[string]*ebiten.Image)
	cacheMutex sync.RWMutex
)

func LoadImage(path string) *ebiten.Image {
	cacheMutex.RLock()
	if img, exists := imageCache[path]; exists {
		cacheMutex.RUnlock()
		return img
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if img, exists := imageCache[path]; exists {
		return img
	}

	image, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil || image == nil {
		// Try to load the default image
		defaultPath := "assets/common/notFound.png"
		if defaultImg, exists := imageCache[defaultPath]; exists {
			return defaultImg
		}

		// Load default image if not in cache
		defaultImage, _, defaultErr := ebitenutil.NewImageFromFile(defaultPath)
		if defaultErr != nil {
			log.Panicf("Failed to load default image %s: %v", defaultPath, defaultErr)
		}

		imageCache[defaultPath] = defaultImage
		return defaultImage
	}

	imageCache[path] = image
	return image
}

func ClearImageCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	for _, img := range imageCache {
		if img != nil {
			img.Deallocate()
		}
	}
	imageCache = make(map[string]*ebiten.Image)
}
