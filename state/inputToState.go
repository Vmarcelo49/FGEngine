package state

import "fgengine/input"

// basic movement

var movementInputs = map[input.GameInput]State{
	input.Right: StateDash | StateForward,
	input.Left:  StateDash | StateBackward,
	input.Down:  StateCrouching | StateNeutral,
	input.Up:    StateJump | StateNeutral,
}

// states are becoming any game action
func inputSequenceToState(inputSeq []input.GameInput) State {
	InputSequenceDash := []input.GameInput{input.Right, input.Right} // this can accept neutral inputs in between, and also be buffered, down/up or reverse direction can cancel it
	InputSequence236 := []input.GameInput{input.Down, input.Right & input.Right, input.Right}
	InputSequence214 := []input.GameInput{input.Down, input.Left & input.Left, input.Left}
	return StateDash | StateForward
}

type InputSequence struct {
	baseInput    []input.GameInput
	tolerance    int // maybe the buffer?
	cancelInputs []input.GameInput
}

var Dash = InputSequence{
	baseInput: []input.GameInput{
		input.Right, input.Right,
		tolerance: 10,
		cancelInputs: []input.GameInput{
			input.Left, input.Down, input.Up,
		},
	},
}
var HalfCircleForward = InputSequence{
	baseInput: []input.GameInput{
		input.Down, input.Right & input.Right, input.Right,
	},
}
