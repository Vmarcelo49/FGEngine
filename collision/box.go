package collision

import (
	"FGEngine/types"
)

type BoxType uint8

const (
	Collision BoxType = iota
	Hit
	Hurt
)

type Box struct {
	Rect    types.Rect
	BoxType BoxType
}
