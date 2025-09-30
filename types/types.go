package types

import "math"

type Rect struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
	W float64 `yaml:"w"`
	H float64 `yaml:"h"`
}

type Vector2 struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
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

func (r Rect) IsOverlapping(other Rect) bool {
	return r.X < other.Right() && r.Right() > other.X &&
		r.Y < other.Bottom() && r.Bottom() > other.Y
}

// AlignCenter centers the rect within the parentRect
func (r *Rect) AlignCenter(parentRect Rect) {
	r.X = parentRect.X + (parentRect.W-r.W)/2
	r.Y = parentRect.Y + (parentRect.H-r.H)/2
}
