package input

type InputSequence struct {
	baseInput []GameInput // the main input sequence
	buffer    int         // tolerance for buffering inputs
	//alias []GameInput // simplified inputs
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

func CheckInputSequences(inputs []GameInput) string {
	detected := []string{}
	for name, seq := range InputSequences {
		if DetectInputSequence(seq, inputs) {
			detected = append(detected, name)
		}
	}
	// Priority order first is higher:
	// 1, Supers and Specials
	// 2, normals
	// 3, movements, walks and jumps

	// if no special move was detected, check for single inputs to trigger normals and movements
	if len(detected) == 0 {
		return CheckSingleInput(inputs[len(inputs)-1])
	}
	return detected[len(detected)-1] // last one probably is the correct one...
}

func CheckSingleInput(inputs GameInput) string {
	// with priority order
	if inputs.IsPressed(D) {
		return "D"
	}
	if inputs.IsPressed(C) {
		return "C"
	}
	if inputs.IsPressed(B) {
		return "B"
	}
	if inputs.IsPressed(A) {
		return "A"
	}
	if inputs.IsPressed(Up) && inputs.IsPressed(Right) {
		return "9"
	}
	if inputs.IsPressed(Up) && inputs.IsPressed(Left) {
		return "7"
	}
	if inputs.IsPressed(Down) && inputs.IsPressed(Right) {
		return "3"
	}
	if inputs.IsPressed(Down) && inputs.IsPressed(Left) {
		return "1"
	}
	if inputs.IsPressed(Up) {
		return "8"
	}
	if inputs.IsPressed(Down) {
		return "2"
	}
	if inputs.IsPressed(Left) {
		return "4"
	}
	if inputs.IsPressed(Right) {
		return "6"
	}

	return ""
}
