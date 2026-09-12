package input

// InputSequence is one named motion: expected corrected inputs, oldest
// first, matched backwards from the newest history entry.
type InputSequence struct {
	Name      string // intent produced on match, e.g. "236A"
	BaseInput []GameInput
	Buffer    int // tolerance for buffering inputs (non-neutral frames)
}

// InputSequences is the engine-standard sequence set in explicit priority
// order: special motions first, then dashes (movement). Evaluation order is
// this declaration order — never map iteration order (SPEC §3.3, §5.5).
// x in the motion names ranges over the A/B/C/D buttons.
var InputSequences = []InputSequence{
	// 236x: Down, Down-Forward, Forward + x
	{Name: "236A", BaseInput: []GameInput{Down, Down | Right, Right, A}, Buffer: 10},
	{Name: "236B", BaseInput: []GameInput{Down, Down | Right, Right, B}, Buffer: 10},
	{Name: "236C", BaseInput: []GameInput{Down, Down | Right, Right, C}, Buffer: 10},
	{Name: "236D", BaseInput: []GameInput{Down, Down | Right, Right, D}, Buffer: 10},
	// 214x: Down, Down-Back, Back + x
	{Name: "214A", BaseInput: []GameInput{Down, Down | Left, Left, A}, Buffer: 10},
	{Name: "214B", BaseInput: []GameInput{Down, Down | Left, Left, B}, Buffer: 10},
	{Name: "214C", BaseInput: []GameInput{Down, Down | Left, Left, C}, Buffer: 10},
	{Name: "214D", BaseInput: []GameInput{Down, Down | Left, Left, D}, Buffer: 10},
	// 623x: Forward, Down, Down-Forward + x
	{Name: "623A", BaseInput: []GameInput{Right, Down, Down | Right, A}, Buffer: 10},
	{Name: "623B", BaseInput: []GameInput{Right, Down, Down | Right, B}, Buffer: 10},
	{Name: "623C", BaseInput: []GameInput{Right, Down, Down | Right, C}, Buffer: 10},
	{Name: "623D", BaseInput: []GameInput{Right, Down, Down | Right, D}, Buffer: 10},
	// 423x: Back, Down, Down-Back + x
	{Name: "423A", BaseInput: []GameInput{Left, Down, Down | Left, A}, Buffer: 10},
	{Name: "423B", BaseInput: []GameInput{Left, Down, Down | Left, B}, Buffer: 10},
	{Name: "423C", BaseInput: []GameInput{Left, Down, Down | Left, C}, Buffer: 10},
	{Name: "423D", BaseInput: []GameInput{Left, Down, Down | Left, D}, Buffer: 10},
	// 22x: Down, Down + x
	{Name: "22A", BaseInput: []GameInput{Down, Down, A}, Buffer: 10},
	{Name: "22B", BaseInput: []GameInput{Down, Down, B}, Buffer: 10},
	{Name: "22C", BaseInput: []GameInput{Down, Down, C}, Buffer: 10},
	{Name: "22D", BaseInput: []GameInput{Down, Down, D}, Buffer: 10},
	// 246x: Down, Back, Forward + x (simplified half-circle, no corners)
	{Name: "246A", BaseInput: []GameInput{Down, Left, Right, A}, Buffer: 10},
	{Name: "246B", BaseInput: []GameInput{Down, Left, Right, B}, Buffer: 10},
	{Name: "246C", BaseInput: []GameInput{Down, Left, Right, C}, Buffer: 10},
	{Name: "246D", BaseInput: []GameInput{Down, Left, Right, D}, Buffer: 10},
	// 642x: Forward, Back, Down + x (simplified half-circle, no corners)
	{Name: "642A", BaseInput: []GameInput{Right, Left, Down, A}, Buffer: 10},
	{Name: "642B", BaseInput: []GameInput{Right, Left, Down, B}, Buffer: 10},
	{Name: "642C", BaseInput: []GameInput{Right, Left, Down, C}, Buffer: 10},
	{Name: "642D", BaseInput: []GameInput{Right, Left, Down, D}, Buffer: 10},
	// Dashes (movement — lowest priority)
	{Name: "66", BaseInput: []GameInput{Right, NoInput, Right}, Buffer: 10},
	{Name: "44", BaseInput: []GameInput{Left, NoInput, Left}, Buffer: 10},
}

// discreteIntents is the exact fire-once set (SPEC §5.5): buttons, dashes,
// and all special motions. Everything else is continuous. Map lookup is
// order-independent and deterministic.
var discreteIntents = func() map[string]bool {
	m := map[string]bool{"A": true, "B": true, "C": true, "D": true}
	for _, seq := range InputSequences {
		m[seq.Name] = true
	}
	return m
}()

// isNonDirectionalInput checks if the input is a non-directional input (A, B, C, D)
func isNonDirectionalInput(input GameInput) bool {
	directionalInputs := Up | Down | Left | Right
	return input != NoInput && (input&directionalInputs) == 0
}

func DetectInputSequence(inputSeq InputSequence, inputs []GameInput) bool {
	if len(inputs) < len(inputSeq.BaseInput) {
		return false
	}

	sequenceLen := len(inputSeq.BaseInput)
	inputPos := len(inputs) - 1
	buffer := inputSeq.Buffer

	// Work backwards through the required sequence
	for seqIndex := sequenceLen - 1; seqIndex >= 0; seqIndex-- {
		expectedInput := inputSeq.BaseInput[seqIndex]
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
	// Priority order is the InputSequences declaration order (first match
	// wins): supers/specials, then dashes. Never map iteration order.
	for _, seq := range InputSequences {
		if DetectInputSequence(seq, inputs) {
			return seq.Name
		}
	}
	// if no special move was detected, check for single inputs to trigger normals and movements
	if len(inputs) == 0 {
		return ""
	}
	return CheckSingleInput(inputs[len(inputs)-1])
}

// CheckInputIntent returns a frame intent while preventing repeated triggers
// for discrete button-based actions (normals/specials) held in input history.
func CheckInputIntent(inputs []GameInput) string {
	if len(inputs) == 0 {
		return ""
	}

	current := CheckInputSequences(inputs)
	if current == "" || !isDiscreteIntent(current) || len(inputs) == 1 {
		return current
	}

	previous := CheckInputSequences(inputs[:len(inputs)-1])
	if previous == current {
		return ""
	}

	return current
}

func isDiscreteIntent(intent string) bool {
	return discreteIntents[intent]
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
