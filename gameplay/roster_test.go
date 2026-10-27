package gameplay

import (
	"fgengine/character"
	"fgengine/constants"
	"fgengine/input"
	"testing"
)

type rosterCase struct {
	name    string
	special string
	motion  []input.GameInput
}

var playableRoster = []rosterCase{
	{"Brasa", "236A", []input.GameInput{input.Down, input.Down | input.Right, input.Right, input.A}},
	{"Muralha", "623B", []input.GameInput{input.Right, input.Down, input.Down | input.Right, input.B}},
	{"Iris", "236C", []input.GameInput{input.Down, input.Down | input.Right, input.Right, input.C}},
	{"Vulto", "214D", []input.GameInput{input.Down, input.Down | input.Left, input.Left, input.D}},
}

func rosterMatch(t *testing.T, first, second string) *GameState {
	t.Helper()
	p1, err := character.LoadCharacter(first, 1)
	if err != nil {
		t.Fatalf("load %s: %v", first, err)
	}
	p2, err := character.LoadCharacter(second, 2)
	if err != nil {
		t.Fatalf("load %s: %v", second, err)
	}
	g := NewGameState(p1, p2, 42)
	return &g
}

func TestPlayableRosterNormalInputs(t *testing.T) {
	t.Chdir(gameplayRepoRoot(t))
	buttons := []struct {
		name  string
		value input.GameInput
	}{{"A", input.A}, {"B", input.B}, {"C", input.C}, {"D", input.D}}
	postures := []struct {
		name     string
		prefix   string
		down     bool
		airborne bool
	}{{"standing", "", false, false}, {"crouching", "2", true, false}, {"airborne", "air", false, true}}

	for _, fighter := range playableRoster {
		for _, posture := range postures {
			for _, button := range buttons {
				t.Run(fighter.name+"/"+posture.name+"/"+button.name, func(t *testing.T) {
					g := rosterMatch(t, fighter.name, "PlaceHolder")
					sm := g.Characters[0].StateMachine
					if posture.airborne {
						sm.Position.Y = constants.GroundLevelY - 20
					}
					press := button.value
					if posture.down {
						press |= input.Down
					}
					g.Update([2]input.GameInput{press, input.NoInput})
					want := posture.prefix + button.name
					if got := sm.AnimPlayer.ActiveAnimationName(); got != want {
						t.Fatalf("input selected %q, want %q", got, want)
					}
					anim := sm.AnimPlayer.Animations[want]
					if len(anim.Sprites) < 4 || len(anim.FrameData) < 4 {
						t.Fatalf("%s needs four sprite and frame-data entries; got %d and %d", want, len(anim.Sprites), len(anim.FrameData))
					}
				})
			}
		}
	}
}

func TestPlayableRosterSpecialInputsBothSides(t *testing.T) {
	t.Chdir(gameplayRepoRoot(t))
	for _, fighter := range playableRoster {
		for _, side := range []int{0, 1} {
			t.Run(fighter.name+map[int]string{0: "/left", 1: "/right"}[side], func(t *testing.T) {
				players := [2]string{"PlaceHolder", "PlaceHolder"}
				players[side] = fighter.name
				g := rosterMatch(t, players[0], players[1])
				for _, step := range fighter.motion {
					if side == 1 {
						step = mirrorRosterInput(step)
					}
					var in [2]input.GameInput
					in[side] = step
					g.Update(in)
				}
				if got := g.Characters[side].StateMachine.AnimPlayer.ActiveAnimationName(); got != fighter.special {
					t.Fatalf("special motion selected %q, want %q", got, fighter.special)
				}
			})
		}
	}
}

func mirrorRosterInput(in input.GameInput) input.GameInput {
	if in&input.Left != 0 {
		return (in &^ input.Left) | input.Right
	}
	if in&input.Right != 0 {
		return (in &^ input.Right) | input.Left
	}
	return in
}

func TestPlayableRosterHitsAndGuards(t *testing.T) {
	t.Chdir(gameplayRepoRoot(t))
	for _, fighter := range playableRoster {
		t.Run(fighter.name+"/hit", func(t *testing.T) {
			g := rosterMatch(t, fighter.name, "PlaceHolder")
			g.Characters[0].StateMachine.Position.X = 300
			def := g.Characters[1].StateMachine
			def.Position.X = 360
			startingHP := def.HP
			for frame := 0; frame < 20; frame++ {
				press := input.NoInput
				if frame == 0 {
					press = input.A
				}
				g.Update([2]input.GameInput{press, input.NoInput})
				if def.HP < startingHP {
					if got := def.AnimPlayer.ActiveAnimationName(); got != "hurt" {
						t.Fatalf("damage caused reaction %q, want hurt", got)
					}
					return
				}
			}
			t.Fatal("standing A never damaged the nearby opponent")
		})
		t.Run(fighter.name+"/guard", func(t *testing.T) {
			g := rosterMatch(t, "PlaceHolder", fighter.name)
			g.Characters[0].StateMachine.Position.X = 300
			def := g.Characters[1].StateMachine
			def.Position.X = 360
			startingHP := def.HP
			for frame := 0; frame < 20; frame++ {
				attack := input.NoInput
				if frame == 5 {
					attack = input.A
				}
				g.Update([2]input.GameInput{attack, input.Right})
				if def.AnimPlayer.ActiveAnimationName() == "blockHit" {
					if def.HP != startingHP || def.StunFrames <= 0 {
						t.Fatalf("guard: HP=%d, want %d; stun=%d", def.HP, startingHP, def.StunFrames)
					}
					return
				}
			}
			t.Fatal("the nearby opponent's attack never entered blockstun")
		})
	}
}
