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

func loadImage(renderable Renderable) *ebiten.Image {
	if imageCache == nil {
		imageCache = make(map[string]*ebiten.Image)
	}

	sprite := renderable.GetSprite()
	if sprite == nil {
		return nil
	}

	imagePath := sprite.ImagePath

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
		log.Printf("Failed to load image %s: %v, using default image", imagePath, err)

		// Try to load the default image
		defaultPath := "assets/common/notFound.png"
		if defaultImg, exists := imageCache[defaultPath]; exists {
			return defaultImg
		}

		// Load default image if not in cache
		defaultImage, _, defaultErr := ebitenutil.NewImageFromFile(defaultPath)
		if defaultErr != nil {
			log.Printf("Failed to load default image %s: %v", defaultPath, defaultErr)
			return nil
		}

		imageCache[defaultPath] = defaultImage
		return defaultImage
	}

	imageCache[imagePath] = image
	return image
}

func LoadImage(path string) *ebiten.Image {
	if imageCache == nil {
		imageCache = make(map[string]*ebiten.Image)
	}

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
	if err != nil {
		log.Printf("Failed to load image %s: %v, using default image", path, err)

		// Try to load the default image
		defaultPath := "assets/common/notFound.png"
		if defaultImg, exists := imageCache[defaultPath]; exists {
			return defaultImg
		}

		// Load default image if not in cache
		defaultImage, _, defaultErr := ebitenutil.NewImageFromFile(defaultPath)
		if defaultErr != nil {
			log.Printf("Failed to load default image %s: %v", defaultPath, defaultErr)
			return nil
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
	imageCache = make(map[string]*ebiten.Image)
}
