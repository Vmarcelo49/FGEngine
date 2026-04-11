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

func checkSOCD(input *GameInput) {
	if input.IsPressed(Left) && input.IsPressed(Right) {
		*input &^= (Left | Right)
	}

	if input.IsPressed(Up) && input.IsPressed(Down) {
		*input &^= (Up | Down)
	}
}
