package graphics

import (
	"fgengine/constants"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Drawable interface {
	Draw(screen *ebiten.Image, camera *Camera)
}

type RenderQueue struct {
	layers [constants.LayerCount][]Drawable
}

func (rq *RenderQueue) Draw(screen *ebiten.Image, camera *Camera) {
	for i := range rq.layers {
		for _, renderable := range rq.layers[i] {
			renderable.Draw(screen, camera)
		}
	}
}

func (rq *RenderQueue) Add(item Drawable, layer int) {
	if layer < 0 || layer >= constants.LayerCount {
		log.Panicf("invalid layer %d (valid: 0-%d)", layer, constants.LayerCount-1)
	}
	rq.layers[layer] = append(rq.layers[layer], item)
}

func (rq *RenderQueue) Remove(item Drawable) {
	for i := range rq.layers {
		for j, renderable := range rq.layers[i] {
			if item == renderable {
				lastIdx := len(rq.layers[i]) - 1
				rq.layers[i][j] = rq.layers[i][lastIdx]
				rq.layers[i] = rq.layers[i][:lastIdx]
				return
			}
		}
	}
}

func (rq *RenderQueue) SetLast(item Drawable) bool {
	for i := range rq.layers {
		for j, renderable := range rq.layers[i] {
			if item == renderable {
				// removes from current position
				rq.layers[i] = append(rq.layers[i][:j], rq.layers[i][j+1:]...)
				// to set it to last
				rq.layers[i] = append(rq.layers[i], item)
				return true
			}
		}
	}
	return false
}

func (rq *RenderQueue) SetFirst(item Drawable) bool {
	for i := range rq.layers {
		for j, renderable := range rq.layers[i] {
			if item == renderable {
				if j != 0 {
					rq.layers[i][0], rq.layers[i][j] = rq.layers[i][j], rq.layers[i][0]

				}
				return true
			}
		}
	}
	return false
}

func (rq *RenderQueue) Clear() {
	for i := range rq.layers {
		rq.layers[i] = rq.layers[i][:0]
	}
}
