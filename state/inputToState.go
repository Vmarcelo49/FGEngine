package state

import (
	"fgengine/input"
	"reflect"
	"slices"
)

// basic movement

var movementInputs = map[input.GameInput]State{
	input.Right: StateDash | StateForward,
	input.Left:  StateDash | StateBackward,
	input.Down:  StateCrouching | StateNeutral,
	input.Up:    StateJump | StateNeutral,
}

// states are becoming any game action
func inputSequenceToState(inputSeq []input.GameInput) State {
	return StateDash | StateForward
}

type InputSequence struct {
	baseInput    []input.GameInput // the main input sequence
	buffer       int               // tolerance for buffering inputs
	cancelInputs []input.GameInput // inputs that can cancel the sequence
}

var InputSequences = map[string]InputSequence{ // This being a map kinda sucks for writing some loops and maybe performance?
	"66": {
		baseInput: []input.GameInput{
			input.Right, input.Right,
		},
		buffer: 10,
		cancelInputs: []input.GameInput{
			input.Left, input.Down, input.Up, // maybe i will remove down
		},
	},
	"236A": {
		baseInput: []input.GameInput{
			input.Down, input.Down & input.Right, input.Right, input.A,
		},
	},
}

func isDashSequence(inputSeq InputSequence) bool {
	return reflect.DeepEqual(inputSeq, InputSequences["66"])
}

func DetectInputSequence(inputSeq InputSequence, inputs []input.GameInput) bool {
	if isDashSequence(inputSeq) { // needs a neutral input in between
		firstRightIndex := -1
		neutralCount := 0

		// Find the first Right input
		for i, in := range inputs {
			// Check for cancel inputs
			if slices.Contains(inputSeq.cancelInputs, in) {
				return false
			}

			if in == input.Right && firstRightIndex == -1 {
				firstRightIndex = i
			} else if firstRightIndex != -1 && in == input.NoInput {
				neutralCount++
			} else if firstRightIndex != -1 && in == input.Right && neutralCount > 0 && neutralCount <= inputSeq.buffer {
				return true
			}
		}
	}
	// non dash sequences
	for i := range inputs {
		if inputs[i] == inputSeq.baseInput[0] {
			// Try to match the complete sequence starting from this position
			sequenceIndex := 1
			searchIndex := i + 1
			bufferUsed := 0

			for sequenceIndex < len(inputSeq.baseInput) && searchIndex < len(inputs) {
				// Check for cancel inputs
				if slices.Contains(inputSeq.cancelInputs, inputs[searchIndex]) {
					break
				}

				if inputs[searchIndex] == inputSeq.baseInput[sequenceIndex] {
					// Found the next input in sequence
					sequenceIndex++
					searchIndex++
					bufferUsed = 0 // Reset buffer count when we find a matching input
				} else if bufferUsed < inputSeq.buffer {
					// Still within buffer tolerance
					bufferUsed++
					searchIndex++
				} else {
					// Exceeded buffer tolerance
					break
				}
			}

			// if matched the complete sequence
			if sequenceIndex == len(inputSeq.baseInput) {
				return true
			}
		}
	}
	return false
}
