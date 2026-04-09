package input

import (
	"slices"
)

type InputSequence struct {
	baseInput    []GameInput // the main input sequence
	buffer       int         // tolerance for buffering inputs
	cancelInputs []GameInput // inputs that can cancel the sequence
	//alias []GameInput
}

// isNonDirectionalInput checks if the input is a non-directional input (A, B, C, D)
func isNonDirectionalInput(input GameInput) bool {
	directionalInputs := Up | Down | Left | Right
	return input != NoInput && (input&directionalInputs) == 0
}

var InputSequences = map[string]InputSequence{ // instead of strings, this should be an enum of common animation names
	"66": {
		baseInput: []GameInput{
			Right, NoInput, Right, // instead of NoInput, this could be an "any non-directional input" placeholder that matches any of A, B, C, D
		},
		buffer: 10,
		cancelInputs: []GameInput{
			Left, Down, Up, // maybe i will remove down
		},
	},
	"236A": {
		baseInput: []GameInput{
			Down, Down | Right, Right, A,
		},
		buffer: 10,
	},
	"426A": {
		baseInput: []GameInput{
			Left, Down, Right, A,
		},
		buffer: 10,
	},
}

func DetectInputSequence(inputSeq InputSequence, inputs []GameInput) bool {
	if len(inputs) < len(inputSeq.baseInput) {
		return false
	}

	sequenceLen := len(inputSeq.baseInput)
	inputPos := len(inputs) - 1
	buffer := inputSeq.buffer

	// Work backwards through the required sequence
	for seqIndex := sequenceLen - 1; seqIndex >= 0; seqIndex-- {
		expectedInput := inputSeq.baseInput[seqIndex]
		found := false
		bufferUsed := 0

		// Look backwards for the expected input within buffer range
		for bufferUsed <= buffer && inputPos >= 0 {
			currentInput := inputs[inputPos]

			// Check for cancel inputs (except NoInput)
			if currentInput != NoInput && slices.Contains(inputSeq.cancelInputs, currentInput) {
				return false
			}

			// for NoInput in the sequence, match any non-directional input
			if (expectedInput == NoInput && isNonDirectionalInput(currentInput)) ||
				currentInput == expectedInput ||
				(expectedInput != NoInput && (currentInput&expectedInput) == expectedInput) {
				found = true
				inputPos-- // Move to the previous input for next sequence element
				break
			}

			// Move backwards
			inputPos--

			// Only count non-neutral inputs against buffer
			if currentInput != NoInput {
				bufferUsed++
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// CheckSpecialMove checks the input history for any special move commands and returns the name of the detected command, or an empty string if none are detected.
func CheckSpecialMove(inputs []GameInput) string {
	for name, seq := range InputSequences {
		if DetectInputSequence(seq, inputs) {
			return name
		}
	}
	return ""
}
