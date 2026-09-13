package gameplay

import (
	"testing"
)

func hudGameState(p1HP, p2HP int, phase MatchPhase, wins [2]int, timer int) GameState {
	p1 := poseState(100, 382, "idle")
	p1.HP = p1HP
	p2 := poseState(200, 382, "idle")
	p2.HP = p2HP
	g := hitGameState(p1, p2)
	g.Phase = phase
	g.Wins = wins
	g.TimerFrames = timer
	return g
}

func TestHUDPassthrough(t *testing.T) {
	g := hudGameState(7500, 3000, PhaseFight, [2]int{1, 0}, 5940)
	s := g.HUD()
	if s.P1HP != 7500 || s.P1MaxHP != 10000 || s.P1Wins != 1 {
		t.Fatalf("P1 wrong: %+v", s)
	}
	if s.P2HP != 3000 || s.P2MaxHP != 10000 || s.P2Wins != 0 {
		t.Fatalf("P2 wrong: %+v", s)
	}
	if s.TimerSeconds != 99 {
		t.Fatalf("seconds = %d, want 99", s.TimerSeconds)
	}
	if s.CenterTextKey != HUDTextNone {
		t.Fatalf("fight phase must show no center text, got %q", s.CenterTextKey)
	}
}

func TestHUDTimerCeil(t *testing.T) {
	for timer, want := range map[int]int{5940: 99, 60: 1, 61: 2, 59: 1, 1: 1, 0: 0} {
		g := hudGameState(10000, 10000, PhaseFight, [2]int{}, timer)
		if got := g.HUD().TimerSeconds; got != want {
			t.Errorf("timer %d -> %d seconds, want %d", timer, got, want)
		}
	}
}

func TestHUDCenterText(t *testing.T) {
	cases := map[string]struct {
		p1HP, p2HP int
		want       string
	}{
		"ko p2":        {5000, 0, HUDTextKO},
		"ko p1":        {0, 5000, HUDTextKO},
		"double ko":    {0, 0, HUDTextDraw},
		"timer tie":    {4000, 4000, HUDTextDraw},
		"timer win p1": {5000, 3000, HUDTextTimeOver},
		"timer win p2": {3000, 5000, HUDTextTimeOver},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			g := hudGameState(tc.p1HP, tc.p2HP, PhaseRoundEnd, [2]int{}, 0)
			if got := g.HUD().CenterTextKey; got != tc.want {
				t.Fatalf("key = %q, want %q", got, tc.want)
			}
		})
	}
}
