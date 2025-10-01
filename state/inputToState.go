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
			input.Down, input.Down | input.Right, input.Right, input.A,
		},
		buffer: 10,
		cancelInputs: []input.GameInput{
			input.Up, input.Left,
		},
	},
	"426A": {
		baseInput: []input.GameInput{
			input.Left, input.Down, input.Right, input.A,
		},
		buffer: 10,
	},
}

func isDashSequence(inputSeq InputSequence) bool {
	return reflect.DeepEqual(inputSeq, InputSequences["66"])
}

func DetectInputSequence(inputSeq InputSequence, inputs []input.GameInput) bool {
	if isDashSequence(inputSeq) { // needs a neutral input in between
		// Look for the pattern: Right -> Neutral(s) -> Right
		// We need to find this pattern in the recent input history

		for i := len(inputs) - 1; i >= 2; i-- { // Start from most recent, need at least 3 inputs
			// Check if current input is Right
			if inputs[i] == input.Right {
				// Look backwards for a sequence of neutrals followed by another Right
				neutralCount := 0
				foundFirstRight := false

				// Count neutrals going backwards
				for j := i - 1; j >= 0; j-- {
					if inputs[j] == input.NoInput {
						neutralCount++
						// Check if we have enough neutrals and there's a buffer limit
						if neutralCount > inputSeq.buffer {
							break // Too many neutrals, sequence invalid
						}
					} else if inputs[j] == input.Right && neutralCount > 0 {
						// Found the first Right after some neutrals
						foundFirstRight = true
						break
					} else if slices.Contains(inputSeq.cancelInputs, inputs[j]) {
						// Found a cancel input, sequence is invalid
						break
					} else {
						// Found some other input, break the sequence
						break
					}
				}

				if foundFirstRight && neutralCount > 0 && neutralCount <= inputSeq.buffer {
					return true
				}
			}
		}
		return false
	}
	// non dash sequences
	// Work backwards from the most recent inputs to find the sequence
	if len(inputs) < len(inputSeq.baseInput) {
		return false
	}

	// Try to find the sequence ending at different positions
	for endPos := len(inputs) - 1; endPos >= len(inputSeq.baseInput)-1; endPos-- {
		if matchSequenceEndingAt(inputSeq, inputs, endPos) {
			return true
		}
	}
	return false
}

// matchSequenceEndingAt checks if a sequence matches ending at the given position
func matchSequenceEndingAt(inputSeq InputSequence, inputs []input.GameInput, endPos int) bool {
	seqLen := len(inputSeq.baseInput)
	if endPos < seqLen-1 {
		return false
	}

	// Work backwards through the sequence
	inputPos := endPos
	for seqIndex := seqLen - 1; seqIndex >= 0; seqIndex-- {
		expectedInput := inputSeq.baseInput[seqIndex]
		found := false
		maxLookback := inputSeq.buffer
		if maxLookback == 0 {
			maxLookback = 3 // default buffer
		}

		// Look for the expected input within buffer range
		for lookback := 0; lookback <= maxLookback && inputPos >= 0; lookback++ {
			currentInput := inputs[inputPos]

			// Check for cancel inputs
			if slices.Contains(inputSeq.cancelInputs, currentInput) {
				return false
			}

			// Check if current input matches expected input
			// For diagonal inputs like Down|Right, allow either exact match or containing the flags
			if currentInput == expectedInput ||
				(expectedInput != input.NoInput && (currentInput&expectedInput) == expectedInput) {
				found = true
				inputPos-- // Move to the previous input for next sequence element
				break
			}

			// Move backwards, but don't count neutral inputs against buffer
			if currentInput != input.NoInput {
				// Non-neutral input that doesn't match, this counts against buffer
			}
			inputPos--
		}

		if !found {
			return false
		}
	}

	return true
}
