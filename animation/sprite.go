package animation

import (
	"fgengine/collision"
	"fgengine/types"
)

// TODO, check if sprite can be simplified and other info be sent do frame properties
type Sprite struct {
	ImagePath      string       `yaml:"imgPath"`
	Duration       uint         `yaml:"duration"`
	SourceSize     types.Rect   `yaml:"sourceSize"`
	Anchor         types.Rect   `yaml:"anchorPoint"`
	CollisionBoxes []types.Rect `yaml:"collisionBoxes"`
	HurtBoxes      []types.Rect `yaml:"hurtBoxes"`
	HitBoxes       []types.Rect `yaml:"hitBoxes"`

	// Cached combined boxes to avoid allocations on every GetAllBoxes call
	allBoxesCache []collision.Box
	cacheValid    bool
}

// GetAllBoxes returns a slice containing all boxes (collision, hurt, hit) for this sprite
// The returned slice is cached to avoid allocations on frequent calls
func (s *Sprite) GetAllBoxes() []collision.Box {
	if !s.cacheValid {
		s.rebuildBoxCache()
	}
	return s.allBoxesCache
}

// InvalidateBoxCache marks the box cache as invalid, forcing a rebuild on next GetAllBoxes call
func (s *Sprite) InvalidateBoxCache() {
	s.cacheValid = false
}

// rebuildBoxCache rebuilds the combined box cache from the individual box slices
func (s *Sprite) rebuildBoxCache() {
	totalBoxes := len(s.CollisionBoxes) + len(s.HurtBoxes) + len(s.HitBoxes)

	// Reuse existing slice if it has sufficient capacity, otherwise allocate new
	if cap(s.allBoxesCache) < totalBoxes {
		s.allBoxesCache = make([]collision.Box, 0, totalBoxes)
	} else {
		s.allBoxesCache = s.allBoxesCache[:0] // Reset length but keep capacity
	}
	for _, rect := range s.CollisionBoxes {
		s.allBoxesCache = append(s.allBoxesCache, collision.Box{
			Rect:    rect,
			BoxType: collision.Collision,
		})
	}
	for _, rect := range s.HurtBoxes {
		s.allBoxesCache = append(s.allBoxesCache, collision.Box{
			Rect:    rect,
			BoxType: collision.Hurt,
		})
	}
	for _, rect := range s.HitBoxes {
		s.allBoxesCache = append(s.allBoxesCache, collision.Box{
			Rect:    rect,
			BoxType: collision.Hit,
		})
	}

	s.cacheValid = true
}
