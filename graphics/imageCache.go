package graphics

import (
	"FGEngine/character"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
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
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	image, _, err := ebitenutil.NewImageFromFile(character.CurrentSprite.ImagePath)
	if err != nil {
		log.Panic(err)
		return nil
	}
	imageCache[character.CurrentSprite.ImagePath] = image
	return image
// TODO, Revise mutex usage
// TODO, implement clearing of the image cache
