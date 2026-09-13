package gameplay

// HUDState is the presentation-ready snapshot for the fight HUD (SPEC §7.8).
// Pure data: HP values (the scene does pixel math), win counts, timer
// seconds, and a locale key for center text (never a locale string —
// gameplay stays locale-free).
type HUDState struct {
	P1HP, P1MaxHP  int
	P2HP, P2MaxHP  int
	P1Wins, P2Wins int
	TimerSeconds   int
	CenterTextKey  string
}

// Center text keys, resolved to locale strings by the scene.
const (
	HUDTextNone     = ""
	HUDTextKO       = "ko"
	HUDTextTimeOver = "time_over"
	HUDTextDraw     = "round_draw"
)

// HUD derives display state from the snapshot. Value receiver: read-only,
// it cannot mutate simulation (SPEC §7.8, §11.3).
func (g GameState) HUD() HUDState {
	s := HUDState{
		TimerSeconds: (g.TimerFrames + 59) / 60,
	}
	if c := g.Characters[0]; c != nil && c.StateMachine != nil {
		s.P1HP, s.P1MaxHP, s.P1Wins = c.StateMachine.HP, c.StateMachine.MaxHP, g.Wins[0]
	}
	if c := g.Characters[1]; c != nil && c.StateMachine != nil {
		s.P2HP, s.P2MaxHP, s.P2Wins = c.StateMachine.HP, c.StateMachine.MaxHP, g.Wins[1]
	}
	if g.Phase == PhaseRoundEnd {
		hp0 := s.P1HP
		hp1 := s.P2HP
		switch {
		case hp0 == hp1:
			// Timer tie or double KO (both clamped 0).
			s.CenterTextKey = HUDTextDraw
		case hp0 <= 0 || hp1 <= 0:
			s.CenterTextKey = HUDTextKO
		default:
			s.CenterTextKey = HUDTextTimeOver
		}
	}
	return s
}
