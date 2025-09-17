package collision

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
