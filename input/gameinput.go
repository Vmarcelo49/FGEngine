package input

type GameInput byte

const (
	NoInput GameInput = 0
)

const (
	Up GameInput = 1 << iota
	Down
	Left
	Right
	A
	B
	C
	D
)

func (gi GameInput) String() string {
	if gi == NoInput {
		return "NoInput"
	}
	str := ""
	if gi&Up != 0 {
		str += "8 "
	}
	if gi&Down != 0 {
		str += "2 "
	}
	if gi&Left != 0 {
		str += "4 "
	}
	if gi&Right != 0 {
		str += "6 "
	}
	if gi&A != 0 {
		str += "A "
	}
	if gi&B != 0 {
		str += "B "
	}
	if gi&C != 0 {
		str += "C "
	}
	if gi&D != 0 {
		str += "D "
	}
	return str
}

func (gi GameInput) IsPressed(input GameInput) bool {
	return gi&input != 0
}

func JustPressed(current, previous, button GameInput) bool {
	return current.IsPressed(button) && !previous.IsPressed(button)
}

func JustReleased(current, previous, button GameInput) bool {
	return !current.IsPressed(button) && previous.IsPressed(button)
}

// ApplySOCD clears simultaneous opposing directions to neutral (SPEC §5.2).
func ApplySOCD(g GameInput) GameInput {
	if g.IsPressed(Left) && g.IsPressed(Right) {
		g &^= (Left | Right)
	}

	if g.IsPressed(Up) && g.IsPressed(Down) {
		g &^= (Up | Down)
	}
	return g
}
