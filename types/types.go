package types

import "math"

type Rect struct {
	X float64 `yaml:"x,omitempty"`
	Y float64 `yaml:"y,omitempty"`
	W float64 `yaml:"w"`
	H float64 `yaml:"h"`
}

type Vector2 struct {
	X float64 `yaml:"x,omitempty"`
	Y float64 `yaml:"y,omitempty"`
}

type BoxType uint8

const (
	Collision BoxType = iota
	Hit
	Hurt
)

func (b BoxType) String() string {
	switch b {
	case Collision:
		return "Collision"
	case Hit:
		return "Hit"
	case Hurt:
		return "Hurt"
	default:
		return "Unknown"
	}
}

func Normalize(v Vector2) Vector2 {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y)
	if length == 0 {
		return Vector2{X: 0, Y: 0}
	}
	return Vector2{X: v.X / length, Y: v.Y / length}
}

func (r Rect) Right() float64 {
	return r.X + r.W
}

func (r Rect) Bottom() float64 {
	return r.Y + r.H
}

// Center returns the center coordinates of the rectangle
func (r Rect) Center() (float64, float64) {
	return r.X + r.W/2, r.Y + r.H/2
}

func (r Rect) Contains(x, y float64) bool {
	return x >= r.X && x < r.Right() && y >= r.Y && y < r.Bottom()
}

// Regular AABB collision check, does not consider edges as overlapping (exclusive)
func (r Rect) IsOverlapping(other Rect) bool {
	return r.X < other.Right() &&
		r.Right() > other.X &&
		r.Y < other.Bottom() &&
		r.Bottom() > other.Y
}

// Checks if two rectangles are overlapping, including edges (inclusive)
func (r Rect) IsOverlapInclusive(other Rect) bool {
	return r.X <= other.Right() &&
		r.Right() >= other.X &&
		r.Y <= other.Bottom() &&
		r.Bottom() >= other.Y
}

// CenterWithin centers the rect within the parent Rect
func (r *Rect) CenterWithin(parent Rect) {
	r.X = parent.X + (parent.W-r.W)/2
	r.Y = parent.Y + (parent.H-r.H)/2
}

func (r *Rect) Pos() Vector2 {
	return Vector2{X: r.X, Y: r.Y}
}

func (r *Rect) Size() Vector2 {
	return Vector2{X: r.W, Y: r.H}
}

func (v Vector2) Add(o Vector2) Vector2 {
	return Vector2{v.X + o.X, v.Y + o.Y}
}

func (v Vector2) Sub(o Vector2) Vector2 {
	return Vector2{v.X - o.X, v.Y - o.Y}
}

func (v Vector2) Mul(s float64) Vector2 {
	return Vector2{v.X * s, v.Y * s}
}
