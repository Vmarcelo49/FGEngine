package types

type Rect struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
	W float64 `yaml:"w"`
	H float64 `yaml:"h"`
}

type Vector2 struct {
	X float64
	Y float64
}

func (r Rect) Right() float64 {
	return r.X + r.W
}

func (r Rect) Bottom() float64 {
	return r.Y + r.H
}

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
