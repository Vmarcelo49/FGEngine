package input

import (
	"testing"
)

func findSequence(name string) InputSequence {
	for _, seq := range InputSequences {
		if seq.Name == name {
			return seq
		}
	}
	panic("unknown sequence in test: " + name)
}

func TestDetectDash66(t *testing.T) {
	history := []GameInput{Right, NoInput, Right}
	if !DetectInputSequence(findSequence("66"), history) {
		t.Fatal("66 should match Right, neutral, Right")
	}
}

func TestDetectDash44(t *testing.T) {
	history := []GameInput{Left, NoInput, Left}
	if !DetectInputSequence(findSequence("44"), history) {
		t.Fatal("44 should match Left, neutral, Left")
	}
}

func TestDetect236A(t *testing.T) {
	history := []GameInput{Down, Down | Right, Right, A}
	if !DetectInputSequence(findSequence("236A"), history) {
		t.Fatal("236A should match the fireball motion")
	}
}

func TestBufferBoundary(t *testing.T) {
	seq := findSequence("66")

	// The NoInput slot needs a real button press; buffer only absorbs
	// extra frames around it. 10 skipped non-neutrals: within buffer.
	ok := []GameInput{Right}
	for i := 0; i < 10; i++ {
		ok = append(ok, Up)
	}
	ok = append(ok, A, Right)
	if !DetectInputSequence(seq, ok) {
		t.Fatal("66 should match with 10 buffered non-neutral inputs")
	}

	// 11 skipped non-neutrals before the button: buffer exceeded.
	tooFar := []GameInput{Right}
	for i := 0; i < 11; i++ {
		tooFar = append(tooFar, Up)
	}
	tooFar = append(tooFar, A, Right)
	if DetectInputSequence(seq, tooFar) {
		t.Fatal("66 must not match with 11 buffered non-neutral inputs")
	}

	// Neutral frames never consume buffer, however many.
	neutral := []GameInput{Right}
	for i := 0; i < 50; i++ {
		neutral = append(neutral, NoInput)
	}
	neutral = append(neutral, A)
	for i := 0; i < 50; i++ {
		neutral = append(neutral, NoInput)
	}
	neutral = append(neutral, Right)
	if !DetectInputSequence(seq, neutral) {
		t.Fatal("66 should match across an arbitrarily long neutral gap")
	}
}

// Every declared sequence must be reachable from its own exact pattern.
// (Catches both dead entries and accidental cross-matches by earlier ones
// when asserted through CheckInputSequences.)
func TestFullTableReachability(t *testing.T) {
	for _, seq := range InputSequences {
		if !DetectInputSequence(seq, seq.BaseInput) {
			t.Errorf("%s does not match its own exact pattern", seq.Name)
		}
		if got := CheckInputSequences(seq.BaseInput); got != seq.Name {
			t.Errorf("exact pattern for %s resolved to %q (priority order suspect)", seq.Name, got)
		}
	}
}

// A|B satisfies both 236A and 236B patterns; declaration order decides.
func TestFirstMatchWinsPriority(t *testing.T) {
	history := []GameInput{Down, Down | Right, Right, A | B}
	if got := CheckInputSequences(history); got != "236A" {
		t.Fatalf("expected 236A (declared first), got %q", got)
	}
}

func TestCheckInputIntentSuppression(t *testing.T) {
	first := CheckInputIntent([]GameInput{NoInput, A})
	if first != "A" {
		t.Fatalf("expected first A press to yield intent A, got %q", first)
	}
	held := CheckInputIntent([]GameInput{NoInput, A, A})
	if held != "" {
		t.Fatalf("expected held A to be suppressed, got %q", held)
	}
	// Continuous intents are never suppressed.
	walk := CheckInputIntent([]GameInput{NoInput, Right, Right})
	if walk != "6" {
		t.Fatalf("expected held Right to keep yielding 6, got %q", walk)
	}
}

func TestApplySOCD(t *testing.T) {
	if got := ApplySOCD(Left | Right); got != NoInput {
		t.Fatalf("Left+Right should clear to neutral, got %v", got)
	}
	if got := ApplySOCD(Up | Down); got != NoInput {
		t.Fatalf("Up+Down should clear to neutral, got %v", got)
	}
	if got := ApplySOCD(Up | A); got != Up|A {
		t.Fatalf("Up+A must pass through, got %v", got)
	}
	if got := ApplySOCD(NoInput); got != NoInput {
		t.Fatalf("neutral must stay neutral, got %v", got)
	}
}

func TestCheckSingleInputPriority(t *testing.T) {
	cases := map[GameInput]string{
		Down | A:   "A",
		Up | Right: "9",
		Up | Left:  "7",
		Down:       "2",
		Left:       "4",
		NoInput:    "",
		D | Up:     "D",
	}
	for in, want := range cases {
		if got := CheckSingleInput(in); got != want {
			t.Errorf("CheckSingleInput(%v) = %q, want %q", in, got, want)
		}
	}
}

func TestIsDiscreteIntent(t *testing.T) {
	discrete := []string{"A", "B", "C", "D", "44", "66", "236A", "214D", "22B", "642C"}
	for _, name := range discrete {
		if !isDiscreteIntent(name) {
			t.Errorf("expected %s to be discrete", name)
		}
	}
	continuous := []string{"4", "6", "8", "idle", "hurt", "airBlock", "crouchBlock", "fall", ""}
	for _, name := range continuous {
		if isDiscreteIntent(name) {
			t.Errorf("expected %s to be continuous, got discrete", name)
		}
	}
}
