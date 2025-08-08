package graphics

import (
	"FGEngine/character"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	imageCache map[string]*ebiten.Image
	cacheMutex sync.RWMutex
)

func loadCharacterImage(character *character.Character) *ebiten.Image {
	if imageCache == nil {
		imageCache = make(map[string]*ebiten.Image)
	}

	cacheMutex.RLock()
	if img, exists := imageCache[character.CurrentSprite.ImagePath]; exists {
		cacheMutex.RUnlock()
		return img
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if img, exists := imageCache[character.CurrentSprite.ImagePath]; exists {
		return img
	}
	image, _, err := ebitenutil.NewImageFromFile(character.CurrentSprite.ImagePath)
	if err != nil {
		log.Panic(err)
		return nil
	}
	imageCache[character.CurrentSprite.ImagePath] = image
	return image
}

// ClearImageCache clears the image cache.
func ClearImageCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	imageCache = make(map[string]*ebiten.Image)
}
